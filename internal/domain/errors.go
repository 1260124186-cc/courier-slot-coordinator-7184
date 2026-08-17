package domain

import "errors"

var (
	ErrNotFound          = errors.New("shipment not found")
	ErrAlreadyExists     = errors.New("shipment already exists")
	ErrConflict          = errors.New("shipment changed while being updated")
	ErrInvalidShipment   = errors.New("invalid shipment")
	ErrInvalidTransition = errors.New("invalid shipment status transition")
	ErrCourierRequired   = errors.New("courier assignment is required before collection")
	ErrZoneCapacity      = errors.New("zone package capacity exceeded")
	ErrShipmentFinalized = errors.New("shipment is already finalized")
	ErrUnknownZone       = errors.New("delivery zone is not configured")
)
