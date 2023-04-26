package validations

import (
	"github.com/HouseCham/VetMate/config"
)

var Config *config.Config

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}