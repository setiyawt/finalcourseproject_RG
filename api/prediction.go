package api

import (
	"encoding/json"
	"net/http"
)

func (api *API) FetchAllPrediction(w http.ResponseWriter, r *http.Request) {
	prediction, err := api.predictionService.FetchAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prediction)
}
