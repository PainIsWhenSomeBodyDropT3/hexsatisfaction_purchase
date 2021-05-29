package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	// Config represents a structure with configs for this microservice.
	Config struct {
		Mongo MongoConfig
		Auth  JWTConfig
		HTTP  HTTPConfig
	}
	// PgConfig represents a structure with configs for mongo database.
	MongoConfig struct {
		URI     string
		Host    string
		Port    int
		Name    string
		Dialect string
	}
	// JWTConfig represents a structure with configs for jwt-token.
	JWTConfig struct {
		SigningKey string
	}
	// HTTPConfig represents a structure with configs for http server.
	HTTPConfig struct {
		Port           int
		MaxHeaderBytes int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
	}
)

// Init populates Config struct with values from config file located at filepath and environment variables.
func Init(path string) (*Config, error) {

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := parseEnv(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(filepath string) error {

	env := ".env"
	envPath, err := os.Getwd()
	if err != nil {
		return err
	}
	envPath = strings.SplitAfter(envPath, "hexsatisfaction_purchase")[0]
	if err := godotenv.Load(envPath + "/" + env); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	configPath := strings.Split(filepath, "/")

	viper.AddConfigPath(envPath + "/" + configPath[0])
	viper.SetConfigName(configPath[1])

	return viper.ReadInConfig()
}

func parseEnv() error {
	if err := parseMongo(); err != nil {
		return err
	}
	return parseJWT()
}

func parseMongo() error {
	viper.SetEnvPrefix("mongo")
	if err := viper.BindEnv("host"); err != nil {
		return err
	}
	return viper.BindEnv("port")
}

func parseJWT() error {
	viper.SetEnvPrefix("jwt")

	return viper.BindEnv("signing_key")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http.port", &cfg.HTTP.Port); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.maxHeaderBytes", &cfg.HTTP.MaxHeaderBytes); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.readTimeout", &cfg.HTTP.ReadTimeout); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.writeTimeout", &cfg.HTTP.WriteTimeout); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mongo.databaseName", &cfg.Mongo.Name); err != nil {
		return err
	}

	return viper.UnmarshalKey("mongo.databaseDialect", &cfg.Mongo.Dialect)
}

func setFromEnv(cfg *Config) {
	cfg.Auth.SigningKey = viper.GetString("signing_key")

	cfg.Mongo.Host = viper.GetString("host")
	cfg.Mongo.Port = viper.GetInt("port")
}
