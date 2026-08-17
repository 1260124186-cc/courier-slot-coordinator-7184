package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/domain"
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

type conflictStore struct{}

func (conflictStore) Create(context.Context, domain.Shipment) error {
	return nil
}

func (conflictStore) Get(context.Context, string) (domain.Shipment, error) {
	return domain.Shipment{
		ID:        "shipment-1",
		Recipient: "Ava",
		Zone:      "west",
		Window:    "09:00-11:00",
		Packages:  []domain.Package{{SKU: "desk", Units: 10}},
		Status:    domain.StatusReady,
		Version:   1,
	}, nil
}

func (conflictStore) Update(context.Context, domain.Shipment, int64) (domain.Shipment, error) {
	return domain.Shipment{}, domain.ErrConflict
}

func (conflictStore) ListByZone(context.Context, string) ([]domain.Shipment, error) {
	return nil, nil
}

func TestCourierConflictKeepsHTTPConflictStatus(t *testing.T) {
	handler := NewHandler(service.NewDispatchService(conflictStore{}))
	request := httptest.NewRequest(http.MethodPost, "/shipments/shipment-1/courier", bytes.NewBufferString(
		`{"courier_id":"courier-7"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body=%s", response.Code, http.StatusConflict, response.Body.String())
	}
}
