package repository

import (
	"errors"
	"finalcourseproject/model"

	"gorm.io/gorm"
)

type ElectricityUsagesRepository interface {
	FetchAll() ([]model.ElectricityUsages, error)
}

type electricityUsagesRepoImpl struct {
	db *gorm.DB
}

func NewElectricityUsagesRepo(db *gorm.DB) *electricityUsagesRepoImpl {
	return &electricityUsagesRepoImpl{db}
}

func (s *electricityUsagesRepoImpl) FetchAll() ([]model.ElectricityUsages, error) {
	var electricityUsages []model.ElectricityUsages
	if err := s.db.Find(&electricityUsages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no electricity usage found")
		}
		return nil, err
	}
	return electricityUsages, nil // TODO: replace this
}
