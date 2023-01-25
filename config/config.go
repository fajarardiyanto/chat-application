package config

import (
	databaseInterfaces "github.com/fajarardiyanto/flt-go-database/interfaces"
	databaseLib "github.com/fajarardiyanto/flt-go-database/lib"
	"github.com/fajarardiyanto/flt-go-logger/interfaces"
	"github.com/fajarardiyanto/flt-go-logger/lib"
	"github.com/fajarardiyanto/flt-go-utils/flags"
)

var (
	cfg      = new(Config)
	logger   interfaces.Logger
	database databaseInterfaces.SQL
	rabbitMQ databaseInterfaces.RabbitMQ
)

type Config struct {
	Version   string `yaml:"version" default:"v1.0.0"`
	Name      string `yaml:"name" default:"App Name"`
	Port      string `yaml:"port" default:"8080"`
	Message   string `yaml:"message" default:"CHAT_MESSAGE"`
	ApiSecret string `yaml:"api_secret" default:"SECRET"`
	Database  struct {
		SQL      databaseInterfaces.SQLConfig              `yaml:"sql"`
		RabbitMQ databaseInterfaces.RabbitMQProviderConfig `yaml:"rabbitmq"`
	} `yaml:"database"`
}

func init() {
	logger = lib.NewLib()
	logger.Init("Testing Chat App")

	flags.Init("config/config.yaml", cfg)
}

func Database(cfg databaseInterfaces.SQLConfig) {
	db := databaseLib.NewLib()
	db.Init(GetLogger())

	InitSQL(db, cfg)
	InitRabbitMQ(db, GetConfig().Database.RabbitMQ)
}

func InitSQL(db databaseInterfaces.Database, cfg databaseInterfaces.SQLConfig) {
	database = db.LoadSQLDatabase(cfg)

	if err := database.LoadSQL(); err != nil {
		logger.Error(err).Quit()
	}
}

func InitRabbitMQ(db databaseInterfaces.Database, cfg databaseInterfaces.RabbitMQProviderConfig) {
	rabbitMQ = db.LoadRabbitMQ(GetConfig().Version, cfg)
}

func GetLogger() interfaces.Logger {
	return logger
}

func GetConfig() *Config {
	return cfg
}

func GetDB() databaseInterfaces.SQL {
	return database
}

func GetRabbitMQ() databaseInterfaces.RabbitMQ {
	return rabbitMQ
}
