package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	KafkaBrokers = "kafka-brokers"
	PgDsn        = "pg-dsn"
)

var conf *Config

type Config struct {
	KafkaBrokers string `yaml:"kafka-brokers"`
	PgDsn        string `yaml:"pg-dsn"`
}

func Init() error {
	ENV := os.Getenv("ENV")

	body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	return err
}

func Get(key string) interface{} {
	switch key {
	case KafkaBrokers:
		return strings.Split(conf.KafkaBrokers, ",")
	case PgDsn:
		return conf.PgDsn
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
