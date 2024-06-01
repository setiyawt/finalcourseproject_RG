package main

import (
	"finalcourseproject/api"
	"finalcourseproject/db"
	"finalcourseproject/model"
	repo "finalcourseproject/repository"
	"finalcourseproject/service"
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

	conn.AutoMigrate(&model.User{}, &model.Session{}, &model.ElectricityUsages{})

	electricity_usages := []model.ElectricityUsages{
		{
			UsageTime:     time.Now(),
			Kwh:           40.0,
			Price_per_kwh: 415,
			CreatedAt:     time.Now(),
		},
		{
			UsageTime:     time.Now(),
			Kwh:           50.0,
			Price_per_kwh: 415,
			CreatedAt:     time.Now(),
		},
		{
			UsageTime:     time.Now(),
			Kwh:           30.0,
			Price_per_kwh: 415,
			CreatedAt:     time.Now(),
		},
	}

	for _, c := range electricity_usages {
		if err := conn.Create(&c).Error; err != nil {
			panic("failed to create default electricity_usage")
		}
	}

	userRepo := repo.NewUserRepo(conn)
	sessionRepo := repo.NewSessionRepo(conn)
	electricityUsagesRepo := repo.NewElectricityUsagesRepo(conn)
	// studentRepo := repo.NewStudentRepo(conn)
	// classRepo := repo.NewClassRepo(conn)

	userService := service.NewUserService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)
	electricityUsagesService := service.NewElectricityUsagesService(electricityUsagesRepo)
	// studentService := service.NewStudentService(studentRepo)
	// classService := service.NewClassService(classRepo)

	mainAPI := api.NewAPI(userService, sessionService, electricityUsagesService)
	mainAPI.Start()
}
