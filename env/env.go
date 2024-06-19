package env

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	CONFIG_URL = ""
	REGION     = ""
)

func Init() error {
	_ = godotenv.Load()

	CONFIG_URL = os.Getenv("CONFIG_URL")
	REGION = os.Getenv("REGION")

	return nil
}
