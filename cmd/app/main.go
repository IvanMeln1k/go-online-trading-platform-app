package main

import (
	"os"
	"strconv"
	"time"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/handler"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/server"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/database"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/password"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err)
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env: %s", err)
	}

	db, err := database.NewPostgresDB(database.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASS"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error connect to postgres db: %s", err)
	}

	rdbNum, err := strconv.Atoi(viper.GetString("rdb.db"))
	if err != nil {
		logrus.Fatalf("invalid number redis db")
	}
	rdb := database.NewRedisDB(database.RedisConfig{
		Host:     viper.GetString("rdb.host"),
		Port:     viper.GetString("rdb.port"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       rdbNum,
	})

	emailSender, err := email.NewEmailSender(email.EmailSenderConfig{
		Email: viper.GetString("smtp.email"),
		Pass:  os.Getenv("SMTP_PASS"),
		Host:  viper.GetString("smtp.host"),
		Port:  viper.GetString("smtp.port"),
	})
	if err != nil {
		logrus.Fatalf("error creating email sender: %s", err)
	}
	accessTTL, err := time.ParseDuration(viper.GetString("tokens.accessTTL"))
	if err != nil {
		logrus.Fatalf("error parsing accessTTL: %s", err)
	}
	emailTTL, err := time.ParseDuration(viper.GetString("tokens.emailTTL"))
	if err != nil {
		logrus.Fatalf("error parsing emailTTL: %s", err)
	}
	refreshTTL, err := time.ParseDuration(viper.GetString("tokens.refreshTTL"))
	if err != nil {
		logrus.Fatalf("error parsing refreshTTL: %s", err)
	}
	tokenManager := tokens.NewTokenManager(os.Getenv("JWT_KEY"), accessTTL, emailTTL)
	passwordManager := password.NewPasswordManager(os.Getenv("SALT"))

	repos := repository.NewRepository(db, rdb)
	services := service.NewService(service.Deps{
		Repo:            repos,
		TokenManager:    tokenManager,
		PasswordManager: passwordManager,
		EmailSender:     emailSender,
		RefreshTTL:      refreshTTL,
	})
	handlers := handler.NewHandler(handler.Deps{
		Services:     services,
		TokenManager: tokenManager,
	})

	srv := new(server.Server)
	srvCfg := server.ServerConfig{
		Host: viper.GetString("app.host"),
		Port: viper.GetString("app.port"),
	}
	if err := srv.Run(srvCfg, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
