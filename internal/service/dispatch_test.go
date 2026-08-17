package service

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
)

func TestDispatchLifecycleAndSummary(t *testing.T) {
	dispatch := NewDispatchService(repository.NewMemoryStore())
	shipment, err := dispatch.CreateShipment(context.Background(), CreateShipmentInput{
		Recipient: "Ava",
		Zone:      "north",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "book", Units: 2}},
	})
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}

	if _, err := dispatch.Collect(context.Background(), shipment.ID); !errors.Is(err, domain.ErrCourierRequired) {
		t.Fatalf("Collect() error = %v, want ErrCourierRequired", err)
	}
	if _, err := dispatch.AssignCourier(context.Background(), shipment.ID, "courier-7"); err != nil {
		t.Fatalf("AssignCourier() error = %v", err)
	}
	if _, err := dispatch.Collect(context.Background(), shipment.ID); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if _, err := dispatch.Deliver(context.Background(), shipment.ID); err != nil {
		t.Fatalf("Deliver() error = %v", err)
	}

	summary, err := dispatch.ZoneSummary(context.Background(), "north")
	if err != nil {
		t.Fatalf("ZoneSummary() error = %v", err)
	}
	if summary.Delivered != 1 || summary.PackageQty != 2 {
		t.Fatalf("ZoneSummary() = %+v, want one delivered shipment with two packages", summary)
	}
}

func TestCreateShipmentRejectsCapacityOverflow(t *testing.T) {
	dispatch := NewDispatchService(repository.NewMemoryStore())
	input := CreateShipmentInput{
		Recipient: "Ava",
		Zone:      "west",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "book", Units: 10}},
	}
	if _, err := dispatch.CreateShipment(context.Background(), input); err != nil {
		t.Fatalf("first CreateShipment() error = %v", err)
	}
	input.Recipient = "Noah"
	if _, err := dispatch.CreateShipment(context.Background(), input); !errors.Is(err, domain.ErrZoneCapacity) {
		t.Fatalf("second CreateShipment() error = %v, want ErrZoneCapacity", err)
	}
}

func TestZoneSummaryRejectsUnconfiguredZone(t *testing.T) {
	dispatch := NewDispatchService(repository.NewMemoryStore())
	if _, err := dispatch.ZoneSummary(context.Background(), "central"); !errors.Is(err, domain.ErrUnknownZone) {
		t.Fatalf("ZoneSummary() error = %v, want ErrUnknownZone", err)
	}
}
