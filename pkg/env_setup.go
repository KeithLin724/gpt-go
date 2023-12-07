package pkg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvSetUp struct {
	ServerURL    string
	ServerApiURL string
}

// The `Init()` function is a method of the `EnvSetUp` struct. It is responsible for initializing the
// environment variables used in the application.
func (envInit *EnvSetUp) Init() error {
	dockerMode := os.Getenv("DOCKER_MODE")

	// throw error
	if dockerMode == "" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
		fmt.Printf("Docker Mode: false\n")
	} else {
		fmt.Printf("Docker Mode: %s\n", dockerMode)
	}

	// init data
	envInit.ServerURL = os.Getenv("SERVER_URL")
	envInit.ServerApiURL = os.Getenv("SERVER_API_URL")

	return nil
}
