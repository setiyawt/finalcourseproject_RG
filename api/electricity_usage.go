package api

import (
	"encoding/json"
	"finalcourseproject/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (api *API) FetchAllElectricityUsages(w http.ResponseWriter, r *http.Request) {
	electricityUsages, err := api.electricityUsagesService.FetchAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(electricityUsages)
}

func (api *API) FetchSElectricityUsagesByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	electricityUsages, err := api.electricityUsagesService.FetchByID(idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(electricityUsages)
}

func (api *API) StoreelectricityUsages(w http.ResponseWriter, r *http.Request) {
	var electricityUsages model.ElectricityUsages

	err := json.NewDecoder(r.Body).Decode(&electricityUsages)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}
	log.Printf("Received data: %+v", electricityUsages)
	usage_kwh := electricityUsages.EndTime.Sub(electricityUsages.StartTime).Hours()
	electricityUsages.Usage_Kwh = usage_kwh

	err = api.electricityUsagesService.Store(&electricityUsages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}
	log.Printf("Data to be saved: %+v", electricityUsages)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(electricityUsages)
}

func (api *API) UpdateelectricityUsages(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var electricityUsages model.ElectricityUsages
	err = json.NewDecoder(r.Body).Decode(&electricityUsages)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Updating electricity usage with ID: %d, Data: %+v", idInt, electricityUsages)

	err = api.electricityUsagesService.Update(idInt, &electricityUsages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(electricityUsages)
}

func (api *API) DeletelectricityUsages(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = api.electricityUsagesService.Delete(idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("Electricity usage with ID=%d berhasil dihapus", idInt)
	json.NewEncoder(w).Encode(model.SuccessResponse{Message: message})

}
