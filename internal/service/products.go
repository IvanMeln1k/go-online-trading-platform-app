package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/sirupsen/logrus"
)

type ProductsService struct {
	productsRepo repository.Products
	usersRepo    repository.Users
}

var (
	ErrProductNotFound = errors.New("product not found")
)

func NewProductsService(ProductsRepo repository.Products, usersRepo repository.Users) *ProductsService {
	return &ProductsService{ProductsRepo, usersRepo}
}

func (s *ProductsService) GetAll(ctx context.Context, userId int) ([]domain.Product, error) {
	_, err := s.usersRepo.GetById(ctx, userId)
	if err != nil {
		logrus.Errorf("Service products error GetAll: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}
	products, err := s.productsRepo.GetAll(ctx, userId)
	if err != nil {
		logrus.Errorf("Service GetAllProducts calling repository error: %s", err)
		return nil, ErrInternal
	}

	return products, nil
}

func (s *ProductsService) Get(ctx context.Context, productId int) (domain.Product, error) {
	product, err := s.productsRepo.Get(ctx, productId)
	if err != nil {
		logrus.Errorf("Service GetProduct calling repository error: %s", err)
		if errors.Is(repository.ErrProductNotFound, err) {
			return product, ErrProductNotFound
		}
		return product, ErrInternal
	}

	return product, nil
}

func (s *ProductsService) Create(ctx context.Context, userId int, product domain.Product) (int, error) {
	_, err := s.usersRepo.GetById(ctx, userId)
	if err != nil {
		logrus.Errorf("Service CreateProduct getUser from repository when creating product error: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return 0, ErrUserNotFound
		}
		return 0, ErrInternal
	}
	id, err := s.productsRepo.Create(ctx, product)
	if err != nil {
		logrus.Errorf("Service CreateProduct are broken when calling repository error: %s", err)
		return 0, ErrInternal
	}
	return id, nil
}

func (s *ProductsService) Delete(ctx context.Context, productId int) error {
	err := s.productsRepo.Delete(ctx, productId)
	if err != nil {
		logrus.Errorf("Service DeleteProduct error when deleting product from repo: %s", err)
		return ErrInternal
	}
	return nil
}
