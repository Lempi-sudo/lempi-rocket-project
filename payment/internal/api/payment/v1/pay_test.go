package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	modelError "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/model"
	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
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

func (s *APISuite) TestPaymentServiceError() {
	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD
		serviceErr    = gofakeit.Error()
	)

	req := &paymentV1.PayOrderRequest{
		Order: &paymentV1.OrderInfo{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		},
	}
	s.paymentService.On("Pay").Return("", serviceErr)

	response, err := s.api.PayOrder(s.ctx, req)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Error(err)

	s.Require().Nil(response)
}

func (s *APISuite) TestPayEmptyUUID() {
	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD
		errEmptyUUID  = modelError.ErrEmptyUUID
	)

	req := &paymentV1.PayOrderRequest{
		Order: &paymentV1.OrderInfo{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		},
	}

	s.paymentService.On("Pay").Return("", errEmptyUUID)
	response, err := s.api.PayOrder(s.ctx, req)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Error(err)

	s.Require().Nil(response)
}
