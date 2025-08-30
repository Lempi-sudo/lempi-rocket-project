package part

import (
	"time"

	repoModel "github.com/Lempi-sudo/lempi-rocket-project/inventory/internal/repository/model"
	"github.com/samber/lo"
)

func NewRepository() *repository {
	parts := map[string]repoModel.Part{
		"550e8400-e29b-41d4-a716-446655440000": {
			Uuid: "550e8400-e29b-41d4-a716-446655440000",
			Info: repoModel.PartInfo{
				Name:          "Main Engine",
				Description:   "Primary propulsion engine",
				Price:         1000000.0,
				StockQuantity: 5,
				Category:      repoModel.CategoryEngine,
				Dimensions:    repoModel.NewDimensions(200.0, 100.0, 100.0, 5000.0),
				Manufacturer:  repoModel.NewManufacturer("RocketCorp", "USA", "https://rocketcorp.com"),
				Tags:          []string{"main", "engine", "propulsion"},
				Metadata: map[string]*repoModel.Value{
					"serial_number": {String_value: lo.ToPtr("SN-001")},
					"max_thrust":    {Float64_value: lo.ToPtr(1500.0)},
				},
				UpdatedAt: lo.ToPtr(time.Now()),
				CreatedAt: lo.ToPtr(time.Now()),
			},
		},
		"550e8400-e29b-41d4-a716-446655440001": {
			Uuid: "550e8400-e29b-41d4-a716-446655440001",
			Info: repoModel.PartInfo{
				Name:          "Porthole",
				Description:   "Window for space view",
				Price:         50000.0,
				StockQuantity: 20,
				Category:      repoModel.CategoryPorthole,
				Dimensions:    repoModel.NewDimensions(50.0, 50.0, 5.0, 10.0),
				Manufacturer:  repoModel.NewManufacturer("SpaceGlass", "Germany", "https://spaceglass.de"),
				Tags:          []string{"window", "main", "glass", "porthole"},
				Metadata: map[string]*repoModel.Value{
					"tint": {String_value: lo.ToPtr("UV-protect")},
				},
				UpdatedAt: lo.ToPtr(time.Now()),
				CreatedAt: lo.ToPtr(time.Now()),
			},
		},
		"550e8400-e29b-41d4-a716-446655440003": {
			Uuid: "550e8400-e29b-41d4-a716-446655440003",
			Info: repoModel.PartInfo{
				Name:          "Fuel Tank",
				Description:   "Liquid fuel storage tank",
				Price:         250000.0,
				StockQuantity: 10,
				Category:      repoModel.CategoryFuel,
				Dimensions:    repoModel.NewDimensions(300.0, 150.0, 150.0, 2000.0),
				Manufacturer:  repoModel.NewManufacturer("FuelTech", "Germany", "https://fueltech.de"),
				Tags:          []string{"fuel", "tank", "main", "storage"},
				Metadata: map[string]*repoModel.Value{
					"capacity": {Float64_value: lo.ToPtr(23.4)},
					"material": {String_value: lo.ToPtr("Titanium")},
				},
				UpdatedAt: lo.ToPtr(time.Now()),
				CreatedAt: lo.ToPtr(time.Now()),
			},
		},
		"550e8400-e29b-41d4-a716-446655440004": {
			Uuid: "550e8400-e29b-41d4-a716-446655440004",
			Info: repoModel.PartInfo{
				Name:          "Wing Panel",
				Description:   "Aerodynamic wing component",
				Price:         75000.0,
				StockQuantity: 15,
				Category:      repoModel.CategoryWing,
				Dimensions:    repoModel.NewDimensions(400.0, 200.0, 20.0, 800.0),
				Manufacturer:  repoModel.NewManufacturer("AeroDynamics", "USA", "https://aerodynamics.com"),
				Tags:          []string{"wing", "panel", "main", "aerodynamic"},
				Metadata: map[string]*repoModel.Value{
					"airfoil_type": {String_value: lo.ToPtr("NACA 2412")},
					"span":         {Float64_value: lo.ToPtr(400.0)},
				},
				UpdatedAt: lo.ToPtr(time.Now()),
				CreatedAt: lo.ToPtr(time.Now()),
			},
		},
	}
	return &repository{data: parts}
}
