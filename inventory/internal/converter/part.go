package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

func ProtoFilterToModel(protoFilter *inventoryV1.PartsFilter) model.PartsFilter {
	var tags []string
	if len(protoFilter.Tags) > 0 {
		tags = make([]string, len(protoFilter.Tags))
		copy(tags, protoFilter.Tags)
	}

	var uuids []string
	if len(protoFilter.Uuids) > 0 {
		uuids = make([]string, len(protoFilter.Uuids))
		copy(uuids, protoFilter.Uuids)
	}

	var categories []model.Category
	if len(protoFilter.Categories) > 0 {
		categories = make([]model.Category, len(protoFilter.Categories))
		for i, protoCategory := range protoFilter.Categories {
			categories[i] = model.Category(protoCategory)
		}
	}

	var names []string
	if len(protoFilter.Names) > 0 {
		names = make([]string, len(protoFilter.Names))
		copy(names, protoFilter.Names)
	}

	var manufacturerCountries []string
	if len(protoFilter.ManufacturerCountries) > 0 {
		manufacturerCountries = make([]string, len(protoFilter.ManufacturerCountries))
		copy(manufacturerCountries, protoFilter.ManufacturerCountries)
	}

	filter := model.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
	return filter
}

func ValueToProto(value *model.Value) *inventoryV1.Value {
	if value == nil {
		return nil
	}

	protoValue := &inventoryV1.Value{}

	if value.String_value != nil {
		protoValue.Kind = &inventoryV1.Value_StringValue{
			StringValue: *value.String_value,
		}
	} else if value.Int64_value != nil {
		protoValue.Kind = &inventoryV1.Value_Int64Value{
			Int64Value: *value.Int64_value,
		}
	} else if value.Float64_value != nil {
		protoValue.Kind = &inventoryV1.Value_DoubleValue{
			DoubleValue: *value.Float64_value,
		}
	} else if value.Bool_value != nil {
		protoValue.Kind = &inventoryV1.Value_BoolValue{
			BoolValue: *value.Bool_value,
		}
	}

	return protoValue
}

func PartToProto(part model.Part) *inventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.Info.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.Info.UpdatedAt)
	}

	var createdAt *timestamppb.Timestamp
	if part.Info.CreatedAt != nil {
		createdAt = timestamppb.New(*part.Info.CreatedAt)
	}

	var tags []string
	if len(part.Info.Tags) > 0 {
		tags = make([]string, len(part.Info.Tags))
		copy(tags, part.Info.Tags)
	}

	d := &inventoryV1.Dimensions{
		Length: part.Info.Dimensions.Length,
		Width:  part.Info.Dimensions.Width,
		Height: part.Info.Dimensions.Height,
		Weight: part.Info.Dimensions.Weight,
	}

	m := &inventoryV1.Manufacturer{
		Name:    part.Info.Manufacturer.Name,
		Country: part.Info.Manufacturer.Country,
		Website: part.Info.Manufacturer.Website,
	}

	var metadata map[string]*inventoryV1.Value
	if len(part.Info.Metadata) > 0 {
		metadata = make(map[string]*inventoryV1.Value, len(part.Info.Metadata))
		for key, value := range part.Info.Metadata {
			metadata[key] = ValueToProto(value)
		}
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Info.Name,
		Description:   part.Info.Description,
		Price:         part.Info.Price,
		StockQuantity: part.Info.StockQuantity,
		Category:      inventoryV1.Category(part.Info.Category),
		Tags:          tags,
		Dimensions:    d,
		Manufacturer:  m,
		Metadata:      metadata,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
