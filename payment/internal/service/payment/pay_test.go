package payment

import (
	"errors"
	"testing"

	modelError "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/model"
	"github.com/Lempi-sudo/lempi-rocket-project/payment/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicePay(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful payment with valid UUID",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{}
			got, err := s.Pay()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Len(t, got, 36) // UUID v4 имеет длину 36 символов
			}
		})
	}
}

func TestServicePayReturnsValidUUID(t *testing.T) {
	s := &service{}

	// Вызываем функцию несколько раз для проверки уникальности
	uuids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		uuid, err := s.Pay()
		require.NoError(t, err)
		require.NotEmpty(t, uuid)
		require.Len(t, uuid, 36)

		// Проверяем, что UUID уникален
		assert.False(t, uuids[uuid], "UUID должен быть уникальным: %s", uuid)
		uuids[uuid] = true

		// Проверяем формат UUID v4 (8-4-4-4-12 символов)
		assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, uuid)
	}
}

func TestServicePayNoEmptyUUID(t *testing.T) {
	s := &service{}

	// Вызываем функцию много раз, чтобы убедиться, что никогда не возвращается пустой UUID
	for i := 0; i < 1000; i++ {
		uuid, err := s.Pay()
		require.NoError(t, err)
		require.NotEmpty(t, uuid)
		require.Greater(t, len(uuid), 0)
	}
}

func TestServicePayConsistentBehavior(t *testing.T) {
	s := &service{}

	// Проверяем, что функция работает консистентно
	uuid1, err1 := s.Pay()
	uuid2, err2 := s.Pay()

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotEmpty(t, uuid1)
	require.NotEmpty(t, uuid2)

	// UUID должны быть разными
	assert.NotEqual(t, uuid1, uuid2)

	// Но должны иметь одинаковую длину
	assert.Equal(t, len(uuid1), len(uuid2))
	assert.Equal(t, 36, len(uuid1))
}

// Тесты с использованием моков для проверки различных сценариев ошибок
func TestPaymentServiceWithMocks(t *testing.T) {
	t.Run("mock successful payment", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
		mockService.EXPECT().Pay().Return(expectedUUID, nil)

		// Используем мок как интерфейс
		uuid, err := mockService.Pay()

		assert.NoError(t, err)
		assert.Equal(t, expectedUUID, uuid)
		mockService.AssertExpectations(t)
	})

	t.Run("mock payment with error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		expectedError := errors.New("failed to generate UUID")
		mockService.EXPECT().Pay().Return("", expectedError)

		uuid, err := mockService.Pay()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})

	t.Run("mock payment with empty UUID error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		// Симулируем случай, когда UUID генерируется, но оказывается пустым
		expectedError := errors.New("generated uuid is empty")
		mockService.EXPECT().Pay().Return("", expectedError)

		uuid, err := mockService.Pay()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})

	t.Run("mock payment with custom error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		customError := errors.New("database connection failed")
		mockService.EXPECT().Pay().Return("", customError)

		uuid, err := mockService.Pay()

		assert.Error(t, err)
		assert.Equal(t, customError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})
}

// Тест для проверки интеграции с реальным сервисом через интерфейс
func TestPaymentServiceInterface(t *testing.T) {
	t.Run("real service implements interface", func(t *testing.T) {
		realService := NewService()

		// Проверяем, что реальный сервис реализует интерфейс
		uuid, err := realService.Pay()
		assert.NoError(t, err)
		assert.NotEmpty(t, uuid)
		assert.Len(t, uuid, 36)
	})
}

// Тест для проверки поведения мока при множественных вызовах
func TestPaymentServiceMockMultipleCalls(t *testing.T) {
	mockService := mocks.NewPaymentService(t)

	// Настраиваем мок для возврата разных значений при каждом вызове
	mockService.EXPECT().Pay().Return("uuid-1", nil).Once()
	mockService.EXPECT().Pay().Return("uuid-2", nil).Once()
	mockService.EXPECT().Pay().Return("", errors.New("error on third call")).Once()

	// Первый вызов - успех
	uuid1, err1 := mockService.Pay()
	assert.NoError(t, err1)
	assert.Equal(t, "uuid-1", uuid1)

	// Второй вызов - успех
	uuid2, err2 := mockService.Pay()
	assert.NoError(t, err2)
	assert.Equal(t, "uuid-2", uuid2)

	// Третий вызов - ошибка
	uuid3, err3 := mockService.Pay()
	assert.Error(t, err3)
	assert.Empty(t, uuid3)

	mockService.AssertExpectations(t)
}

// Тест для проверки различных типов ошибок через моки
func TestPaymentServiceMockErrorScenarios(t *testing.T) {

	t.Run("empty uuid", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		errEmptyUUID := modelError.ErrEmptyUUID
		mockService.EXPECT().Pay().Return("", errEmptyUUID)

		uuid, err := mockService.Pay()
		assert.Error(t, err)
		assert.Equal(t, errEmptyUUID, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})
	t.Run("mock network error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		networkError := errors.New("network timeout")
		mockService.EXPECT().Pay().Return("", networkError)

		uuid, err := mockService.Pay()
		assert.Error(t, err)
		assert.Equal(t, networkError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})

	t.Run("mock system error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		systemError := errors.New("system resources exhausted")
		mockService.EXPECT().Pay().Return("", systemError)

		uuid, err := mockService.Pay()
		assert.Error(t, err)
		assert.Equal(t, systemError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})

	t.Run("mock validation error", func(t *testing.T) {
		mockService := mocks.NewPaymentService(t)

		validationError := errors.New("invalid payment method")
		mockService.EXPECT().Pay().Return("", validationError)

		uuid, err := mockService.Pay()
		assert.Error(t, err)
		assert.Equal(t, validationError, err)
		assert.Empty(t, uuid)
		mockService.AssertExpectations(t)
	})
}
