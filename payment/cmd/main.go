package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (p *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Генерируем UUID для нового наблюдения
	orderUuid := req.GetOrder().OrderUuid
	if len(orderUuid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bad uuid")
	}

	payment_method := req.GetOrder().PaymentMethod
	if payment_method == paymentV1.PaymentMethod_UNKNOWN_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "Payment method unspecified ")
	}

	userUuid := req.GetOrder().UserUuid
	if len(userUuid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bad uuid")
	}

	paymentUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", paymentUUID)
	return &paymentV1.PayOrderResponse{
		Uuid: paymentUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listen: %v\n", err)
		}
	}()

	s := grpc.NewServer()

	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
