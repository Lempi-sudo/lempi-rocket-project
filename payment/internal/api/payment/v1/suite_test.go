package v1

import (
	"context"
	"testing"

	"github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite
	ctx            context.Context
	paymentService *mocks.PaymentService
	api            *paymentAPI
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewPaymentService(s.T())

	s.api = NewPaymentAPI(s.paymentService)
}

func (s *APISuite) TearDownTest() {

}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
