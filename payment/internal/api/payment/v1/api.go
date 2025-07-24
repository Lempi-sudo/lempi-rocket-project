package v1

import paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *paymentService {
	return &paymentService{}
}
