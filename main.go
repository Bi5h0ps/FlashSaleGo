package main

import (
	"FlashSaleGo/common"
	"FlashSaleGo/grpc/order"
	"google.golang.org/grpc"
	"log"
	"net"
)

var localHost = ""
var port = "9093"

// distributed machine's ip address
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

func main() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}

	localHost, err = common.GetIntranceIp()
	if err != nil {
		log.Fatal(err)
	}

	s := order.NewOrderServer(localHost, port, hostArray)
	defer s.Destroy()

	grpcServer := grpc.NewServer()

	order.RegisterOrderServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %v:%v", port, err)
	}
}
