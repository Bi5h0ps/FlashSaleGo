package main

import (
	"FlashSaleGo/rabbitmq"
	"FlashSaleGo/repository"
	"FlashSaleGo/service"
)

func main() {
	repoOrder := repository.NewOrderRepository(nil)
	serviceOrder := service.NewOrderService(repoOrder)
	repoProduct := repository.NewProductRepository(nil)
	serviceProduct := service.NewProductService(repoProduct)
	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("iRaidenProduct")
	rabbitmqConsumeSimple.ConsumeSimple(serviceOrder, serviceProduct)
}
