package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "couldn't parse config file")
	}

	if err := parseEnv(); err != nil {
		return nil, errors.Wrap(err, "couldn't parse env file")
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "couldn't parse config file")
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(filepath string) error {
	envPath, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "couldn't get path to current directory")
	}
	envPath = strings.SplitAfter(envPath, "hexsatisfaction_purchase")[0]
	if err := godotenv.Load(envPath + "/" + ".env"); err != nil {
		return errors.Wrap(err, "couldn't load env file")
	}

	configPath := strings.Split(filepath, "/")

	viper.AddConfigPath(envPath + "/" + configPath[0])
	viper.SetConfigName(configPath[1])

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "couldn't load config file")
	}

	return nil
}

func parseEnv() error {
	if err := parseMongo(); err != nil {
		return errors.Wrap(err, "couldn't parse mongo")
	}

	if err := parseJWT(); err != nil {
		return errors.Wrap(err, "couldn't parse jwt")
	}

	return nil
}

func parseMongo() error {
	viper.SetEnvPrefix("mongo")
	if err := viper.BindEnv("host"); err != nil {
		return errors.Wrap(err, "couldn't bind host")
	}

	if err := viper.BindEnv("port"); err != nil {
		return errors.Wrap(err, "couldn't bind port")
	}

	return nil
}

func parseJWT() error {
	viper.SetEnvPrefix("jwt")
	if err := viper.BindEnv("signing_key"); err != nil {
		return errors.Wrap(err, "couldn't bind signing_key")
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http.port", &cfg.HTTP.Port); err != nil {
		return errors.Wrap(err, "couldn't get port")
	}

	if err := viper.UnmarshalKey("http.maxHeaderBytes", &cfg.HTTP.MaxHeaderBytes); err != nil {
		return errors.Wrap(err, "couldn't get maxHeaderBytes")
	}

	if err := viper.UnmarshalKey("http.readTimeout", &cfg.HTTP.ReadTimeout); err != nil {
		return errors.Wrap(err, "couldn't get readTimeout")
	}

	if err := viper.UnmarshalKey("http.writeTimeout", &cfg.HTTP.WriteTimeout); err != nil {
		return errors.Wrap(err, "couldn't get writeTimeout")
	}

	if err := viper.UnmarshalKey("mongo.databaseName", &cfg.Mongo.Name); err != nil {
		return errors.Wrap(err, "couldn't get databaseName")
	}

	if err := viper.UnmarshalKey("mongo.databaseDialect", &cfg.Mongo.Dialect); err != nil {
		return errors.Wrap(err, "couldn't get databaseDialect")
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Auth.SigningKey = viper.GetString("signing_key")

	cfg.Mongo.Host = viper.GetString("host")
	cfg.Mongo.Port = viper.GetInt("port")
}
