package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type DbEnvs struct {
	DbUser        string `mapstructure:"db_user"`
	DbPassword    string `mapstructure:"db_password"`
	ContainerName string `mapstructure:"container_name"`
	DbDriver      string `mapstructure:"db_driver"`
	DbName        string `mapstructure:"db_name"`
	DbPort        string `mapstructure:"db_port"`
}

type ServerEnvs struct {
	Path    string `mapstructure:"path"`
	Port    string `mapstructure:"port"`
	Address string // Derived field.
}

type TokenEnvs struct {
	SymmetricKey   string        `mapstructure:"symmetric_key"`
	AccessDuration time.Duration `mapstructure:"access_duration"`
}

type ServerConfig struct {
	MigrationUrl string `mapstructure:"migration_url"`
	Db           map[string]*DbEnvs
	Server       ServerEnvs
	Token        TokenEnvs
}

var config ServerConfig

func init() {
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Join(filepath.Dir(filename), "..")

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't find config.yaml file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	config.Server.Address = fmt.Sprintf("%s:%s", config.Server.Path, config.Server.Port)
}

// Function to retrieve the entire server configuration.
func GetConfig() ServerConfig {
	return config
}

// Function to retrieve development database environment variables.
func GetDevDbEnvs() DbEnvs {
	return *config.Db["dev_db"]
}

// Function to retrieve test containers database environment variables.
func GetTestcontainersEnvs() DbEnvs {
	return *config.Db["testcontainers_db"]
}

// Function to construct and retrieve the database source connection string.
func GetDevDbSource() string {
	env := GetDevDbEnvs()
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		env.DbDriver,
		env.DbUser,
		env.DbPassword,
		env.ContainerName,
		env.DbPort,
		env.DbName,
	)
}

// Function to retrieve server environment variables.
func GetServerEnvs() ServerEnvs {
	return config.Server
}

// Function to retrieve token environment variables.
func GetTokenEnvs() TokenEnvs {
	return config.Token
}
