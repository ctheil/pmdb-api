package main

import (
	"log"

	"github.com/ctheil/pmdb-api/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file... %e", err)
	}
	var a app.App
	a.CreateConnection()
	// a.Migrate()
	a.CreateRoutes()
	a.Run()
}
