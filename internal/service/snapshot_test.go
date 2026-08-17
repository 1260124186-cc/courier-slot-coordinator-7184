package service

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
)

func TestReadingShipmentCannotChangeZoneCapacity(t *testing.T) {
	dispatch := NewDispatchService(repository.NewMemoryStore())
	first, err := dispatch.CreateShipment(context.Background(), CreateShipmentInput{
		Recipient: "Ava",
		Zone:      "west",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "parcel", Units: 10}},
	})
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}

	snapshot, err := dispatch.GetShipment(context.Background(), first.ID)
	if err != nil {
		t.Fatalf("GetShipment() error = %v", err)
	}
	snapshot.Packages[0].Units = 1

	_, err = dispatch.CreateShipment(context.Background(), CreateShipmentInput{
		Recipient: "Noah",
		Zone:      "west",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "parcel", Units: 10}},
	})
	if !errors.Is(err, domain.ErrZoneCapacity) {
		t.Fatalf("CreateShipment() error = %v, want ErrZoneCapacity after a caller mutates a read result", err)
	}
}
