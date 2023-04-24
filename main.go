package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/HouseCham/VetMate/config"
	"github.com/HouseCham/VetMate/controllers"
	"github.com/HouseCham/VetMate/routes"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB

func main() {
	app := fiber.New()

	config, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	DB, err := LoadDbConnection(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	controllers.ShareDbConnection(DB)
	controllers.ShareConfigFile(&config)
	validations.ShareConfigFile(&config)
	routes.SetAllRoutes(app)

	log.Printf("Server is running on http://%s:%s", config.DevConfiguration.Server.Host, config.DevConfiguration.Server.Port)
	app.Listen(fmt.Sprintf(":%s",config.DevConfiguration.Server.Port))
	
}

func LoadDbConnection(config config.Config) (*sql.DB, error) {
	DB, err := sql.Open(config.DevConfiguration.Database.DriverName, config.DevConfiguration.Database.DNS)
	if err != nil {
		return nil, err
	}
	return DB, nil
}