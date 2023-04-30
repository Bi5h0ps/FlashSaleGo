package main

import (
	"FlashSaleGo/grpc/inventory"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}

	s := inventory.NewServerInventory(1000)

	grpcServer := grpc.NewServer()
	inventory.RegisterInventoryServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000:%v", err)
	}
}
