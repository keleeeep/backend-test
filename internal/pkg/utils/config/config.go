/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.43
 */

package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Datasource string `yaml:"datasource"`
	} `yaml:"database"`

	SecretKey struct{
		AccessSecret string `yaml:"access_secret"`
	} `yaml:"secret_key"`
}

func New(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml failed: %v", err)
	}

	return &cfg, nil
}
