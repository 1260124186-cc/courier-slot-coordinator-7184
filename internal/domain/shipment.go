package domain

import (
	"fmt"
	"strings"
	"time"
)

type ShipmentStatus string

const (
	StatusReady     ShipmentStatus = "ready"
	StatusCollected ShipmentStatus = "collected"
	StatusDelivered ShipmentStatus = "delivered"
	StatusCancelled ShipmentStatus = "cancelled"
)

type Package struct {
	SKU   string `json:"sku"`
	Units int    `json:"units"`
}

type Shipment struct {
	ID        string         `json:"id"`
	Recipient string         `json:"recipient"`
	Zone      string         `json:"zone"`
	Window    string         `json:"window"`
	CourierID string         `json:"courier_id,omitempty"`
	Status    ShipmentStatus `json:"status"`
	Packages  []Package      `json:"packages"`
	Version   int64          `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (s Shipment) PackageUnits() int {
	total := 0
	for _, item := range s.Packages {
		total += item.Units
	}
	return total
}

func (s Shipment) Validate() error {
	if strings.TrimSpace(s.ID) == "" ||
		strings.TrimSpace(s.Recipient) == "" ||
		strings.TrimSpace(s.Zone) == "" ||
		strings.TrimSpace(s.Window) == "" {
		return ErrInvalidShipment
	}
	if len(s.Packages) == 0 {
		return ErrInvalidShipment
	}
	for _, item := range s.Packages {
		if strings.TrimSpace(item.SKU) == "" || item.Units <= 0 {
			return fmt.Errorf("%w: package sku and units are required", ErrInvalidShipment)
		}
	}
	return nil
}

func (s Shipment) CanTransition(next ShipmentStatus) bool {
	switch s.Status {
	case StatusReady:
		return next == StatusCollected || next == StatusCancelled
	case StatusCollected:
		return next == StatusDelivered || next == StatusCancelled
	default:
		return false
	}
}

func CloneShipment(source Shipment) Shipment {
	clone := source
	clone.Packages = append([]Package(nil), source.Packages...)
	return clone
}
