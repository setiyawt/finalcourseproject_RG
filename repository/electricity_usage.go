package repository

import (
	"errors"
	"finalcourseproject/model"

	"gorm.io/gorm"
)

type ElectricityUsagesRepository interface {
	FetchAll() ([]model.ElectricityUsages, error)
	FetchByID(id int) (*model.ElectricityUsages, error)
	Store(s *model.ElectricityUsages) error
	Update(id int, s *model.ElectricityUsages) error
	Delete(id int) error
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
		return nil, err
	}
	return electricityUsages, nil
}

func (s *electricityUsagesRepoImpl) Store(electricityUsages *model.ElectricityUsages) error {

	if err := s.db.Save(&electricityUsages); err != nil {
		return nil
	}
	return nil
}

func (s *electricityUsagesRepoImpl) Update(id int, electricityUsages *model.ElectricityUsages) error {

	if err := s.db.Model(&model.ElectricityUsages{}).Where("id = ?", id).Updates(&electricityUsages).Error; err != nil {
		return err
	}
	return nil
}

func (s *electricityUsagesRepoImpl) Delete(id int) error {
	if err := s.db.Where("id = ?", id).Delete(&model.ElectricityUsages{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *electricityUsagesRepoImpl) FetchByID(id int) (*model.ElectricityUsages, error) {
	var electricityUsages model.ElectricityUsages
	if err := s.db.First(&electricityUsages, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Electricity usages not found")
		}
		return nil, err
	}
	return &electricityUsages, nil
}
