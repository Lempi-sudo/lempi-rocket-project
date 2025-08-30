package part

import (
	"sync"

	def "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository"
	repoModel "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}
