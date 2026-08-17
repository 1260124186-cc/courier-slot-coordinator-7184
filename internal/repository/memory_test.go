package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
)

func TestMemoryStoreUsesOptimisticVersioning(t *testing.T) {
	store := NewMemoryStore()
	shipment := domain.Shipment{
		ID:        "shipment-1",
		Recipient: "Ava",
		Zone:      "north",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "book", Units: 2}},
		Status:    domain.StatusReady,
		Version:   1,
	}
	if err := store.Create(context.Background(), shipment); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	shipment.CourierID = "courier-7"
	updated, err := store.Update(context.Background(), shipment, 1)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Version != 2 {
		t.Fatalf("Update() version = %d, want 2", updated.Version)
	}
	if _, err := store.Update(context.Background(), shipment, 1); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("stale Update() error = %v, want ErrConflict", err)
	}
}
