package config

import (
	"fmt"
	"os"

	yaml3 "gopkg.in/yaml.v3"
)

type CalendarConf struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type SchedulerConf struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
}

type SenderConf struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
}

type HTTPConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level    string
	Filename string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type RabbitConf struct {
	Dsn      string
	Exchange string
	Queue    string
}

func LoadCalendar(configFile string) (*CalendarConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := new(CalendarConf)
	err = yaml3.Unmarshal(content, newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return newConfig, nil
}

func LoadScheduler(configFile string) (*SchedulerConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := new(SchedulerConf)
	err = yaml3.Unmarshal(content, newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return newConfig, nil
}

func LoadSender(configFile string) (*SenderConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := new(SenderConf)
	err = yaml3.Unmarshal(content, newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return newConfig, nil
}
