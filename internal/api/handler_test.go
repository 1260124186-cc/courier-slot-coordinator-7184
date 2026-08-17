package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/repository"
	"github.com/1260124186-cc/courier-slot-coordinator/internal/service"
)

func TestCreateShipmentEndpoint(t *testing.T) {
	handler := NewHandler(service.NewDispatchService(repository.NewMemoryStore()))
	request := httptest.NewRequest(http.MethodPost, "/shipments", bytes.NewBufferString(
		`{"recipient":"Ava","zone":"north","window":"09:00-11:00","packages":[{"sku":"book","units":2}]}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", response.Code, http.StatusCreated, response.Body.String())
	}
}

func TestUnknownZoneSummaryEndpointIsRejected(t *testing.T) {
	handler := NewHandler(service.NewDispatchService(repository.NewMemoryStore()))
	request := httptest.NewRequest(http.MethodGet, "/zones/central/summary", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", response.Code, http.StatusNotFound, response.Body.String())
	}
}
