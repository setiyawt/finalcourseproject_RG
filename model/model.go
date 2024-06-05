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
	ID        uint      `gorm:"not null"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	UsageTime float64   `gorm:"type:decimal(10,2)"`
	Kwh       float64   `gorm:"not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
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
	Message string `json:"message"`
}
