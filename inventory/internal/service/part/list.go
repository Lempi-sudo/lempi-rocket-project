package part

import (
	"context"
	"log"

	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/model"
	"github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/service/converter"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	repoParts, err := s.inventoryRepository.GetAllParts(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		parts := make([]model.Part, 0, len(repoParts))
		for _, part := range repoParts {
			parts = append(parts, converter.PartRepositoryToModel(part))
		}
		return parts, nil
	}

	filteredParts := make([]model.Part, 0, len(repoParts))
	for _, part := range repoParts {
		filteredParts = append(filteredParts, converter.PartRepositoryToModel(part))
	}

	if len(filter.Uuids) > 0 {
		uuidSet := make(map[string]bool)
		for _, uuid := range filter.Uuids {
			uuidSet[uuid] = true
		}

		tempParts := make([]model.Part, 0)
		for _, part := range filteredParts {
			if uuidSet[part.Uuid] {
				tempParts = append(tempParts, part)
			}
		}
		filteredParts = tempParts
	}

	if len(filter.Names) > 0 {
		nameSet := make(map[string]bool)
		for _, name := range filter.Names {
			nameSet[name] = true
		}

		tempParts := make([]model.Part, 0)
		for _, part := range filteredParts {
			if nameSet[part.Info.Name] {
				tempParts = append(tempParts, part)
			}
		}
		filteredParts = tempParts
	}

	if len(filter.Categories) > 0 {
		categorySet := make(map[model.Category]bool)
		for _, category := range filter.Categories {
			categorySet[category] = true
		}

		tempParts := make([]model.Part, 0)
		for _, part := range filteredParts {
			if categorySet[part.Info.Category] {
				tempParts = append(tempParts, part)
			}
		}
		filteredParts = tempParts
	}

	if len(filter.ManufacturerCountries) > 0 {
		countrySet := make(map[string]bool)
		for _, country := range filter.ManufacturerCountries {
			countrySet[country] = true
		}

		tempParts := make([]model.Part, 0)
		for _, part := range filteredParts {
			if part.Info.Manufacturer != nil && countrySet[part.Info.Manufacturer.Country] {
				tempParts = append(tempParts, part)
			}
		}
		filteredParts = tempParts
	}

	if len(filter.Tags) > 0 {
		tagSet := make(map[string]bool)
		for _, tag := range filter.Tags {
			tagSet[tag] = true
		}

		tempParts := make([]model.Part, 0)
		for _, part := range filteredParts {
			for _, partTag := range part.Info.Tags {
				if tagSet[partTag] {
					tempParts = append(tempParts, part)
					break
				}
			}
		}
		filteredParts = tempParts
	}
	log.Println("Inventory return order following by filter rules")
	return filteredParts, nil
}
