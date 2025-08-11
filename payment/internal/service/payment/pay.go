package payment

import (
	"github.com/google/uuid"

	modelError "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/model"
)

func (p *service) Pay() (string, error) {
	v4, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuid := v4.String()
	if len(uuid) == 0 {
		return "", modelError.ErrEmptyUUID
	}
	return uuid, nil
}
