package payment

import (
	modelError "github.com/Lempi-sudo/lempi-rocket-project/payment/internal/model"
	"github.com/google/uuid"
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
