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
	DbUser     string `mapstructure:"db_user"`
	DbPassword string `mapstructure:"db_password"`
	DbDriver   string `mapstructure:"db_driver"`
	DbName     string `mapstructure:"db_name"`
	DbPort     string `mapstructure:"db_port"`
}

type DbConfig struct {
	envs map[string]*DbEnvs
}

type ServerEnvs struct {
	Path    string `mapstructure:"path"`
	Port    string `mapstructure:"port"`
	Address string
}

type TokenEnvs struct {
	SymmetricKey   string        `mapstructure:"symmetric_key"`
	AccessDuration time.Duration `mapstructure:"access_duration"`
}

type ServerConfig struct {
	Db    DbConfig
	Ports ServerEnvs
	Token TokenEnvs
}

var dbConfig DbConfig
var serverEnvs ServerEnvs
var tokenEnvs TokenEnvs
var Config ServerConfig = ServerConfig{}

func init() {
	dbConfig = DbConfig{}
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(filename)
	configPath = filepath.Clean(filepath.Join(configPath, ".."))
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find config.yaml file: ", err)
	}

	err = viper.UnmarshalKey("db", &dbConfig.envs)
	if err != nil {
		log.Fatal("DB Configs can't be loaded: ", err)
	}
	Config.Db = dbConfig

	err = viper.UnmarshalKey("server", &serverEnvs)
	if err != nil {
		log.Fatal("Server Configs can't be loaded: ", err)
	}
	serverEnvs.Address = serverEnvs.Path + ":" + serverEnvs.Port
	Config.Ports = serverEnvs

	err = viper.UnmarshalKey("token", &tokenEnvs)
	if err != nil {
		log.Fatal("Token Configs can't be loaded: ", err)
	}
	Config.Token = tokenEnvs
}

func GetConfig() ServerConfig {
	return Config
}

func GetDevDbEnvs() DbEnvs {
	return *Config.Db.envs["dev_db"]
}

func GetTestcontainersEnvs() DbEnvs {
	return *Config.Db.envs["testcontainers_db"]
}

func GetDevDbSource() string {
	env := GetDevDbEnvs()
	dbAddr := fmt.Sprintf(
		"%s://%s:%s@localhost:%s/%s?sslmode=disable",
		env.DbDriver,
		env.DbUser,
		env.DbPassword,
		env.DbPort,
		env.DbName,
	)
	return dbAddr
}

func GetServerEnvs() ServerEnvs {
	return Config.Ports
}

func GetTokenEnvs() TokenEnvs {
	return Config.Token
}
