package v1

import (
	"context"
	"log"

	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PayOrder обрабатывает запрос на оплату заказа.
//
// Проверяет корректность UUID заказа и пользователя, а также наличие метода оплаты.
// В случае успешной оплаты возвращает сгенерированный UUID транзакции.
func (p *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
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

	paymentUUID := p.serverPayment.Pay()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", paymentUUID)
	return &paymentV1.PayOrderResponse{
		Uuid: paymentUUID,
	}, nil
}
