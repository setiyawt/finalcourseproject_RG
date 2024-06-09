package api

import (
	"encoding/json"
	"finalcourseproject/model"
	"finalcourseproject/repository"
	"net/http"
)

type UsageRequest struct {
	DeviceName string  `json:"device_name"`
	Cost       float64 `json:"cost"`
	Usage_Kwh  float64 `json:"usage_kwh"`
}

type UsageHandler struct {
	usageRepository repository.ElectricityUsagesRepository
}

func NewUsageHandler(usageRepository repository.ElectricityUsagesRepository) *UsageHandler {
	return &UsageHandler{usageRepository: usageRepository}
}

func (h *UsageHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req UsageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usage := model.ElectricityUsages{
		Name: req.DeviceName,
		//Cost:      req.Cost,
		Usage_Kwh: req.Usage_Kwh,
	}

	if err := h.usageRepository.Create(usage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"id": usage.ID})
}
