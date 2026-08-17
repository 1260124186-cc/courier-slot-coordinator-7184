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
)

type OperationError struct {
	Operation string
	Cause     error
}

func (e *OperationError) Error() string {
	return e.Operation + ": " + e.Cause.Error()
}

func (e *OperationError) Unwrap() error {
	return e.Cause
}

func WrapOperation(operation string, cause error) error {
	if cause == nil {
		return nil
	}
	return &OperationError{Operation: operation, Cause: cause}
}
