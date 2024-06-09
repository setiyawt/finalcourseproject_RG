package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"finalcourseproject/model"
	"finalcourseproject/repository"
	"io/ioutil"
	"log"
	"os"

	"net/http"
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
	log.Println("Starting prediction job scheduler")
	c := cron.New()
	_, err := c.AddFunc("31 15 9 6 *", func() { // Setiap tanggal 1 pada pukul 00:00
		log.Println("Running scheduled prediction job")
		err := s.PredicElectricityUsages()
		if err != nil {
			// Tangani error
			log.Println("Error dalam proses prediksi:", err)
		}
	})
	if err != nil {
		// Tangani error
		log.Println("Error dalam menjadwalkan proses prediksi:", err)
		return // exit if cron job scheduling fails
	}
	c.Start()
	log.Println("Prediction job scheduler started")
}

func (s *predictionService) PredicElectricityUsages() error {
	log.Println("Fetching all electricity usages")
	usages, err := s.electricityUsagesRepository.FetchAll()
	if err != nil {
		log.Println("Error fetching electricity usages:", err)
		return err
	}

	token := os.Getenv("HF_API_TOKEN")
	if token == "" {
		log.Println("Hugging Face API token is not set in the environment variables")
		return errors.New("Hugging Face API token is not set")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(usages))
	defer close(errChan)

	for _, usage := range usages {
		wg.Add(1)
		go func(usage model.ElectricityUsages) {
			defer wg.Done()
			log.Println("Processing usage ID:", usage.ID)
			data := []float64{usage.Usage_Kwh}
			payload := Payload{Inputs: data}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Println("Error marshaling payload:", err)
				errChan <- err
				return
			}

			log.Println("Payload to be sent:", string(payloadBytes))
			body := bytes.NewReader(payloadBytes)
			apiURL := "https://api-inference.huggingface.co/models/kashif/autoformer-electricity-hourly"
			req, err := http.NewRequest("POST", apiURL, body)
			if err != nil {
				log.Println("Error creating request:", err)
				errChan <- err
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error sending request:", err)
				errChan <- err
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				errChan <- err
				return
			}

			log.Println("Response from model API:", string(respBody))
			var hfResponse HFResponse
			err = json.Unmarshal(respBody, &hfResponse)
			if err != nil {
				log.Println("Error unmarshaling response:", err)
				errChan <- err
				return
			}

			if len(hfResponse.Predictions) == 0 {
				log.Println("No predictions returned from model")
				errChan <- errors.New("no predictions returned from model")
				return
			}

			predictedKwh := hfResponse.Predictions[0]
			pricePerKwh := 1352.0
			PredictedCost := predictedKwh * pricePerKwh

			prediction := model.Prediction{
				ID:            usage.ID,
				PredictedKwh:  predictedKwh,
				PredictedCost: PredictedCost,
			}

			log.Println("Saving prediction for usage ID:", usage.ID)
			if err := s.predictionRepository.Create(prediction); err != nil {
				log.Println("Error saving prediction:", err)
				errChan <- err
				return
			}
		}(usage)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		log.Println("Error occurred during prediction:", err)
		return err
	default:
		log.Println("All predictions processed successfully")
		return nil
	}
}
