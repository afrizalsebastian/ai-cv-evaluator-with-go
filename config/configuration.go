package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

var (
	ErrOsGetPwd            = errors.New("failed to get working directory")
	ErrReadConfigFile      = errors.New("failed to read config file")
	ErrUnmarshalConfigFile = errors.New("failed to read config file")
)

type Config struct {
	AppPort                    int64    `mapstructure:"PORT"`
	GeminiApiKey               string   `mapstructure:"GEMINI_API_KEY"`
	ChromaUrl                  string   `mapstructure:"CHROMA_URL"`
	GeminiModel                string   `mapstructure:"GEMINI_MODEl"`
	DBUser                     string   `mapstructure:"DB_USER"`
	DBPassword                 string   `mapstructure:"DB_PASSWORD"`
	DBHost                     string   `mapstructure:"DB_HOST"`
	DBPort                     string   `mapstructure:"DB_PORT"`
	DBName                     string   `mapstructure:"DB_NAME"`
	KafkaTLS                   bool     `mapstructure:"KAFKA_TLS"`
	KafkaSASLEnable            bool     `mapstructure:"KAFKA_SASL_ENABLE"`
	KafkaSASLHandshake         bool     `mapstructure:"KAFKA_SASL_HANDSHAKE"`
	KafkaBroker                []string `mapstructure:"KAFKA_BROKER"`
	KafkaMaxRetryPolicy        int      `mapstructure:"KAFKA_MAX_RETRY_POLICY"`
	KafkaCvEvaluatorTopic      string   `mapstructure:"KAFKA_CV_EVALUATOR_TOPIC"`
	KafkaCvEvaluatorTopicGroup string   `mapstructure:"KAFKA_CV_EVALUATOR_TOPIC_GROUP"`
}

var appConfig Config

func Init() error {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory")

		return ErrOsGetPwd
	}
	v.AddConfigPath(wd)

	if err := v.ReadInConfig(); err != nil {
		var errConfigFileNotFound viper.ConfigFileNotFoundError
		if errors.As(err, &errConfigFileNotFound) {
			log.Println("No .env file found")
		} else {
			log.Println("Failed to read .env file")
			return ErrReadConfigFile
		}
	}

	bindEnvs(v, &appConfig)

	if err := v.Unmarshal(&appConfig); err != nil {
		log.Println("error to unmarshal config file")
		return ErrUnmarshalConfigFile
	}

	return nil

}

func bindEnvs(v *viper.Viper, config interface{}) {
	cfgType := reflect.TypeOf(config).Elem()
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		envKey := field.Tag.Get("mapstructure")
		if envKey != "" {
			_ = v.BindEnv(envKey)
		}
	}
}

func Get() *Config { return &appConfig }
