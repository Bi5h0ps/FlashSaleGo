package service

import (
	"FlashSaleGo/model"
	"FlashSaleGo/repository"
)

type IProductService interface {
	GetProductById(int64) (*model.Product, error)
	GetAllProduct() ([]*model.Product, error)
	DeleteProductById(int64) bool
	InsertProduct(*model.Product) (int64, error)
	UpdateProduct(*model.Product) error
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productRepository repository.IProduct
}

func (p ProductService) GetProductById(i int64) (*model.Product, error) {
	return p.productRepository.SelectByKey(i)
}

func (p ProductService) GetAllProduct() ([]*model.Product, error) {
	return p.productRepository.SelectAll()
}

func (p ProductService) DeleteProductById(i int64) bool {
	return p.productRepository.Delete(i)
}

func (p ProductService) InsertProduct(product *model.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p ProductService) UpdateProduct(product *model.Product) error {
	return p.productRepository.Update(product)
}

func (p ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}

func NewProductService(repository repository.IProduct) IProductService {
	return &ProductService{productRepository: repository}
}
