package order

import (
	"FlashSaleGo/rabbitmq"
	"golang.org/x/net/context"
	"log"
)

type ServerOrder struct {
	// RabbitMQ
	rabbitMq *rabbitmq.RabbitMQ
	// LocalHost
	localHost string
	port      string
}

func (s *ServerOrder) mustEmbedUnimplementedOrderServiceServer() {
	panic("implement me")
}

func (s *ServerOrder) MakeOrder(ctx context.Context, message *OrderInfo) (*OrderResult, error) {
	log.Printf("From Client: username: %v, productid: %v", message.Username, message.ProductID)
	return &OrderResult{IsOrderSuccess: "true"}, nil
}

func (s *ServerOrder) Destroy() {
	s.rabbitMq.Destroy()
}

func NewOrderServer(localHost, port string) (server *ServerOrder) {
	server = &ServerOrder{
		rabbitMq:  rabbitmq.NewRabbitMQSimple("iRaidenProduct"),
		localHost: localHost,
		port:      port,
	}
	return
}
