package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"finalcourseproject/model"
	"finalcourseproject/repository"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
)

type PredictionService interface {
	FetchAll() ([]model.Prediction, error)
	PredicElectricityUsages() error
	StartPredictionJob()
}

type predictionService struct {
	predictionRepository        repository.PredictionRepository
	electricityUsagesRepository repository.ElectricityUsagesRepository
}

func NewPredictionService(predictionRepository repository.PredictionRepository, electricityUsagesRepository repository.ElectricityUsagesRepository) PredictionService {
	return &predictionService{
		predictionRepository:        predictionRepository,
		electricityUsagesRepository: electricityUsagesRepository,
	}
}

func (s *predictionService) FetchAll() ([]model.Prediction, error) {
	var predictions []model.Prediction
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		predictions, err = s.predictionRepository.FetchAll()
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return predictions, nil
}

type Payload struct {
	Inputs []float64 `json:"inputs"`
}

type HFResponse struct {
	Predictions []float64 `json:"predictions"`
}

func (s *predictionService) StartPredictionJob() {
	c := cron.New()
	_, err := c.AddFunc("0 0 1 * *", func() { // Setiap tanggal 1 pada pukul 00:00
		err := s.PredicElectricityUsages()
		if err != nil {
			// Tangani error
			log.Println("Error dalam proses prediksi:", err)
		}
	})
	if err != nil {
		// Tangani error
		log.Println("Error dalam menjadwalkan proses prediksi:", err)
	}
	c.Start()
}

func (s *predictionService) PredicElectricityUsages() error {
	usages, err := s.electricityUsagesRepository.FetchAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(usages))
	defer close(errChan)

	for _, usage := range usages {
		wg.Add(1)
		go func(usage model.ElectricityUsages) {
			defer wg.Done()
			data := []float64{usage.Kwh, usage.UsageTime}
			payload := Payload{Inputs: data}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				errChan <- err
				return
			}

			body := bytes.NewReader(payloadBytes)
			apiURL := "https://api-inference.huggingface.co/models/Ankur87/Llama2_Time_series_forecasting"
			req, err := http.NewRequest("POST", apiURL, body)
			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+os.Getenv("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"))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errChan <- err
				return
			}

			var hfResponse HFResponse
			if err := json.Unmarshal(respBody, &hfResponse); err != nil {
				errChan <- err
				return
			}

			if len(hfResponse.Predictions) == 0 {
				errChan <- errors.New("no predictions returned from model")
				return
			}

			predictedKwh := hfResponse.Predictions[0]
			pricePerKwh := 1350.9
			PredictedCost := predictedKwh * pricePerKwh

			prediction := model.Prediction{
				ID:            usage.ID,
				PredictedKwh:  predictedKwh,
				PredictedCost: PredictedCost,
			}

			if err := s.predictionRepository.Create(prediction); err != nil {
				errChan <- err
				return
			}
		}(usage)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
