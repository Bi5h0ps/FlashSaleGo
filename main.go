package main

import (
	"FlashSaleGo/common"
	"FlashSaleGo/grpc/order"
	"google.golang.org/grpc"
	"log"
	"net"
)

var localHost = ""
var inventoryIP = ":9094"
var localHostPort = ":9093"

// distributed machine's ip address
var hostArray = []string{"127.0.0.1:9093"}

func main() {
	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatalf("Failed to listen on port 9093: %v", err)
	}

	//get localHost address
	localHost, err = common.GetIntranceIp()
	if err != nil {
		log.Fatal(err)
	}

	s := order.NewOrderServer(localHost+localHostPort, inventoryIP, hostArray)
	defer s.Destroy()

	grpcServer := grpc.NewServer()

	order.RegisterOrderServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9093:%v", err)
	}
}
