package service

import (
	"finalcourseproject/model"
	"finalcourseproject/repository"
	"sync"
)

type PredictionService interface {
	FetchAll() ([]model.Prediction, error)
}

type predictionService struct {
	predictionRepository repository.PredictionRepository
}

func NewPredictionService(predictionRepository repository.PredictionRepository) PredictionService {
	return &predictionService{predictionRepository}
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
