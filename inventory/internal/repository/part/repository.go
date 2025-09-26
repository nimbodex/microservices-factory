package part

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// MemoryPartRepository implements PartRepository using in-memory storage
type MemoryPartRepository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

// NewMemoryPartRepository creates a new in-memory part repository
func NewMemoryPartRepository() *MemoryPartRepository {
	repo := &MemoryPartRepository{
		parts: make(map[string]*model.Part),
	}

	repo.initSampleData()
	return repo
}

// GetByUUID retrieves a part by its UUID
func (r *MemoryPartRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid.String()]
	if !exists {
		return nil, fmt.Errorf("part with UUID %s not found", uuid)
	}

	// Return a copy to avoid external modifications
	partCopy := *part
	return &partCopy, nil
}

// List retrieves parts matching the filter criteria
func (r *MemoryPartRepository) List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Part

	if isEmptyFilter(filter) {
		for _, part := range r.parts {
			partCopy := *part
			result = append(result, &partCopy)
		}
		return result, nil
	}

	for _, part := range r.parts {
		if matchesPart(part, filter) {
			partCopy := *part
			result = append(result, &partCopy)
		}
	}

	return result, nil
}

// Create creates a new part
func (r *MemoryPartRepository) Create(ctx context.Context, part *model.Part) error {
	if part == nil {
		return fmt.Errorf("part cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	partKey := part.UUID.String()
	if _, exists := r.parts[partKey]; exists {
		return fmt.Errorf("part with UUID %s already exists", part.UUID)
	}

	// Create a copy to avoid external modifications
	partCopy := *part
	r.parts[partKey] = &partCopy

	return nil
}

// Update updates an existing part
func (r *MemoryPartRepository) Update(ctx context.Context, part *model.Part) error {
	if part == nil {
		return fmt.Errorf("part cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	partKey := part.UUID.String()
	if _, exists := r.parts[partKey]; !exists {
		return fmt.Errorf("part with UUID %s not found", part.UUID)
	}

	// Create a copy to avoid external modifications
	partCopy := *part
	r.parts[partKey] = &partCopy

	return nil
}

// Delete removes a part by its UUID
func (r *MemoryPartRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	partKey := uuid.String()
	if _, exists := r.parts[partKey]; !exists {
		return fmt.Errorf("part with UUID %s not found", uuid)
	}

	delete(r.parts, partKey)
	return nil
}

func isEmptyFilter(filter *model.PartsFilter) bool {
	if filter == nil {
		return true
	}

	return len(filter.UUIDs) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func matchesPart(part *model.Part, filter *model.PartsFilter) bool {
	return matchesUUIDs(part, filter.UUIDs) &&
		matchesNames(part, filter.Names) &&
		matchesCategories(part, filter.Categories) &&
		matchesManufacturerCountries(part, filter.ManufacturerCountries) &&
		matchesTags(part, filter.Tags)
}

func matchesUUIDs(part *model.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}

	for _, uuid := range uuids {
		if part.UUID.String() == uuid {
			return true
		}
	}
	return false
}

func matchesNames(part *model.Part, names []string) bool {
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

func matchesCategories(part *model.Part, categories []inventoryv1.Category) bool {
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

func matchesManufacturerCountries(part *model.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}

	if part.Manufacturer == nil {
		return false
	}

	for _, country := range countries {
		if strings.EqualFold(part.Manufacturer.Country, country) {
			return true
		}
	}
	return false
}

func matchesTags(part *model.Part, filterTags []string) bool {
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

func (r *MemoryPartRepository) initSampleData() {
	now := time.Now()

	parts := []*model.Part{
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion system for long-distance space travel",
			Price:         1500000.50,
			StockQuantity: 5,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			Dimensions: &model.Dimensions{
				Length: 250.0,
				Width:  120.0,
				Height: 180.0,
				Weight: 5000.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "Germany",
				Website: "https://spacetech-industries.com",
			},
			Tags: []string{"quantum", "engine", "premium", "long-range"},
			Metadata: map[string]interface{}{
				"power_output":  "15.2 TW",
				"efficiency":    98.5,
				"certification": "ISO-SPACE-9001",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "Liquid Hydrogen Fuel Cell",
			Description:   "Clean-burning hydrogen fuel for environmental sustainability",
			Price:         25000.75,
			StockQuantity: 50,
			Category:      inventoryv1.Category_CATEGORY_FUEL,
			Dimensions: &model.Dimensions{
				Length: 80.0,
				Width:  80.0,
				Height: 120.0,
				Weight: 150.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "EcoFuel Corp",
				Country: "Japan",
				Website: "https://ecofuel.jp",
			},
			Tags: []string{"hydrogen", "fuel", "eco-friendly", "clean"},
			Metadata: map[string]interface{}{
				"energy_density": "142 MJ/kg",
				"purity":         99.99,
				"storage_temp":   "-253°C",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:          "Reinforced Observation Porthole",
			Description:   "Ultra-strong transparent aluminum porthole for safe space observation",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
			Dimensions: &model.Dimensions{
				Length: 60.0,
				Width:  60.0,
				Height: 15.0,
				Weight: 45.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "ClearSpace Optics",
				Country: "USA",
				Website: "https://clearspace-optics.com",
			},
			Tags: []string{"porthole", "observation", "reinforced", "transparent"},
			Metadata: map[string]interface{}{
				"material":            "Transparent Aluminum",
				"pressure_resistance": "15 ATM",
				"transparency":        99.8,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:          "Adaptive Solar Wing",
			Description:   "Self-adjusting solar panel wing for maximum energy efficiency",
			Price:         850000.25,
			StockQuantity: 8,
			Category:      inventoryv1.Category_CATEGORY_WING,
			Dimensions: &model.Dimensions{
				Length: 1200.0,
				Width:  300.0,
				Height: 25.0,
				Weight: 2500.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "SolarWings Ltd",
				Country: "Germany",
				Website: "https://solarwings.de",
			},
			Tags: []string{"solar", "wing", "adaptive", "energy"},
			Metadata: map[string]interface{}{
				"power_generation": "500 kW",
				"efficiency":       45.2,
				"auto_tracking":    true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:          "Unknown Component XJ-2024",
			Description:   "Mysterious component found in deep space wreckage",
			Price:         999999.99,
			StockQuantity: 1,
			Category:      inventoryv1.Category_CATEGORY_UNKNOWN,
			Dimensions: &model.Dimensions{
				Length: 42.0,
				Width:  42.0,
				Height: 42.0,
				Weight: 424.2,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "Unknown",
				Country: "Unknown",
				Website: "",
			},
			Tags: []string{"unknown", "mysterious", "alien", "rare"},
			Metadata: map[string]interface{}{
				"energy_signature": "Unidentified",
				"material":         "Unknown alloy",
				"age":              "> 1000 years",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, part := range parts {
		r.parts[part.UUID.String()] = part
	}
}
