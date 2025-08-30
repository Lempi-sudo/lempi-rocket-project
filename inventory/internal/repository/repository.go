package repository

import (
	"context"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	GetAllParts(ctx context.Context) ([]model.Part, error)
}
