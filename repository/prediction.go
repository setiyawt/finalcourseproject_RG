package repository

import (
	"errors"
	"finalcourseproject/model"

	"gorm.io/gorm"
)

type PredictionRepository interface {
	Create(prediction model.Prediction) error
	FetchAll() ([]model.Prediction, error)
}

type predictionRepoImpl struct {
	db *gorm.DB
}

func NewPredictionRepo(db *gorm.DB) *predictionRepoImpl {
	return &predictionRepoImpl{db}
}

func (s *predictionRepoImpl) Create(prediction model.Prediction) error {
	return s.db.Create(&prediction).Error
}

func (s *predictionRepoImpl) FetchAll() ([]model.Prediction, error) {
	var prediction []model.Prediction
	if err := s.db.Find(&prediction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no electricity usage found")
		}
		return nil, err
	}
	return prediction, nil
}
