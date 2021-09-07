// package config release work with yaml config file
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	DefaultPath = "" // default prefix for path to real config file
)

// subtype for config, describing dsn parameters
type DSNType struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DBName string `yaml:"dbname"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

// subtype for config, describing constraint parameters
type ConstraintsType struct {
	Dimension struct {
		Max int `yaml:"max"`
	} `yaml:"dimension"`
	Point struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"point"`
}

// describing config structure
type ConfType struct {
	Database    DSNType         `yaml:"db"`
	Constraints ConstraintsType `yaml:"constraints"`
}

/*
open yaml file, read and decode to config structure.
return config instance
*/
func BuildConfig(path string) (cfg *ConfType) {
	var (
		err error

		file    *os.File
		decoder *yaml.Decoder
	)

	// open config file
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

	// decode opened file into config structure
	decoder = yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("error [read config file]:", err)
		return nil
	}

	return cfg
}
