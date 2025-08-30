package main

import (
	// 1. Стандартные библиотеки Go.

	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// 2. Сторонние библиотеки (начинаются с домена).
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/api/inventory/v1"
	repository "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/part"
	serviceInventory "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/service/part"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	repo := repository.NewRepository()
	service := serviceInventory.NewService(repo)
	apiInv := inventoryApi.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, apiInv)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
