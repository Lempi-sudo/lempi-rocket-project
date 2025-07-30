package payment

import def "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
