package main

import (
	"FlashSaleGo/common"
	"FlashSaleGo/grpc/order"
	"google.golang.org/grpc"
	"log"
	"net"
)

var LocalHost = ""
var Port = "9093"

func main() {
	lis, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", Port, err)
	}

	LocalHost, err = common.GetIntranceIp()
	if err != nil {
		log.Fatal(err)
	}

	s := order.NewOrderServer(LocalHost, Port)
	defer s.Destroy()

	grpcServer := grpc.NewServer()

	order.RegisterOrderServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000:%v", err)
	}
}
