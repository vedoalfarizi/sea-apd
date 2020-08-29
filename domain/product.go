package domain

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/williamchang80/sea-apd/dto/request"
)

type Product struct {
	gorm.Model
	Name        string `json:"id"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type ProductUsecase interface {
	GetProducts() ([]Product, error)
	GetProductById(string) (Product, error)
	CreateProduct(request.Product) error
	UpdateProduct(string, request.Product) error
	DeleteProduct(string) error
}

type ProductRepository interface {
	GetProducts() ([]Product, error)
	GetProductById(string) (Product, error)
	CreateProduct(Product) error
	UpdateProduct(string, Product) error
	DeleteProduct(string) error
}

type ProductController interface {
	GetProducts(echo.Context) error
	GetProductById(echo.Context) error
	CreateProduct(echo.Context) error
	UpdateProduct(echo.Context) error
	DeleteProduct(echo.Context) error
}

func NewProduct(name string, desc string, price int, image string, stock int) *Product {
	return &Product{
		Name:        name,
		Description: desc,
		Price:       price,
		Image:       image,
		Stock:       stock,
	}
}
