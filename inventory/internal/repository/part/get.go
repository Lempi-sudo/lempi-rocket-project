package part

import (
	"context"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	repoModel "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
)

// GetPart возвращает информацию о детали по её UUID.
//
// Если UUID отсутствует в запросе или не найден в хранилище, возвращается ошибка.
func (r *repository) GetPart(ctx context.Context, uuid string) (repoModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(uuid) == 0 {
		return repoModel.Part{}, model.ErrBadUuid
	}

	repoPart, ok := r.data[uuid]
	if !ok {
		return repoModel.Part{}, model.ErrPartNotFound
	}

	return repoPart, nil
}
