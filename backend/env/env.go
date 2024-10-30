package env

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	BASE_URI                    string
	UI_URI                      string
	OAUTH2_GOOGLE_CLIENT_ID     string
	OAUTH2_GOOGLE_CLIENT_SECRET string
	POSTGRES_USER               string
	POSTGRES_PASSWORD           string
	POSTGRES_HOST               string
	POSTGRES_PORT               string
	POSTGRES_DB                 string
	REDIS_ADDR                  string
	MAIL_USER                   string
	MAIL_PASSWORD               string
	MAIL_HOST                   string
	MAIL_PORT                   string
}

const (
	DEV  string = "dev"
	PROD        = "prod"
)

func GetConfig(prefix string) Config {
	envDir := os.Getenv("ENV_DIR")

	if err := godotenv.Load(filepath.Join(envDir, "base.env")); err != nil {
		log.Fatal("Failed to load base env file!")
	}

	if err := godotenv.Load(filepath.Join(envDir, prefix+".env")); err == nil {
		slog.Info("Loaded environment", "EXT_ENVIRONMENT", prefix+".env")
	}

	configData := Config{}
	configStruct := reflect.ValueOf(&configData).Elem()
	types := configStruct.Type()

	for i := 0; i < configStruct.NumField(); i++ {
		configStruct.Field(i).SetString(getEnvOrFail(types.Field(i).Name))
	}

	return configData
}

func getEnvOrFail(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal(key + " not set!")
	}
	return val
}
