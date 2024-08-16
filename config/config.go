package config

import (
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLogLevel = "DEBUG"
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" valid:"optional"`
	} `json:"log"`

	AWS struct {
		AwsAccessKeyId     string `yaml:"access_key_id" env:"AWS_ACCESS_KEY_ID" valid:"required"`
		AwsSecretAccessKey string `yaml:"secret_access_key" env:"AWS_SECRET_ACCESS_KEY" valid:"required"`
		Region             string `yaml:"aws_region" env:"AWS_REGION" valid:"required"`
	} `yaml:"aws"`

	Kubernetes struct {
		ClusterName string `yaml:"cluster_name" env:"CLUSTER_NAME" valid:"required"`
		Region      string `yaml:"k8s_region" env:"K8S_REGION" valid:"required"`
	} `yaml:"kubernetes"`

	//Slack struct {
	//	URL string `yaml:"url" env:"SLACK_URL" valid:"required"`
	//} `yaml:"slack"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("could not load environment: %v", err)
	}

	return nil
}
