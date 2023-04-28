package repository

import (
	"FlashSaleGo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	Conn() error
	Insert(*model.Order) (int64, error)
	Delete(int64) bool
	Update(order *model.Order) error
	SelectByKey(int64) (*model.Order, error)
	SelectAll() ([]*model.Order, error)
}

type OrderRepository struct {
	myGormConn *gorm.DB
}

const defaultOrderDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/flashSaleDB?charset=utf8"

func (o *OrderRepository) Conn() (err error) {
	if o.myGormConn == nil {
		o.myGormConn, err = gorm.Open(mysql.Open(defaultOrderDB), &gorm.Config{})
		if err != nil {
			return
		}
		o.myGormConn.AutoMigrate(&model.Product{})
	}
	return nil
}

func (o *OrderRepository) Insert(order *model.Order) (orderID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	//user are not allowed to set product id themselves
	order.ID = 0
	result := o.myGormConn.Create(order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.ID, nil
}

func (o *OrderRepository) Delete(id int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	var order model.Order
	result := o.myGormConn.Delete(&order, id)
	if result.Error != nil {
		return false
	}
	return true
}

func (o *OrderRepository) Update(order *model.Order) (err error) {
	if err = o.Conn(); err != nil {
		return
	}
	if result := o.myGormConn.Save(order); result.Error != nil {
		return result.Error
	}
	return
}

func (o *OrderRepository) SelectByKey(id int64) (order *model.Order, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	order = &model.Order{}
	if result := o.myGormConn.First(order, id); result.Error != nil {
		return nil, result.Error
	}
	return
}

func (o *OrderRepository) SelectAll() (orderArray []*model.Order, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	if result := o.myGormConn.Find(&orderArray); result.Error != nil {
		return nil, result.Error
	}
	return
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		myGormConn: db,
	}
}
