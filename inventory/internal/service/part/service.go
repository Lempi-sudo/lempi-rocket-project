package part

import (
	repoModel "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository"
	def "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repoModel.InventoryRepository
}

func NewService(ufoRepository repoModel.InventoryRepository) *service {
	return &service{
		inventoryRepository: ufoRepository,
	}
}
