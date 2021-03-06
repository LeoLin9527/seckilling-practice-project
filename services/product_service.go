package services

import (
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"seckilling-practice-project/respsoiories"
)

type IProductService interface {
	GetProductByID(int64) (*models.Product, error)
	GetAllProduct() ([]*models.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *models.Product) (int64, error)
	UpdateProduct(product *models.Product) error
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productRepository respsoiories.Iproduct
}

func (p *ProductService) GetProductByID(productID int64) (*models.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*models.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *models.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *models.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}

func NewProductService(productRepository respsoiories.Iproduct) IProductService {
	return &ProductService{productRepository: productRepository}
}

func DefaultProductService() IProductService {
	mysqlCon, err := common.DefaultDb()
	if err != nil {
		panic(err)
	}
	productRepo := respsoiories.NewProductManager("product", mysqlCon)
	return NewProductService(productRepo)

}
