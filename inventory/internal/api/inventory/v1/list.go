package v1

import (
	"context"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/converter"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

func (a *inventoryApi) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ProtoFilterToModel(req.GetFilter())
	parts, err := a.service.ListParts(ctx, &filter)
	if err != nil {
		return nil, err
	}

	partsResult := make([]*inventoryV1.Part, len(parts))

	for ind, part := range parts {
		partsResult[ind] = converter.PartToProto(part)
	}

	return &inventoryV1.ListPartsResponse{Parts: partsResult}, nil
}
