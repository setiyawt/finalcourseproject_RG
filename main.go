package main

import (
	"finalcourseproject/api"
	"finalcourseproject/db"
	"finalcourseproject/model"
	repo "finalcourseproject/repository"
	"finalcourseproject/service"
	"fmt"
	"sync"
	"time"
)

func main() {
	db := db.NewDB()
	dbCredential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "finalcourseproject",
		Port:         5432,
		Schema:       "public",
	}

	conn, err := db.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.User{}, &model.Session{}, &model.Prediction{}, &model.ElectricityUsages{})

	prediction := []model.Prediction{
		{
			PredictedKwh: 2240.0,
			PredictedAt:  time.Now(),
			CreatedAt:    time.Now(),
		},
		{
			PredictedKwh: 2133.0,
			PredictedAt:  time.Now(),
			CreatedAt:    time.Now(),
		},
		{
			PredictedKwh: 1294.0,
			PredictedAt:  time.Now(),
			CreatedAt:    time.Now(),
		},
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(prediction))

	for _, usage := range prediction {
		wg.Add(1)
		go func(usage model.Prediction) {
			defer wg.Done()
			if err := conn.Create(&usage).Error; err != nil {
				errChan <- fmt.Errorf("failed to create default electricity_usage: %w", err)
			}
		}(usage)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
		}
	}

	userRepo := repo.NewUserRepo(conn)
	sessionRepo := repo.NewSessionRepo(conn)
	predictionRepo := repo.NewPredictionRepo(conn)
	electricityUsagesRepo := repo.NewElectricityUsagesRepo(conn)

	userService := service.NewUserService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)
	predictionService := service.NewPredictionService(predictionRepo)
	electricityUsagesService := service.NewElectricityUsagesService(electricityUsagesRepo)

	mainAPI := api.NewAPI(userService, sessionService, predictionService, electricityUsagesService)
	mainAPI.Start()
}
