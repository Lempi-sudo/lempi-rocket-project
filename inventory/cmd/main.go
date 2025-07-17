package main

import (
	// 1. Стандартные библиотеки Go.
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// 2. Сторонние библиотеки (начинаются с домена).
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	// 3. Локальные импорты.
	inventoryV1 "github.com/Lempi-sudo/lempi-rocket-project/shared/pkg/proto/inventory/v1"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	service := NewInventoryService()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

const grpcPort = 50052

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	parts map[string]*inventoryV1.Part
	mu    sync.RWMutex
}

// NewInventoryService returns inventoryService instanse.
func NewInventoryService() *inventoryService {
	parts := map[string]*inventoryV1.Part{
		"550e8400-e29b-41d4-a716-446655440000": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440000",
			Name:          "Main Engine",
			Description:   "Primary propulsion engine",
			Price:         1000000.0,
			StockQuantity: 5,
			Category:      inventoryV1.Category_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 200.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 5000.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "RocketCorp",
				Country: "USA",
				Website: "https://rocketcorp.com",
			},
			Tags: []string{"main", "engine", "propulsion"},
			Metadata: map[string]*inventoryV1.Value{
				"serial_number": {Kind: &inventoryV1.Value_StringValue{StringValue: "SN-001"}},
				"max_thrust":    {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 1500.0}},
			},
		},
		"550e8400-e29b-41d4-a716-446655440001": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Porthole",
			Description:   "Window for space view",
			Price:         50000.0,
			StockQuantity: 20,
			Category:      inventoryV1.Category_PORTHOLE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 5.0,
				Weight: 10.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "SpaceGlass",
				Country: "Germany",
				Website: "https://spaceglass.de",
			},
			Tags: []string{"window", "glass", "porthole"},
			Metadata: map[string]*inventoryV1.Value{
				"tint": {Kind: &inventoryV1.Value_StringValue{StringValue: "UV-protect"}},
			},
		},
		"550e8400-e29b-41d4-a716-446655440003": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Fuel Tank",
			Description:   "Liquid fuel storage tank",
			Price:         250000.0,
			StockQuantity: 10,
			Category:      inventoryV1.Category_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 300.0,
				Width:  150.0,
				Height: 150.0,
				Weight: 2000.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "FuelTech",
				Country: "Germany",
				Website: "https://fueltech.de",
			},
			Tags: []string{"fuel", "tank", "storage"},
			Metadata: map[string]*inventoryV1.Value{
				"capacity": {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 5000.0}},
				"material": {Kind: &inventoryV1.Value_StringValue{StringValue: "Titanium"}},
			},
		},
		"550e8400-e29b-41d4-a716-446655440004": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Wing Panel",
			Description:   "Aerodynamic wing component",
			Price:         75000.0,
			StockQuantity: 15,
			Category:      inventoryV1.Category_WING,
			Dimensions: &inventoryV1.Dimensions{
				Length: 400.0,
				Width:  200.0,
				Height: 20.0,
				Weight: 800.0,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "AeroDynamics",
				Country: "USA",
				Website: "https://aerodynamics.com",
			},
			Tags: []string{"wing", "panel", "aerodynamic"},
			Metadata: map[string]*inventoryV1.Value{
				"airfoil_type": {Kind: &inventoryV1.Value_StringValue{StringValue: "NACA 2412"}},
				"span":         {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 400.0}},
			},
		},
	}
	return &inventoryService{parts: parts}
}

// GetPart возвращает информацию о детали по её UUID.
//
// Если UUID отсутствует в запросе или не найден в хранилище, возвращается ошибка с соответствующим кодом gRPC.
func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uid := req.GetUuid()
	if len(uid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bad uuid")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.parts[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", uid)
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

// ListParts возвращает список деталей, соответствующих заданным фильтрам.
//
// Если фильтр не указан, возвращаются все детали. Поддерживается фильтрация по UUID, имени, категории,
// стране производителя и тегам.
func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()
	if filter == nil {
		parts := make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	filteredParts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, part := range s.parts {
		filteredParts = append(filteredParts, part)
	}

	if len(filter.Uuids) > 0 {
		uuidSet := make(map[string]bool)
		for _, uuid := range filter.Uuids {
			uuidSet[uuid] = true
		}

		tempParts := make([]*inventoryV1.Part, 0)
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

		tempParts := make([]*inventoryV1.Part, 0)
		for _, part := range filteredParts {
			if nameSet[part.Name] {
				tempParts = append(tempParts, part)
			}
		}
		filteredParts = tempParts
	}

	if len(filter.Categories) > 0 {
		categorySet := make(map[inventoryV1.Category]bool)
		for _, category := range filter.Categories {
			categorySet[category] = true
		}

		tempParts := make([]*inventoryV1.Part, 0)
		for _, part := range filteredParts {
			if categorySet[part.Category] {
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

		tempParts := make([]*inventoryV1.Part, 0)
		for _, part := range filteredParts {
			if part.Manufacturer != nil && countrySet[part.Manufacturer.Country] {
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

		tempParts := make([]*inventoryV1.Part, 0)
		for _, part := range filteredParts {
			for _, partTag := range part.Tags {
				if tagSet[partTag] {
					tempParts = append(tempParts, part)
					break
				}
			}
		}
		filteredParts = tempParts
	}

	return &inventoryV1.ListPartsResponse{Parts: filteredParts}, nil
}
