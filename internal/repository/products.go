package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ProductsRepository struct {
	db *sqlx.DB
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func NewProductsRepository(db *sqlx.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) Create(ctx context.Context, product domain.Product) (int, error) {
	var id int
	row := r.db.QueryRow(`INSERT INTO products (article, name, price, manufacturer, sellerId, deleted, rating) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, product.Article, product.Name,
		product.Price, product.Manufacturer, product.SellerId, product.Deleted, product.Rating)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("Creation error of product: %s", err)
		return 0, ErrInternal
	}
	return id, nil
}

func (r *ProductsRepository) Get(ctx context.Context, productId int) (domain.Product, error) {
	var product domain.Product

	row := r.db.QueryRowxContext(ctx, `SELECT * FROM products WHERE id = $1`, productId)
	if err := row.StructScan(&product); err != nil {
		logrus.Errorf("Error get product from postgresql: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			return product, ErrProductNotFound
		}
		return product, ErrInternal
	}

	return product, nil

}

func (r *ProductsRepository) GetMyAll(ctx context.Context, userId int) ([]domain.Product, error) {
	products := []domain.Product{}
	err := r.db.Select(&products, "SELECT * FROM products WHERE user_id = $1", userId)
	if err != nil {
		logrus.Errorf("Error get products from postgresql: %s", err)
		return nil, ErrInternal
	}
	return products, nil
}

func (r *ProductsRepository) Delete(ctx context.Context, productId int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", productId)
	if err != nil {
		logrus.Errorf("Error deleting product from posgresql: %s", err)
		return ErrInternal
	}
	return nil
}

func (r *ProductsRepository) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Product, error) {
	err := filter.Fill_defaults()
	if err != nil {
		logrus.Errorf("Error limit too high: %s", err)
		return nil, ErrInternal
	}

	var products []domain.Product
	var names []string
	var values []interface{}
	argId := 1

	addProp := func(name string, sign string, value interface{}) {
		names = append(names, fmt.Sprintf("%s %s $%d", name, sign, argId))
		values = append(values, value)
		argId++
	}

	if filter.Article != nil {
		addProp("article", "=", *filter.Article)
	}
	if filter.Name != nil {
		addProp("name", "=", *filter.Name)
	}
	if filter.MinPrice != nil {
		addProp("price", ">", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		addProp("price", "<", *filter.MaxPrice)
	}
	if filter.Manufacturer != nil {
		addProp("manafacturer", "=", *filter.Manufacturer)
	}
	if filter.Rating != nil {
		addProp("article", "=", *filter.Rating)
	}

	setQuery := strings.Join(names, ", ")
	values = append(values, *filter.Limit, *filter.Offset)
	query := fmt.Sprintf(`SELECT * FROM products WHERE %s ORDER BY rating DESC LIMIT $%d OFFSET $%d`,
		setQuery, argId, argId+1)
	err = r.db.Select(&products, query, values...)
	if err != nil {
		logrus.Errorf("sql-query: %s", query)
		logrus.Errorf("error get products: %s", err)
		return products, ErrInternal
	}

	return products, nil
}
