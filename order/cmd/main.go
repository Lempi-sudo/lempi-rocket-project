package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/payment/v1"
)

const (
	orderHttpPort    = "8080"
	inventoryAddress = "localhost:50052"
	paymentAddress   = "localhost:50051"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var (
	ErrOrderAlreadyPaid error = errors.New("order has already paid")
	ErrOrderNotFound    error = errors.New("order not found")
)

type OrderStorage struct {
	orders map[string]*orderV1.OrderDto
	mu     sync.RWMutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) UpdateOrder(uuid string, order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = order
}

func (s *OrderStorage) CancelOrder(uuid string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.orders[uuid]; !ok {
		return ErrOrderNotFound
	}
	if s.orders[uuid].Status == orderV1.OrderStatusPAID {
		return ErrOrderAlreadyPaid
	}
	if s.orders[uuid].Status == orderV1.OrderStatusPENDINGPAYMENT {
		s.orders[uuid].Status = orderV1.OrderStatusCANCELLED
	}
	return nil
}

type OrderHandler struct {
	orderV1.UnimplementedHandler
	orderStorage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		orderStorage: storage,
	}
}

func (h *OrderHandler) CreateOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return &orderV1.CreateOrderInternalServerError{
			Message: "Internal server error",
			Code:    500,
		}, nil
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент
	client := inventoryV1.NewInventoryServiceClient(conn)

	// Преобразуем UUID в строки для запроса
	var partUuids []string
	for _, uuid := range req.GetPartUuids() {
		partUuids = append(partUuids, uuid.String())
	}

	// Создаем запрос к inventory сервису
	listPartsReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUuids,
		},
	}

	// Отправляем запрос
	response, err := client.ListParts(ctx, listPartsReq)
	if err != nil {
		log.Printf("Error calling inventory service: %v", err)
		return &orderV1.CreateOrderInternalServerError{
			Message: "Failed to get parts from inventory",
			Code:    500,
		}, nil
	}

	// Проверяем, что все запрошенные детали найдены
	if len(response.Parts) != len(partUuids) {
		return &orderV1.CreateOrderBadRequest{
			Message: "Some parts not found in inventory",
			Code:    400,
		}, nil
	}

	totalPrice := 0.0
	for _, part := range response.GetParts() {
		totalPrice += part.Price
	}

	orderUUID, _ := uuid.NewRandom()
	h.orderStorage.UpdateOrder(orderUUID.String(),
		&orderV1.OrderDto{
			OrderUUID:  orderUUID,
			TotalPrice: totalPrice,
			UserUUID:   req.GetUserUUID(),
			PartUuids:  req.GetPartUuids(),
			Status:     orderV1.OrderStatusPENDINGPAYMENT,
		})
	return &orderV1.CreateOrderResponse{
		TotalPrice: totalPrice,
		OrderUUID:  orderUUID,
	}, nil
}

func (h *OrderHandler) PayOrder(_ context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Message: "Internal server error",
			Code:    500,
		}, nil
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	orderUUUID := params.OrderUUID.String()
	h.orderStorage.mu.RLock()
	// defer h.orderStorage.mu.RUnlock() если делать так, то лочится updateOrder
	order := h.orderStorage.GetOrder(orderUUUID)
	if order == nil {
		h.orderStorage.mu.RUnlock()
		return &orderV1.PayOrderNotFound{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	h.orderStorage.mu.RUnlock()

	var paymentMethod paymentV1.PaymentMethod
	switch req.GetPaymentMethod() {
	case orderV1.PaymentMethodCARD:
		paymentMethod = paymentV1.PaymentMethod_CARD
	case orderV1.PaymentMethodCREDITCARD:
		paymentMethod = paymentV1.PaymentMethod_CREDIT_CARD
	case orderV1.PaymentMethodSBP:
		paymentMethod = paymentV1.PaymentMethod_SBP
	case orderV1.PaymentMethodINVESTORMONEY:
		paymentMethod = paymentV1.PaymentMethod_INVESTOR_MONEY
	case orderV1.PaymentMethodUNKNOWN:
		paymentMethod = paymentV1.PaymentMethod_UNKNOWN_UNSPECIFIED
	default:
		return &orderV1.PayOrderBadRequest{
			Code:    400,
			Message: "Order not found",
		}, nil
	}

	userUUID := order.GetUserUUID().String()

	// Создаем gRPC клиент
	client := paymentV1.NewPaymentServiceClient(conn)

	payOrderRequest := &paymentV1.PayOrderRequest{
		Order: &paymentV1.OrderInfo{
			OrderUuid:     orderUUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		},
	}

	payOrderResponse, err := client.PayOrder(ctx, payOrderRequest)
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Message: "Internal server error",
			Code:    500,
		}, nil
	}

	transactionUUIDStr := payOrderResponse.GetUuid()

	// Преобразуем строку в UUID
	transactionUUID, err := uuid.Parse(transactionUUIDStr)
	if err != nil {
		return &orderV1.PayOrderInternalServerError{
			Message: "Invalid transaction UUID",
			Code:    500,
		}, nil
	}

	// Создаем новый заказ с обновленными данными
	newOrder := &orderV1.OrderDto{
		OrderUUID:       order.GetOrderUUID(),
		TotalPrice:      order.GetTotalPrice(),
		UserUUID:        order.GetUserUUID(),
		PartUuids:       order.GetPartUuids(),
		Status:          orderV1.OrderStatusPAID,
		TransactionUUID: orderV1.OptNilUUID{Set: true, Value: transactionUUID},
		PaymentMethod:   orderV1.OptPaymentMethod{Set: true, Value: req.GetPaymentMethod()},
	}

	// Обновляем заказ в хранилище
	h.orderStorage.UpdateOrder(orderUUUID, newOrder)

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (h *OrderHandler) CancelOrderByUUID(_ context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	// Преобразуем UUID в строку
	orderUUIDStr := params.OrderUUID.String()

	err := h.orderStorage.CancelOrder(orderUUIDStr)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.CancelOrderByUUIDNotFound{
				Message: "Order not found",
				Code:    404,
			}, nil
		}
		if errors.Is(err, ErrOrderAlreadyPaid) {
			return &orderV1.CancelOrderByUUIDConflict{
				Message: "Order has already paid",
				Code:    409,
			}, nil
		}
		// Возвращаем ошибку 500 - внутренняя ошибка сервера
		return &orderV1.CancelOrderByUUIDInternalServerError{
			Message: "Internal server error",
			Code:    500,
		}, nil
	}
	// Возвращаем успешный ответ 204 - заказ отменен
	return &orderV1.CancelOrderByUUIDNoContent{}, nil
}

func (h *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	uuid := params.OrderUUID.String()
	h.orderStorage.mu.RLock()
	defer h.orderStorage.mu.RUnlock()
	order := h.orderStorage.GetOrder(uuid)
	if order == nil {
		return &orderV1.GetOrderByUUIDNotFound{
			Code:    404,
			Message: "Order not found",
		}, nil
	}
	return order, nil
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", orderHttpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", orderHttpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
