package config

import (
	"github.com/num30/config"
)

type Config struct {
	Dev            bool   `default:"false" envvar:"APP_DEV"`
	DBHost         string `default:"postgres" envvar:"DB_HOST"`
	DBPort         int    `default:"5432" envvar:"DB_PORT"`
	DBName         string `default:"postgres" envvar:"DB_NAME"`
	DBUser         string `default:"postgres" envvar:"DB_USER"`
	DBPassword     string `default:"postgres" envvar:"DB_PASSWORD"`
	KafkaHost      string `default:"kafka" envvar:"KAFKA_HOST"`
	KafkaPort      int    `default:"9092" envvar:"KAFKA_PORT"`
	KafkaTopicName string `default:"currency-rates" envvar:"CURRENCY_KAFKA_TOPIC_NAME"`
}

func New() (*Config, error) {
	var conf Config

	err := config.NewConfReader("config").Read(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
