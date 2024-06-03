package main

import (
	"encoding/json"
	"finalcourseproject/api"
	"finalcourseproject/db"
	"finalcourseproject/model"
	repo "finalcourseproject/repository"
	"finalcourseproject/service"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Struktur untuk mendeskripsikan payload MQTT
type UsagePayload struct {
	UsageTime float64   `json:"usage_time"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Kwh       float64   `json:"kwh"`
	Name      string    `json:"name"`
}

func main() {
	// Inisialisasi database
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

	// MQTT Configuration
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("736cb8ed-ba7e-4e1b-800e-ad69a6f90ff5")

	// Menambahkan log tambahan
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("Connected to the MQTT broker successfully.")
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("Connection to the MQTT broker lost: %v\n", err)
	}

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Fatal("Error connecting to the MQTT broker:", token.Error())
		os.Exit(1)
	}

	log.Println("MQTT connection established successfully.")

	// Subscribe to the IoT devices topic
	token = client.Subscribe("kwh", 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

		// Parse the JSON payload
		var usagePayload UsagePayload
		err := json.Unmarshal(msg.Payload(), &usagePayload)
		if err != nil {
			log.Printf("Error parsing JSON payload: %v\n", err)
			return
		}

		usageTime := usagePayload.EndTime.Sub(usagePayload.StartTime).Hours()
		// Log the parsed data
		log.Printf("Parsed data - Name: %s, StartTime: %v, EndTime: %v, Kwh: %f\n", usagePayload.Name, usagePayload.StartTime, usagePayload.EndTime, usagePayload.Kwh)

		// Create a new ElectricityUsages entry

		usage := model.ElectricityUsages{
			Name:      usagePayload.Name,
			StartTime: usagePayload.StartTime,
			EndTime:   usagePayload.EndTime,
			UsageTime: usageTime,
			Kwh:       usagePayload.Kwh,
		}

		// Store the data in the database
		if err := conn.Create(&usage).Error; err != nil {
			log.Printf("Failed to store message to database: %v\n", err)
		} else {
			log.Println("Data stored successfully in the database.")
		}
	})

	token.Wait()

	if token.Error() != nil {
		log.Fatal("Error subscribing to the MQTT topic:", token.Error())
		os.Exit(1)
	} else {
		log.Println("Subscribed to the MQTT topic successfully.")
	}

	// Inisialisasi data prediksi
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
