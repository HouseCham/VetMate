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

func main() {
	app := fiber.New()

	// Load the configuration file from the config folder/config.json
	config, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// setting up the database connection, with a proper connection pool
	db, err := setupDB(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer db.Close()

	controllers.ShareDbConnection(db)

	ShareConfigFile(config)

	routes.SetAllRoutes(app)

	log.Printf("Server is running on http://%s:%s", config.DevConfiguration.Server.Host, config.DevConfiguration.Server.Port)
	app.Listen(fmt.Sprintf(":%s", config.DevConfiguration.Server.Port))

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

// LoadDbConnection loads the database connection
// and returns a pointer to it
func setupDB(config config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.DevConfiguration.Database.DriverName, config.DevConfiguration.Database.DNS)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections
    db.SetMaxOpenConns(10)

    // Set the maximum number of idle connections
    db.SetMaxIdleConns(5)

	// Ping the database to ensure a connection is established
    if err := db.Ping(); err != nil {
        return nil, err
    }

	return db, nil
}