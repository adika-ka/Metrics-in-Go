package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"task4.2.3/internal/controller"
	"task4.2.3/internal/models"
	"task4.2.3/internal/monitoring"
)

type AddressHandler struct {
	controller *controller.AddressController
}

func NewAddressHandler(controller *controller.AddressController) *AddressHandler {
	return &AddressHandler{controller: controller}
}

// @Summary Поиск адреса
// @Description Принимает текстовый запрос и возвращает список адресов
// @Tags address
// @Accept json
// @Produce json
// @Param request body models.SearchRequest true "Запрос"
// @Success 200 {object} models.SearchResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/search [post]
// @Security BearerAuth
func (h *AddressHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	monitoring.HTTPRequestTotal.WithLabelValues("/api/address/search").Inc()
	startTime := time.Now()

	var req models.SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	addresses, err := h.controller.Search(req.Query)
	if err != nil {
		handleAPIError(w, err)
		return
	}

	resp := models.SearchResponse{Addresses: convertToPointerSlice(addresses)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	duration := time.Since(startTime).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("/api/address/search").Observe(duration)
}

// @Summary Поиск адреса
// @Description Принимает текстовый запрос с координатами ("lat" - широта, "lng" - долгота) и возвращает список адресов
// @Tags address
// @Accept json
// @Produce json
// @Param request body models.GeocodeRequest true "Запрос"
// @Success 200 {object} models.GeocodeResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/geocode [post]
// @Security BearerAuth
func (h *AddressHandler) GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	monitoring.HTTPRequestTotal.WithLabelValues("/api/address/geocode").Inc()
	startTime := time.Now()

	var req models.GeocodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	addresses, err := h.controller.Geocode(req.Lat, req.Lng)
	if err != nil {
		handleAPIError(w, err)
		return
	}

	resp := models.GeocodeResponse{Addresses: convertToPointerSlice(addresses)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	duration := time.Since(startTime).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("/api/address/geocode").Observe(duration)
}

func handleAPIError(w http.ResponseWriter, err error) {
	if err.Error() == "no addresses found" {
		http.Error(w, "No addresses found", http.StatusNotFound)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func convertToPointerSlice(addresses []models.Address) []*models.Address {
	addressPointers := make([]*models.Address, len(addresses))
	for i := range addresses {
		addressPointers[i] = &addresses[i]
	}
	return addressPointers
}
