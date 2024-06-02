package service

import (
	"finalcourseproject/model"
	"finalcourseproject/repository"
)

type ElectricityUsagesService interface {
	FetchAll() ([]model.ElectricityUsages, error)
	FetchByID(id int) (*model.ElectricityUsages, error)
	Store(s *model.ElectricityUsages) error
	Update(id int, s *model.ElectricityUsages) error
	Delete(id int) error
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

func (s *electricityUsagesService) FetchByID(id int) (*model.ElectricityUsages, error) {
	electricityUsages, err := s.electricityUsagesRepository.FetchByID(id)
	if err != nil {
		return nil, err
	}

	return electricityUsages, nil
}

func (s *electricityUsagesService) Store(electricityUsages *model.ElectricityUsages) error {
	err := s.electricityUsagesRepository.Store(electricityUsages)
	if err != nil {
		return err
	}

	return nil
}

func (s *electricityUsagesService) Update(id int, electricityUsages *model.ElectricityUsages) error {
	err := s.electricityUsagesRepository.Update(id, electricityUsages)
	if err != nil {
		return err
	}

	return nil
}

func (s *electricityUsagesService) Delete(id int) error {
	err := s.electricityUsagesRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
