package repository

import (
	"errors"
	"finalcourseproject/model"

	"gorm.io/gorm"
)

type PredictionRepository interface {
	FetchAll() ([]model.Prediction, error)
}

type predictionRepoImpl struct {
	db *gorm.DB
}

func NewPredictionRepo(db *gorm.DB) *predictionRepoImpl {
	return &predictionRepoImpl{db}
}

func (s *predictionRepoImpl) FetchAll() ([]model.Prediction, error) {
	var electricityUsages []model.Prediction
	if err := s.db.Find(&electricityUsages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no electricity usage found")
		}
		return nil, err
	}
	return electricityUsages, nil // TODO: replace this
}
