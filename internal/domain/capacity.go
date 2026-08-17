package domain

type ZoneCapacity struct {
	Limit int
	Used  int
}

func NewZoneCapacity(limit int) ZoneCapacity {
	return ZoneCapacity{Limit: limit}
}

func (c *ZoneCapacity) Include(shipment Shipment) {
	if shipment.CountsAgainstCapacity() {
		c.Used += shipment.PackageUnits()
	}
}

func (c ZoneCapacity) CanAccept(shipment Shipment) bool {
	return shipment.CountsAgainstCapacity() && c.Used+shipment.PackageUnits() <= c.Limit
}
