package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type MessageQueue struct {
	Host      string `env:"mq_host"`
	Port      int    `env:"mq_port"`
	User      string `env:"mq_user"`
	Pass      string `env:"mq_pass"`
	Queues    Queue
	Exchanges Exchange
}

type Queue struct {
	NewVaccineRegistration string `env:"resident_queue"`
}

type Exchange struct {
	ResidentVaccination string `env:"resident_exchange"`
}

type Database struct {
	Host string `env:"db_host"`
	Port int    `env:"db_port"`
	User string `env:"db_user"`
	Pass string `env:"db_pass"`
	Name string `env:"db_name"`
}

type AppConfig struct {
	MQ MessageQueue
	DB Database
}

func NewConfig() (conf *AppConfig, err error) {
	v := viper.New()

	// handle config path for unit test
	dirPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error get working dir: %s", err))
	}
	dirPaths := strings.Split(dirPath, "/internal")
	godotenv.Load(fmt.Sprintf("%s/.env", dirPaths[0]))

	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	cfg := AppConfig{}
	if err = env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg,nil
}