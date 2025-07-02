package config

import (
	"base/internal/infrastructures/log"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type MainConfig struct {
	Server   ServerConfig `yaml:"Server"`
	Kafka    KafkaConfig  `yaml:"Kafka"`
	Database DBConfig     `yaml:"Database"`
}

type (
	ServerConfig struct {
		Port            string        `yaml:"Port"`
		GracefulTimeout time.Duration `yaml:"GracefulTimeout"`
		ReadTimeout     time.Duration `yaml:"ReadTimeout"`
		WriteTimeout    time.Duration `yaml:"WriteTimeout"`
		EnableSwagger   bool          `yaml:"EnableSwagger" default:"true"`
	}

	DBConfig struct {
		SlaveDSN        string `yaml:"SlaveDSN"`
		MasterDSN       string `yaml:"MasterDSN"`
		RetryInterval   int    `yaml:"RetryInterval"`
		MaxIdleConn     int    `yaml:"MaxIdleConn"`
		MaxConn         int    `yaml:"MaxConn"`
		ConnMaxLifetime string `yaml:"ConnMaxLifetime"`
	}

	KafkaConfig struct {
		Producer KafkaProducerConfig `yaml:"Producer"`
		Consumer KafkaConsumerConfig `yaml:"Consumer"`
		Topics   KafkaTopics         `yaml:"Topics"`
	}

	KafkaTopics struct {
		Topic1            string `yaml:"Topic1"`
		Topic2            string `yaml:"Topic2"`
		Backoff1ndAttempt string `yaml:"Backoff1ndAttempt"`
	}

	RetryConfig struct {
		MaxRetry          int           `yaml:"MaxRetry"`
		RetryInitialDelay time.Duration `yaml:"RetryInitialDelay"`
		MaxJitter         time.Duration `yaml:"MaxJitter"`
		HandlerTimeout    time.Duration `yaml:"HandlerTimeout"`
		DisableBackoff    bool          `yaml:"DisableBackoff" default:"true"`
	}
	KafkaProducerConfig struct {
		Brokers        []string `yaml:"Brokers"`
		MaxMessageSize int      `yaml:"MaxMessageSize"`
	}

	KafkaConsumerConfig struct {
		Brokers       []string    `yaml:"Brokers"`
		ConsumerGroup string      `yaml:"ConsumerGroup"`
		MaxAttempt    int         `yaml:"MaxAttempt"`
		Retry         RetryConfig `yaml:"Retry"`
	}
)

func ReadModuleConfig(cfg interface{}, module, configLocation string) interface{} {
	if configLocation == "" {
		configLocation = "config/files"
	}

	err := Read(cfg, configLocation, module)
	if err != nil {
		log.WithError(err).Fatalf("failed to read config for %v", module)
	}
	return cfg
}

func Read(cfg interface{}, path string, module string) error {
	getFormatFile := filePath(path)

	switch getFormatFile {
	case ".json":
		fname := path + "/" + module + ".json"
		jsonFile, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonFile, cfg)
	default:
		fname := path + "/" + module + ".yaml"
		yamlFile, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(yamlFile, cfg)
	}

}

func filePath(root string) string {
	var file string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file = filepath.Ext(info.Name())
		return nil
	})
	return file
}
