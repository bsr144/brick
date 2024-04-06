package config

import (
	"brick/internal/pkg/utils"
	"log"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func New() *Driver {
	return &Driver{
		Database: loadConfig("database").(Database),
		App:      loadConfig("app").(App),
		Logging:  loadConfig("logging").(Logging),
	}
}

func loadConfig(configType string) interface{} {
	switch strings.ToLower(configType) {
	case "database":
		return loadDatabaseConfig()
	case "app":
		return loadAppConfig()
	case "logging":
		return loadLoggingConfig()
	default:
		log.Fatalf("config type named %s is not found", configType)
	}
	return nil
}

func loadDatabaseConfig() Database {
	dbName, err := utils.ReadStringEnvKey("DB_NAME", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbPort, err := utils.ReadStringEnvKey("DB_PORT", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbHost, err := utils.ReadStringEnvKey("DB_HOST", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbPassword, err := utils.ReadStringEnvKey("DB_PASSWORD", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbUsername, err := utils.ReadStringEnvKey("DB_USERNAME", true)
	if err != nil {
		log.Fatal(err.Error())
	}
	return Database{
		Name:     dbName,
		Port:     dbPort,
		Host:     dbHost,
		Password: dbPassword,
		Username: dbUsername,
	}
}

func loadAppConfig() App {
	appName, err := utils.ReadStringEnvKey("APP_NAME", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	appVersion, err := utils.ReadStringEnvKey("APP_VERSION", false)
	if err != nil {
		log.Fatal(err.Error())
	}

	appPort, err := utils.ReadIntEnvKey("APP_PORT", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	return App{
		Name:    appName,
		Version: appVersion,
		Port:    appPort,
	}
}

func loadLoggingConfig() Logging {
	loggingLevel, err := utils.ReadIntEnvKey("LOGGING_LEVEL", true)
	if err != nil {
		log.Fatal(err.Error())
	}

	return Logging{
		Level: loggingLevel,
	}
}
