package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);unique"`
	Password string `json:"password"`
}
type Session struct {
	gorm.Model
	Token    string    `json:"token"`
	Username string    `json:"username"`
	Expiry   time.Time `json:"expiry"`
}

type ElectricityUsages struct {
	gorm.Model
	ID            uint      `gorm:"not null"`
	Name          string    `gorm:"not null"`
	UsageTime     time.Time `gorm:"not null"`
	Kwh           float64   `gorm:"not null"`
	Price_per_kwh float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

type Prediction struct {
	gorm.Model
	ID           uint      `gorm:"not null"`
	PredictedKwh float64   `gorm:"not null"`
	PredictedAt  time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
