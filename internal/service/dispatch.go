package service

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
)

type CreateShipmentInput struct {
	Recipient string           `json:"recipient"`
	Zone      string           `json:"zone"`
	Window    string           `json:"window"`
	Packages  []domain.Package `json:"packages"`
}

type DispatchService struct {
	store            repository.Store
	clock            func() time.Time
	sequence         atomic.Uint64
	zonePackageLimit map[string]int
}

func NewDispatchService(store repository.Store) *DispatchService {
	return &DispatchService{
		store: store,
		clock: time.Now,
		zonePackageLimit: map[string]int{
			"north": 24,
			"south": 18,
			"east":  20,
			"west":  16,
		},
	}
}

func (s *DispatchService) CreateShipment(ctx context.Context, input CreateShipmentInput) (domain.Shipment, error) {
	now := s.clock().UTC()
	shipment := domain.Shipment{
		ID:        fmt.Sprintf("shipment-%06d", s.sequence.Add(1)),
		Recipient: strings.TrimSpace(input.Recipient),
		Zone:      strings.ToLower(strings.TrimSpace(input.Zone)),
		Window:    strings.TrimSpace(input.Window),
		Packages:  domain.ClonePackages(input.Packages),
		Status:    domain.StatusReady,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := shipment.Validate(); err != nil {
		return domain.Shipment{}, err
	}
	if err := s.ensureZoneCapacity(ctx, shipment.Zone, shipment.PackageUnits()); err != nil {
		return domain.Shipment{}, err
	}
	if err := s.store.Create(ctx, shipment); err != nil {
		return domain.Shipment{}, err
	}
	return domain.CloneShipment(shipment), nil
}

func (s *DispatchService) AssignCourier(ctx context.Context, shipmentID, courierID string) (domain.Shipment, error) {
	courierID = strings.TrimSpace(courierID)
	if courierID == "" {
		return domain.Shipment{}, domain.ErrCourierRequired
	}

	shipment, err := s.store.Get(ctx, shipmentID)
	if err != nil {
		return domain.Shipment{}, err
	}
	if shipment.Status == domain.StatusDelivered || shipment.Status == domain.StatusCancelled {
		return domain.Shipment{}, domain.ErrShipmentFinalized
	}
	shipment.CourierID = courierID
	shipment.UpdatedAt = s.clock().UTC()
	return s.store.Update(ctx, shipment, shipment.Version)
}

func (s *DispatchService) Collect(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	return s.transition(ctx, shipmentID, domain.StatusCollected)
}

func (s *DispatchService) Deliver(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	return s.transition(ctx, shipmentID, domain.StatusDelivered)
}

func (s *DispatchService) Cancel(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	return s.transition(ctx, shipmentID, domain.StatusCancelled)
}

func (s *DispatchService) transition(ctx context.Context, shipmentID string, next domain.ShipmentStatus) (domain.Shipment, error) {
	shipment, err := s.store.Get(ctx, shipmentID)
	if err != nil {
		return domain.Shipment{}, err
	}
	if next == domain.StatusCollected && shipment.CourierID == "" {
		return domain.Shipment{}, domain.ErrCourierRequired
	}
	if !shipment.CanTransition(next) {
		return domain.Shipment{}, domain.ErrInvalidTransition
	}
	shipment.Status = next
	shipment.UpdatedAt = s.clock().UTC()
	return s.store.Update(ctx, shipment, shipment.Version)
}

func (s *DispatchService) GetShipment(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	shipment, err := s.store.Get(ctx, shipmentID)
	if err != nil {
		return domain.Shipment{}, err
	}
	return domain.CloneShipment(shipment), nil
}

func (s *DispatchService) ensureZoneCapacity(ctx context.Context, zone string, requestedUnits int) error {
	limit, knownZone := s.zonePackageLimit[zone]
	if !knownZone {
		return fmt.Errorf("%w: %s", domain.ErrInvalidShipment, zone)
	}
	shipments, err := s.store.ListByZone(ctx, zone)
	if err != nil {
		return err
	}

	usedUnits := 0
	for _, shipment := range shipments {
		if shipment.Status != domain.StatusCancelled && shipment.Status != domain.StatusDelivered {
			usedUnits += shipment.PackageUnits()
		}
	}
	if usedUnits+requestedUnits > limit {
		return domain.ErrZoneCapacity
	}
	return nil
}
