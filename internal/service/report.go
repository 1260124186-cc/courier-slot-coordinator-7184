package service

import (
	"context"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
)

type ZoneSummary struct {
	Zone       string `json:"zone"`
	Ready      int    `json:"ready"`
	Collected  int    `json:"collected"`
	Delivered  int    `json:"delivered"`
	Cancelled  int    `json:"cancelled"`
	PackageQty int    `json:"package_qty"`
}

func (s *DispatchService) ZoneSummary(ctx context.Context, zone string) (ZoneSummary, error) {
	shipments, err := s.store.ListByZone(ctx, zone)
	if err != nil {
		return ZoneSummary{}, err
	}

	summary := ZoneSummary{Zone: zone}
	for _, shipment := range shipments {
		summary.PackageQty += shipment.PackageUnits()
		switch shipment.Status {
		case domain.StatusReady:
			summary.Ready++
		case domain.StatusCollected:
			summary.Collected++
		case domain.StatusDelivered:
			summary.Delivered++
		case domain.StatusCancelled:
			summary.Cancelled++
		}
	}
	return summary, nil
}
