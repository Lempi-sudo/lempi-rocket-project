package payment

import (
	"testing"

	s "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("creates new service instance", func(t *testing.T) {
		// Вызываем конструктор
		svc := NewService()

		// Проверяем, что возвращается не nil
		require.NotNil(t, svc)

		// Проверяем тип возвращаемого значения
		assert.IsType(t, &service{}, svc)
	})

	t.Run("service implements PaymentService interface", func(t *testing.T) {
		// Создаем экземпляр сервиса
		svc := NewService()
		require.NotNil(t, svc)

		// Проверяем, что сервис реализует интерфейс PaymentService
		var _ s.PaymentService = svc
	})
}
