package order

import (
	"FlashSaleGo/distributed"
	"FlashSaleGo/rabbitmq"
	"errors"
	"golang.org/x/net/context"
	"strconv"
)

type ServerOrder struct {
	// RabbitMQ
	rabbitMq *rabbitmq.RabbitMQ
	// LocalHost
	localHost string
	port      string
	// Distributed
	accessControlUnit *distributed.AccessControl
	// Consistant Hashing
	hashConsistent *distributed.Consistent
}

func (s *ServerOrder) mustEmbedUnimplementedOrderServiceServer() {
	panic("implement me")
}

func (s *ServerOrder) MakeOrder(ctx context.Context, message *OrderInfo) (*OrderResult, error) {
	//1.distributed auth
	right := s.accessControlUnit.GetDistributedRight(message.Username, s.hashConsistent, s.localHost)
	if !right {
		return &OrderResult{IsOrderSuccess: "false"}, errors.New("distributed user auth failed")
	}
	return &OrderResult{IsOrderSuccess: "success"}, nil
}

func (s *ServerOrder) GetUserCloudData(ctx context.Context, message *UserInfo) (*UserCloudData, error) {
	data, err := s.accessControlUnit.GetDataFromMap(strconv.FormatInt(message.UserID, 10))
	return &UserCloudData{TimeStamp: data.LastOrderTime}, err
}

func (s *ServerOrder) Destroy() {
	s.rabbitMq.Destroy()
}

func NewOrderServer(localHost, port string, hostArray []string) (server *ServerOrder) {
	hashConsistent := distributed.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	server = &ServerOrder{
		rabbitMq:          rabbitmq.NewRabbitMQSimple("iRaidenProduct"),
		localHost:         localHost,
		port:              port,
		accessControlUnit: distributed.NewAccessControlUnit(),
		hashConsistent:    hashConsistent,
	}
	return
}
