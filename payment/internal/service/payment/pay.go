package payment

import "github.com/google/uuid"

func (p *service) Pay() string {
	return uuid.NewString()
}
