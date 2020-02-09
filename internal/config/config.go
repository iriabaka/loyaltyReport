package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	DataSource struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"datasource"`
	Mail struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		SendFrom string `yaml:"send_from"`
	} `yaml:"mail"`
	Reports map[string]struct {
		Query       string   `yaml:"query"`
		Subject     string   `yaml:"subject"`
		SendTo      []string `yaml:"send_to"`
		Text        string   `yaml:"text"`
		WinEncoding bool     `yaml:"win_encoding"`
	} `yaml:"reports"`
}

func GetConfig(path string, log *log.Logger) (Config, error) {
	var config Config

	log.Println("Read settings from configuration file")
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return config, errors.WithStack(err)
	}

	if err := yaml.Unmarshal(bs, &config); err != nil {
		return config, errors.WithStack(err)
	}

	return config, nil
}
