package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	DefaultPath = ""
)

type DSNType struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DBName string `yaml:"dbname"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

type ConstraintsType struct{
	Dimension struct{
		Max int `yaml:"max"`
	} `yaml:"dimension"`
	Point struct{
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"point"`
}

type Config struct {
	Database DSNType `yaml:"db"`
	Constraints ConstraintsType `yaml:"constraints"`
}

func Conf(path string) (cfg *Config) {
	var (
		err error

		file    *os.File
		decoder *yaml.Decoder
	)

	file, err = os.Open(path + "config.yml")
	if err != nil {
		fmt.Println("error [open config file]:", err)
		return nil
	}

	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println("error [close config file]:", err)
		}
	}()

	decoder = yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("error [read config file]:", err)
		return nil
	}

	return cfg
}
