package part

import (
	"context"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/service/converter"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	p, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return converter.PartRepositoryToModel(p), nil
}
