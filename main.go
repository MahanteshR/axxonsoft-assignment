package main

import (
	"log"
	"net/http"
	"os"

	"github.com/axxonsoft-assignment/pkg/cache"
	tasksHandler "github.com/axxonsoft-assignment/pkg/http/handlers"
	"github.com/axxonsoft-assignment/pkg/http/routes"
	taskService "github.com/axxonsoft-assignment/pkg/service"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	LoadEnv()

	// Initialize redis
	redisClient := NewRedisClient()
	defer redisClient.Close()

	// Initialize layers
	cacheLayer := cache.New(redisClient)
	service := taskService.New(cacheLayer)
	handler := tasksHandler.New(service)

	router := mux.NewRouter()

	// Initialize routes
	routes.New(router, handler)

	port := os.Getenv("HTTP_PORT")
	log.Printf("Server is running on http://localhost:%v\n", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Error running the server: %v", err)
		return
	}
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	return client
}

func LoadEnv() {
	envPath := "./config/.env"

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("No .env file found")
	}
}
