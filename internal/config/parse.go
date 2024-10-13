package config

import (
	"flag"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Db struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
	ExternalAPI struct {
		BaseURL string `yaml:"base-url"`
		Retries int    `yaml:"retries"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"external-api"`
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, err = replaceEnvVars(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseCLI() (string, error) {
	var path string

	flag.StringVar(&path, "config", "./config.yaml", "path to config file")
	flag.Parse()

	s, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if s.IsDir() {
		return "", err
	}

	return path, nil
}

func replaceEnvVars(input []byte) ([]byte, error) {
	envVarRegexp := regexp.MustCompile(`\$\{(\w+)\}`)
	return envVarRegexp.ReplaceAllFunc(input, func(match []byte) []byte {
		key := string(match[2 : len(match)-1])
		return []byte(os.Getenv(key))
	}), nil
}
