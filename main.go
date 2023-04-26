package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/HouseCham/VetMate/auth"
	"github.com/HouseCham/VetMate/config"
	"github.com/HouseCham/VetMate/controllers"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/middleware"
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

	ShareConfigFile(config)

	routes.SetAllRoutes(app)

	log.Printf("Server is running on http://%s:%s", config.DevConfiguration.Server.Host, config.DevConfiguration.Server.Port)
	app.Listen(fmt.Sprintf(":%s", config.DevConfiguration.Server.Port))

}

// LoadDbConnection loads the database connection
// and returns a pointer to it
func LoadDbConnection(config config.Config) (*sql.DB, error) {
	DB, err := sql.Open(config.DevConfiguration.Database.DriverName, config.DevConfiguration.Database.DNS)
	if err != nil {
		return nil, err
	}
	return DB, nil
}

// ShareConfigFile shares the config file with all the packages
// that need it
func ShareConfigFile(config config.Config) {
	controllers.ShareConfigFile(&config)
	validations.ShareConfigFile(&config)
	auth.ShareConfigFile(&config)
	middleware.ShareConfigFile(&config)
	db.ShareConfigFile(&config)
}