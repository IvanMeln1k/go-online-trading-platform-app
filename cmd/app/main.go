package main

import (
	"os"
	"strconv"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/handler"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/server"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/service"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/database"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
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
		Password: viper.GetString("rdb.password"),
		DB:       rdbNum,
	})

	emailSender, err := email.NewEmailSender("IvanMelnikovF@gmail.com", os.Getenv("SMTP_PASS"),
		viper.GetString("smtp.host"), viper.GetString("smtp.port"))
	if err != nil {
		logrus.Fatalf("error create email sender: %s", err)
	}

	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos, emailSender)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	srvCfg := server.ServerConfig{
		Host: "localhost",
		Port: "8000",
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
