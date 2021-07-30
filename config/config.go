package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	cfgFile = flag.String("config", "./config.yaml", "配置文件路径")

	cfg *Config
)

//Config example config
type Config struct {
	Listen string `yaml:"listen"`
	App    struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Url  string `yaml:"url"`
	} `yaml:"app"`

	Database struct {
		Connection string `yaml:"connection"`
		Host       string `yaml:"host"`
		Database   string `yaml:"database"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
	} `yaml:"database"`

	RabbitMQ struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
	} `yaml:"rabbitMQ"`

	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

//GetConfig 获取配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	bytes, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	cfgData := &Config{}
	err = yaml.Unmarshal(bytes, cfgData)
	if err != nil {
		panic(err)
	}
	cfg = cfgData
	return cfg
}
