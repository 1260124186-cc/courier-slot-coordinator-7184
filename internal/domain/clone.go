package domain

// ClonePackages returns a deep copy of the package slice so callers can never
// mutate the stored shipment through a returned read result.
func ClonePackages(packages []Package) []Package {
	if packages == nil {
		return nil
	}
	clone := make([]Package, len(packages))
	copy(clone, packages)
	return clone
}

// CloneShipment returns a copy of the shipment with a fully isolated package slice.
func CloneShipment(source Shipment) Shipment {
	clone := source
	clone.Packages = ClonePackages(source.Packages)
	return clone
}
