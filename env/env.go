package env

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	CONFIG_URL     = ""
	REGION         = ""
	CLICKHOUSE_DSN = ""
)

func Init() error {
	_ = godotenv.Load()

	CONFIG_URL = os.Getenv("CONFIG_URL")
	REGION = os.Getenv("REGION")
	CLICKHOUSE_DSN = os.Getenv("CLICKHOUSE_DSN")

	return nil
}
