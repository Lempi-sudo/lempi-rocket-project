package v1

import (
	"github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service"
	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
)

type paymentAPI struct {
	paymentV1.UnimplementedPaymentServiceServer

	serverPayment service.PaymentService
}

func NewPaymentAPI(server service.PaymentService) *paymentAPI {
	return &paymentAPI{
		serverPayment: server,
	}
}
