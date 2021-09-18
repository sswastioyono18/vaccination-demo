package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Provider the config provider
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	Set(key string, value interface{})
	AutomaticEnv()
}

var defaultConfig *viper.Viper

func init() {
	defaultConfig = readViperConfig()
}

func readViperConfig() *viper.Viper {
	v := viper.New()

	// handle config path for unit test
	dirPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Error get working dir: %s", err))
	}
	dirPaths := strings.Split(dirPath, "/internal")

	godotenv.Load(fmt.Sprintf("%s/params/.env", dirPaths[0]))
	godotenv.Load("./config/.env")

	v.AllowEmptyEnv(true)
	v.AutomaticEnv()
	return v
}

// Config return provider so that you can read config anywhere
func Config() Provider {
	return defaultConfig
}
