package domain

// CountsTowardZoneCapacity reports whether a shipment still consumes a zone's
// active delivery capacity.
func (s Shipment) CountsTowardZoneCapacity() bool {
	return !s.IsFinalized()
}
