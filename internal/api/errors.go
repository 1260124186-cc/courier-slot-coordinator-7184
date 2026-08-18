package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
)

func serviceErrorStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrUnknownZone):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidShipment),
		errors.Is(err, domain.ErrInvalidTransition),
		errors.Is(err, domain.ErrCourierRequired),
		errors.Is(err, domain.ErrShipmentFinalized),
		errors.Is(err, domain.ErrZoneCapacity):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeServiceError(writer http.ResponseWriter, err error) {
	writeError(writer, serviceErrorStatus(err), err)
}

func writeError(writer http.ResponseWriter, status int, err error) {
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
