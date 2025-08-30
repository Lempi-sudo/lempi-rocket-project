package converter

import (
	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	repoModel "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
)

func PartRepositoryToModel(sighting repoModel.Part) model.Part {
	return model.Part{
		Uuid: sighting.Uuid,
		Info: partInfoToModel(sighting.Info),
	}
}

func partInfoToModel(info repoModel.PartInfo) model.PartInfo {
	// copy tags
	var tags []string
	if len(info.Tags) > 0 {
		tags = make([]string, len(info.Tags))
		copy(tags, info.Tags)
	}

	// convert metadata
	var metadata map[string]*model.Value
	if info.Metadata != nil {
		metadata = make(map[string]*model.Value, len(info.Metadata))
		for k, v := range info.Metadata {
			metadata[k] = valueToModel(v)
		}
	}

	d := model.NewDimensions(info.Dimensions.Length, info.Dimensions.Width, info.Dimensions.Height, info.Dimensions.Weight)
	m := model.NewManufacturer(info.Manufacturer.Name, info.Manufacturer.Country, info.Manufacturer.Website)

	return model.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Tags:          tags,
		Category:      model.Category(info.Category),
		Dimensions:    d,
		Manufacturer:  m,
		Metadata:      metadata,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

func valueToModel(v *repoModel.Value) *model.Value {
	if v == nil {
		return nil
	}
	return &model.Value{
		String_value:  v.String_value,
		Int64_value:   v.Int64_value,
		Float64_value: v.Float64_value,
		Bool_value:    v.Bool_value,
	}
}
