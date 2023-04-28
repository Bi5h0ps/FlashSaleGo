package repository

import (
	"FlashSaleGo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IProduct interface {
	Conn() error
	Insert(*model.Product) (int64, error)
	Delete(int64) bool
	Update(*model.Product) error
	SelectByKey(int64) (*model.Product, error)
	SelectAll() ([]*model.Product, error)
	SubProductNum(productID int64) error
}

type ProductRepository struct {
	table      string
	myGormConn *gorm.DB
}

func (p *ProductRepository) Conn() (err error) {
	dsn := "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/flashSaleDB?charset=utf8"
	if p.myGormConn == nil {
		p.myGormConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}
		p.myGormConn.AutoMigrate(&model.Product{})
	}
	if p.table == "" {
		p.table = "product"
	}
	return nil
}

func (p *ProductRepository) Insert(product *model.Product) (productId int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	//user are not allowed to set product id themselves
	product.ID = 0
	result := p.myGormConn.Create(product)
	if result.Error != nil {
		return 0, result.Error
	}
	return product.ID, nil
}

func (p *ProductRepository) Delete(id int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	var product model.Product
	result := p.myGormConn.Delete(&product, id)
	if result.Error != nil {
		return false
	}
	return true
}

func (p *ProductRepository) Update(product *model.Product) (err error) {
	if err = p.Conn(); err != nil {
		return
	}
	if result := p.myGormConn.Save(product); result.Error != nil {
		return result.Error
	}
	return
}

func (p *ProductRepository) SelectByKey(productID int64) (productResult *model.Product, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	productResult = &model.Product{}
	if result := p.myGormConn.First(productResult, productID); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (p *ProductRepository) SelectAll() (productArray []*model.Product, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	if result := p.myGormConn.Find(&productArray); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (p *ProductRepository) SubProductNum(productID int64) (err error) {
	if err = p.Conn(); err != nil {
		return
	}
	result := p.myGormConn.Model(&model.Product{}).Where("ID", productID).
		UpdateColumn("productNum", gorm.Expr("productNum - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	return
}

func NewProductRepository(table string, db *gorm.DB) IProduct {
	return &ProductRepository{table: table, myGormConn: db}
}
