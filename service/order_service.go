package service

import (
	"FlashSaleGo/model"
	"FlashSaleGo/repository"
)

type IOrderService interface {
	GetOrderByID(int64) (*model.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*model.Order) error
	InsertOrder(*model.Order) (int64, error)
	GetAllOrder() ([]*model.Order, error)
	//GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*model.Message) (int64, error)
}

type OrderService struct {
	OrderRepository repository.IOrderRepository
}

func (o *OrderService) GetOrderByID(id int64) (*model.Order, error) {
	return o.OrderRepository.SelectByKey(id)
}

func (o *OrderService) DeleteOrderByID(id int64) bool {
	return o.OrderRepository.Delete(id)
}

func (o *OrderService) UpdateOrder(order *model.Order) error {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *model.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*model.Order, error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) InsertOrderByMessage(message *model.Message) (int64, error) {
	order := &model.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: model.OrderSuccess,
	}
	return o.InsertOrder(order)
}

func NewOrderService(repository repository.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repository}
}
