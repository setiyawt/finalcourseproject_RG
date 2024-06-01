package service

import (
	"finalcourseproject/model"
	"finalcourseproject/repository"
)

type ElectricityUsagesService interface {
	FetchAll() ([]model.ElectricityUsages, error)
}

type electricityUsagesService struct {
	electricityUsagesRepository repository.ElectricityUsagesRepository
}

func NewElectricityUsagesService(electricityUsagesRepository repository.ElectricityUsagesRepository) ElectricityUsagesService {
	return &electricityUsagesService{electricityUsagesRepository}
}

func (s *electricityUsagesService) FetchAll() ([]model.ElectricityUsages, error) {
	electricityUsages, err := s.electricityUsagesRepository.FetchAll()
	if err != nil {
		return nil, err
	}

	return electricityUsages, nil
}
