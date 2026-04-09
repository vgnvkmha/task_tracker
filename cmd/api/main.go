package api

import (
	"log"
	"task_tracker/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
