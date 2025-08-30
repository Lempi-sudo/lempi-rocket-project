package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/converter"
	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

func (a *inventoryApi) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.service.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrBadUuid) {
			return nil, status.Errorf(codes.InvalidArgument, "Bad uuid")
		}
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, err

	}
	return &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}, nil
}
