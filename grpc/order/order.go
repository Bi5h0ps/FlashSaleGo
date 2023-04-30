package order

import (
	"FlashSaleGo/distributed"
	"FlashSaleGo/grpc/inventory"
	"FlashSaleGo/model"
	"FlashSaleGo/rabbitmq"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerOrder struct {
	// RabbitMQ
	rabbitMq *rabbitmq.RabbitMQ
	// LocalHost
	localHost string
	// Inventory Server
	inventoryIP string
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
	right := s.accessControlUnit.GetDistributedRight(message.UserID, s.hashConsistent, s.localHost)
	if !right {
		return &OrderResult{IsOrderSuccess: "false"},
			errors.New("distributed user auth failed")
	}
	//2.product number control
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.inventoryIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Fail to dial on inventory server, err: %v", err))
	}
	defer conn.Close()
	c := inventory.NewInventoryServiceClient(conn)
	inventoryControlResult, err := c.UpdateProductCount(ctx,
		&inventory.ProductInfo{ProductID: message.ProductID})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Inventory Server Error, err: %v", err))
	}
	if !inventoryControlResult.IsInventorySuccess {
		return &OrderResult{IsOrderSuccess: "false"}, nil
	}
	//3.write to rabbitmq
	rabbitmqMessage := model.NewMessage(message.UserID, message.ProductID)
	byteMessage, err := json.Marshal(rabbitmqMessage)
	if err != nil {
		return nil, errors.New("rabbitmq message type conversion failed")
	}
	err = s.rabbitMq.PublishSimple(string(byteMessage))
	if err != nil {
		return nil, errors.New("rabbitmq failed to publish")
	}
	return &OrderResult{IsOrderSuccess: "true"}, nil
}

func (s *ServerOrder) GetUserCloudData(ctx context.Context, message *UserInfo) (*UserCloudData, error) {
	data, err := s.accessControlUnit.GetDataFromMap(message.UserID)
	return &UserCloudData{TimeStamp: data.LastOrderTime}, err
}

func (s *ServerOrder) Destroy() {
	s.rabbitMq.Destroy()
}

func NewOrderServer(localHost, inventoryIP string, hostArray []string) (server *ServerOrder) {
	hashConsistent := distributed.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	server = &ServerOrder{
		rabbitMq:          rabbitmq.NewRabbitMQSimple("iRaidenProduct"),
		localHost:         localHost,
		inventoryIP:       inventoryIP,
		accessControlUnit: distributed.NewAccessControlUnit(),
		hashConsistent:    hashConsistent,
	}
	return
}
