package domain

import "testing"

func TestShipmentValidationAndClone(t *testing.T) {
	shipment := Shipment{
		ID:        "shipment-1",
		Recipient: "Ava",
		Zone:      "north",
		Window:    "09:00-11:00",
		Packages:  []Package{{SKU: "book", Units: 2}},
		Status:    StatusReady,
	}
	if err := shipment.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	clone := CloneShipment(shipment)
	clone.Packages[0].Units = 9
	if shipment.Packages[0].Units != 2 {
		t.Fatalf("CloneShipment() did not isolate package slice")
	}
}
