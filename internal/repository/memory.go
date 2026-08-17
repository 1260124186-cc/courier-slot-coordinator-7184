package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
)

type Store interface {
	Create(context.Context, domain.Shipment) error
	Get(context.Context, string) (domain.Shipment, error)
	Update(context.Context, domain.Shipment, int64) (domain.Shipment, error)
	ListByZone(context.Context, string) ([]domain.Shipment, error)
}

type MemoryStore struct {
	mu        sync.RWMutex
	shipments map[string]domain.Shipment
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{shipments: make(map[string]domain.Shipment)}
}

func (s *MemoryStore) Create(ctx context.Context, shipment domain.Shipment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := shipment.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.shipments[shipment.ID]; exists {
		return domain.ErrAlreadyExists
	}
	s.shipments[shipment.ID] = domain.CloneShipment(shipment)
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, id string) (domain.Shipment, error) {
	if err := ctx.Err(); err != nil {
		return domain.Shipment{}, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	shipment, exists := s.shipments[id]
	if !exists {
		return domain.Shipment{}, domain.ErrNotFound
	}
	return domain.CloneShipment(shipment), nil
}

func (s *MemoryStore) Update(ctx context.Context, shipment domain.Shipment, expectedVersion int64) (domain.Shipment, error) {
	if err := ctx.Err(); err != nil {
		return domain.Shipment{}, err
	}
	if err := shipment.Validate(); err != nil {
		return domain.Shipment{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.shipments[shipment.ID]
	if !exists {
		return domain.Shipment{}, domain.ErrNotFound
	}
	if current.Version != expectedVersion {
		return domain.Shipment{}, fmt.Errorf("update shipment conflict: %v", domain.ErrConflict)
	}
	shipment.Version = current.Version + 1
	s.shipments[shipment.ID] = domain.CloneShipment(shipment)
	return domain.CloneShipment(shipment), nil
}

func (s *MemoryStore) ListByZone(ctx context.Context, zone string) ([]domain.Shipment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	shipments := make([]domain.Shipment, 0)
	for _, shipment := range s.shipments {
		if shipment.Zone == zone {
			shipments = append(shipments, domain.CloneShipment(shipment))
		}
	}
	return shipments, nil
}
