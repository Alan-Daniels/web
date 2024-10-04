package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port     string `yaml:"port" envconfig:"SERVER_PORT"`
		HostName string `yaml:"hostname" envconfig:"SERVER_HOSTNAME"`
	} `yaml:"server"`
	Database struct {
		Namespace string `yaml:"namespace" envconfig:"DB_NAMESPACE"`
		Name      string `yaml:"name" envconfig:"DB_NAME"`
		Username  string `yaml:"username,omitempty" envconfig:"DB_USERNAME"`
		Password  string `yaml:"password,omitempty" envconfig:"DB_PASSWORD"`
		Uri       string `yaml:"uri" envconfig:"DB_URI"`
	} `yaml:"database"`
}

func Init(file string) (*Config, error) {
	cfg := new(Config)
	if err := fileInit(cfg, file); err != nil {
		return nil, err
	}
	//if err := envInit(cfg); err != nil {
	//	return nil, err
	//}

	return cfg, nil
}

func fileInit(cfg *Config, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func envInit(cfg *Config) error {
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}
	return nil
}
