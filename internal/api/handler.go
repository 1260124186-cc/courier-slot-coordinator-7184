package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/1260124186-cc/courier-slot-coordinator/internal/service"
)

type Handler struct {
	dispatch *service.DispatchService
	mux      *http.ServeMux
}

func NewHandler(dispatch *service.DispatchService) http.Handler {
	handler := &Handler{dispatch: dispatch, mux: http.NewServeMux()}
	handler.mux.HandleFunc("POST /shipments", handler.createShipment)
	handler.mux.HandleFunc("GET /shipments/{shipmentID}", handler.getShipment)
	handler.mux.HandleFunc("POST /shipments/{shipmentID}/courier", handler.assignCourier)
	handler.mux.HandleFunc("POST /shipments/{shipmentID}/collect", handler.collectShipment)
	handler.mux.HandleFunc("POST /shipments/{shipmentID}/deliver", handler.deliverShipment)
	handler.mux.HandleFunc("GET /zones/{zone}/summary", handler.zoneSummary)
	return handler
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.mux.ServeHTTP(writer, request)
}

func (h *Handler) createShipment(writer http.ResponseWriter, request *http.Request) {
	var input service.CreateShipmentInput
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	shipment, err := h.dispatch.CreateShipment(request.Context(), input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, shipment)
}

func (h *Handler) getShipment(writer http.ResponseWriter, request *http.Request) {
	shipment, err := h.dispatch.GetShipment(request.Context(), request.PathValue("shipmentID"))
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, shipment)
}

func (h *Handler) assignCourier(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		CourierID string `json:"courier_id"`
	}
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, err)
		return
	}
	shipment, err := h.dispatch.AssignCourier(request.Context(), request.PathValue("shipmentID"), input.CourierID)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, shipment)
}

func (h *Handler) collectShipment(writer http.ResponseWriter, request *http.Request) {
	shipment, err := h.dispatch.Collect(request.Context(), request.PathValue("shipmentID"))
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, shipment)
}

func (h *Handler) deliverShipment(writer http.ResponseWriter, request *http.Request) {
	shipment, err := h.dispatch.Deliver(request.Context(), request.PathValue("shipmentID"))
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, shipment)
}

func (h *Handler) zoneSummary(writer http.ResponseWriter, request *http.Request) {
	zone := strings.ToLower(strings.TrimSpace(request.PathValue("zone")))
	summary, err := h.dispatch.ZoneSummary(request.Context(), zone)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, summary)
}

func decodeJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
