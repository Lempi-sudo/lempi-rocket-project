package v1

import (
	"github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service"
	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
)

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer

	serverPayment service.PaymentService
}

func NewPaymentService(server service.PaymentService) *paymentService {
	return &paymentService{
		serverPayment: server,
	}
}
