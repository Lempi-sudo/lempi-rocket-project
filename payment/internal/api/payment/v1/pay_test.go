package v1

import (
	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APISuite) TestPaySuccess() {
	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD
		uuid          = gofakeit.UUID()
	)

	req := &paymentV1.PayOrderRequest{
		Order: &paymentV1.OrderInfo{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		},
	}

	expectedPayOrderResponse := &paymentV1.PayOrderResponse{
		Uuid: uuid,
	}

	s.paymentService.On("Pay").Return(uuid, nil)

	response, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().Equal(response.GetUuid(), expectedPayOrderResponse.Uuid)

}
