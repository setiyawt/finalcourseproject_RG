package api

import (
	"encoding/json"
	"finalcourseproject/service"
	"net/http"
)

type PredictHandler struct {
	predictionService service.PredictionService
}

func NewPredictHandler(predictionService service.PredictionService) *PredictHandler {
	return &PredictHandler{predictionService: predictionService}
}

func (h *PredictHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := h.predictionService.PredicElectricityUsages(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "prediction successful"})
}

func (api *API) FetchAllPrediction(w http.ResponseWriter, r *http.Request) {
	prediction, err := api.predictionService.FetchAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prediction)
}
