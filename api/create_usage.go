package api

import (
	"encoding/json"
	"finalcourseproject/model"
	"finalcourseproject/repository"
	"net/http"
)

type UsageRequest struct {
	DeviceName string  `json:"device_name"`
	Kwh        float64 `json:"kwh"`
	UsageTime  float64 `json:"usage_time"`
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
		Name:      req.DeviceName,
		Kwh:       req.Kwh,
		UsageTime: req.UsageTime,
	}

	err := h.usageRepository.Create(usage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usage)
}
