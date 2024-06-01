package repository

import (
	"errors"
	"finalcourseproject/model"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailName(name string) error
	SessionAvailToken(token string) (model.Session, error)
}

type sessionsRepoImpl struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *sessionsRepoImpl {
	return &sessionsRepoImpl{db}
}

func (s *sessionsRepoImpl) AddSessions(session model.Session) error {
	if err := s.db.Create(&session).Error; err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (s *sessionsRepoImpl) DeleteSession(token string) error {
	if err := s.db.Where("token = ?", token).Delete(&model.Session{}).Error; err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	if err := s.db.Model(&model.Session{}).Where("username = ?", session.Username).Updates(&session).Error; err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (s *sessionsRepoImpl) SessionAvailName(name string) error {
	var session model.Session
	if err := s.db.Where("username = ?", name).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("session not found")
		}
		return err
	}
	return nil
}

func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Session{}, err
		}
		return model.Session{}, err
	}
	currentTime := time.Now().UTC()
	if session.Expiry.Before(currentTime) {
		return model.Session{}, errors.New("Token is Expired!")
	}

	return session, nil // TODO: replace this
}
