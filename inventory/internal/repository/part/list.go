package part

import (
	"context"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
)

func (r *repository) GetAllParts(ctx context.Context) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, part := range r.data {
		parts = append(parts, part)
	}

	return parts, nil
}
