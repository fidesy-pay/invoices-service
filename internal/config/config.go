package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	KafkaBrokers   = "kafka-brokers"
	PgDsn          = "pg-dsn"
	ExpireInterval = "expire-interval"
)

var conf *Config

type Config struct {
	KafkaBrokers   string        `yaml:"kafka-brokers"`
	PgDsn          string        `yaml:"pg-dsn"`
	ExpireInterval time.Duration `yaml:"expire-interval"`
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
	case ExpireInterval:
		return conf.ExpireInterval
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
