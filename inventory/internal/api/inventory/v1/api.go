package v1

import (
	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/service"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

type inventoryApi struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *inventoryApi {
	return &inventoryApi{
		service: inventoryService,
	}
}
