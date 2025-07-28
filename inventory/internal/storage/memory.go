package storage

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/nexarise/microservices-factory/shared/pkg/proto/inventory/v1"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	parts map[string]*inventoryv1.Part
}

func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		parts: make(map[string]*inventoryv1.Part),
	}

	storage.initSampleData()

	return storage
}

func (s *MemoryStorage) GetPart(uuid string) (*inventoryv1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, exists := s.parts[uuid]
	if !exists {
		return nil, fmt.Errorf("part with UUID %s not found", uuid)
	}

	return part, nil
}

func (s *MemoryStorage) ListParts(filter *inventoryv1.PartsFilter) ([]*inventoryv1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*inventoryv1.Part

	if isEmptyFilter(filter) {
		for _, part := range s.parts {
			result = append(result, part)
		}
		return result, nil
	}

	for _, part := range s.parts {
		if matchesPart(part, filter) {
			result = append(result, part)
		}
	}

	return result, nil
}

func isEmptyFilter(filter *inventoryv1.PartsFilter) bool {
	if filter == nil {
		return true
	}

	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func matchesUuids(part *inventoryv1.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}

	for _, uuid := range uuids {
		if part.Uuid == uuid {
			return true
		}
	}
	return false
}

func matchesNames(part *inventoryv1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}

	for _, name := range names {
		if strings.Contains(strings.ToLower(part.Name), strings.ToLower(name)) {
			return true
		}
	}
	return false
}

func matchesCategories(part *inventoryv1.Part, categories []inventoryv1.Category) bool {
	if len(categories) == 0 {
		return true
	}

	for _, category := range categories {
		if part.Category == category {
			return true
		}
	}
	return false
}

func matchesManufacturerCountries(part *inventoryv1.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}

	for _, country := range countries {
		if strings.EqualFold(part.Manufacturer.Country, country) {
			return true
		}
	}
	return false
}

func matchesTags(part *inventoryv1.Part, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	for _, filterTag := range filterTags {
		for _, partTag := range part.Tags {
			if strings.EqualFold(partTag, filterTag) {
				return true
			}
		}
	}
	return false
}

func matchesPart(part *inventoryv1.Part, filter *inventoryv1.PartsFilter) bool {
	return matchesUuids(part, filter.Uuids) &&
		matchesNames(part, filter.Names) &&
		matchesCategories(part, filter.Categories) &&
		matchesManufacturerCountries(part, filter.ManufacturerCountries) &&
		matchesTags(part, filter.Tags)
}

func (s *MemoryStorage) initSampleData() {
	now := timestamppb.New(time.Now())

	parts := []*inventoryv1.Part{
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion system for long-distance space travel",
			Price:         1500000.50,
			StockQuantity: 5,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 250.0,
				Width:  120.0,
				Height: 180.0,
				Weight: 5000.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "Germany",
				Website: "https://spacetech-industries.com",
			},
			Tags: []string{"quantum", "engine", "premium", "long-range"},
			Metadata: map[string]*structpb.Value{
				"power_output":  structpb.NewStringValue("15.2 TW"),
				"efficiency":    structpb.NewNumberValue(98.5),
				"certification": structpb.NewStringValue("ISO-SPACE-9001"),
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Liquid Hydrogen Fuel Cell",
			Description:   "Clean-burning hydrogen fuel for environmental sustainability",
			Price:         25000.75,
			StockQuantity: 50,
			Category:      inventoryv1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryv1.Dimensions{
				Length: 80.0,
				Width:  80.0,
				Height: 120.0,
				Weight: 150.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "EcoFuel Corp",
				Country: "Japan",
				Website: "https://ecofuel.jp",
			},
			Tags: []string{"hydrogen", "fuel", "eco-friendly", "clean"},
			Metadata: map[string]*structpb.Value{
				"energy_density": structpb.NewStringValue("142 MJ/kg"),
				"purity":         structpb.NewNumberValue(99.99),
				"storage_temp":   structpb.NewStringValue("-253°C"),
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Reinforced Observation Porthole",
			Description:   "Ultra-strong transparent aluminum porthole for safe space observation",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 60.0,
				Width:  60.0,
				Height: 15.0,
				Weight: 45.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "ClearSpace Optics",
				Country: "USA",
				Website: "https://clearspace-optics.com",
			},
			Tags: []string{"porthole", "observation", "reinforced", "transparent"},
			Metadata: map[string]*structpb.Value{
				"material":            structpb.NewStringValue("Transparent Aluminum"),
				"pressure_resistance": structpb.NewStringValue("15 ATM"),
				"transparency":        structpb.NewNumberValue(99.8),
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Adaptive Solar Wing",
			Description:   "Self-adjusting solar panel wing for maximum energy efficiency",
			Price:         850000.25,
			StockQuantity: 8,
			Category:      inventoryv1.Category_CATEGORY_WING,
			Dimensions: &inventoryv1.Dimensions{
				Length: 1200.0,
				Width:  300.0,
				Height: 25.0,
				Weight: 2500.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SolarWings Ltd",
				Country: "Germany",
				Website: "https://solarwings.de",
			},
			Tags: []string{"solar", "wing", "adaptive", "energy"},
			Metadata: map[string]*structpb.Value{
				"power_generation": structpb.NewStringValue("500 kW"),
				"efficiency":       structpb.NewNumberValue(45.2),
				"auto_tracking":    structpb.NewBoolValue(true),
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Unknown Component XJ-2024",
			Description:   "Mysterious component found in deep space wreckage",
			Price:         999999.99,
			StockQuantity: 1,
			Category:      inventoryv1.Category_CATEGORY_UNKNOWN,
			Dimensions: &inventoryv1.Dimensions{
				Length: 42.0,
				Width:  42.0,
				Height: 42.0,
				Weight: 424.2,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Unknown",
				Country: "Unknown",
				Website: "",
			},
			Tags: []string{"unknown", "mysterious", "alien", "rare"},
			Metadata: map[string]*structpb.Value{
				"energy_signature": structpb.NewStringValue("Unidentified"),
				"material":         structpb.NewStringValue("Unknown alloy"),
				"age":              structpb.NewStringValue("> 1000 years"),
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, part := range parts {
		s.parts[part.Uuid] = part
	}
}
