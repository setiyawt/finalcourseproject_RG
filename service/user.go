package service

import (
	"finalcourseproject/model"
	"finalcourseproject/repository"
)

type UserService interface {
	Login(user model.User) error
	Register(user model.User) error
	CheckPassLength(pass string) bool
	CheckPassAlphabet(pass string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Login(user model.User) error {

	err := s.userRepository.CheckAvail(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Register(user model.User) error {
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }
	// user.Password = string(hashedPassword)

	err := s.userRepository.Add(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CheckPassLength(pass string) bool {
	return len(pass) >= 5
}

func (s *userService) CheckPassAlphabet(pass string) bool {
	for _, charVariable := range pass {
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
			return false
		}
	}
	return true
}
