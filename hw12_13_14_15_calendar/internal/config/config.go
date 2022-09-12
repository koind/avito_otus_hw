package config

import (
	"os"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type SchedulerConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
}

type SenderConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
}

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

type LoggerConf struct {
	Level    Level
	Filename string
}

type StorageType string

const (
	SQL      StorageType = "sql"
	InMemory StorageType = "in-memory"
)

type StorageConf struct {
	Type StorageType
	Dsn  string
}

type HTTPConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Host string
	Port string
}

type RabbitConf struct {
	Dsn      string
	Exchange string
	Queue    string
}

func NewConfig() Config {
	return Config{}
}

func NewSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{}
}

func NewSenderConfig() SenderConfig {
	return SenderConfig{}
}

func LoadConfig() (*Config, error) {
	return &Config{
		Logger: LoggerConf{
			Level:    Level(os.Getenv("LOG_LEVEL")),
			Filename: os.Getenv("LOG_FILENAME"),
		},
		Storage: StorageConf{
			Type: StorageType(os.Getenv("STORAGE_TYPE")),
			Dsn:  os.Getenv("STORAGE_DSN"),
		},
		HTTP: HTTPConf{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		GRPC: GRPCConf{
			Host: os.Getenv("GRPC_HOST"),
			Port: os.Getenv("GRPC_PORT"),
		},
	}, nil
}

func LoadSchedulerConfig() (*SchedulerConfig, error) {
	return &SchedulerConfig{
		Logger: LoggerConf{
			Level:    Level(os.Getenv("LOG_LEVEL")),
			Filename: os.Getenv("LOG_FILENAME"),
		},
		Storage: StorageConf{
			Type: StorageType(os.Getenv("STORAGE_TYPE")),
			Dsn:  os.Getenv("STORAGE_DSN"),
		},
		Rabbit: RabbitConf{
			Dsn:      os.Getenv("RABBIT_DSN"),
			Exchange: os.Getenv("RABBIT_EXCHANGE"),
			Queue:    os.Getenv("RABBIT_QUEUE"),
		},
	}, nil
}

func LoadSenderConfig() (*SenderConfig, error) {
	return &SenderConfig{
		Logger: LoggerConf{
			Level:    Level(os.Getenv("LOG_LEVEL")),
			Filename: os.Getenv("LOG_FILENAME"),
		},
		Storage: StorageConf{
			Type: StorageType(os.Getenv("STORAGE_TYPE")),
			Dsn:  os.Getenv("STORAGE_DSN"),
		},
		Rabbit: RabbitConf{
			Dsn:      os.Getenv("RABBIT_DSN"),
			Exchange: os.Getenv("RABBIT_EXCHANGE"),
			Queue:    os.Getenv("RABBIT_QUEUE"),
		},
	}, nil
}
