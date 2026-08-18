package service

import (
	"fmt"
	"strings"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
)

func normalizeZone(zone string) string {
	return strings.ToLower(strings.TrimSpace(zone))
}

func (s *DispatchService) validateZone(zone string) error {
	if _, err := s.zoneLimit(zone); err != nil {
		return err
	}
	return nil
}

func (s *DispatchService) zoneLimit(zone string) (int, error) {
	limit, known := s.zonePackageLimit[zone]
	if !known {
		return 0, fmt.Errorf("%w: %s", domain.ErrUnknownZone, zone)
	}
	return limit, nil
}
