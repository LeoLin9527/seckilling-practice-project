package configs

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type Config struct {
	Server struct {
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
	} `yaml:"database"`
	MQ struct {
		Username string `yaml:"user", envconfig:"MQ_USERNAME"`
		Password string `yaml:"pass", envconfig:"MQ_PASSWORD"`
	} `yaml:"mq"`
}

var Cfg Config

func init() {
	Cfg = Config{}
	readFile(&Cfg)
	readEnv(&Cfg)
	fmt.Printf("%+v", Cfg)

}

func readEnv(c *Config) {
	err := envconfig.Process("", c)
	if err != nil {
		processError(err)
	}
}

func readFile(c *Config) {
	rootPath, err := GetProjectPath()
	if err != nil {
		panic("error while reading config")
	}
	f, err := os.Open(path.Join(rootPath, "configs", "config.yml"))
	if err != nil {
		processError(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
func GetProjectPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}