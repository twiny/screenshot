package api

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	Port  string `yaml:"port"`
	Debug bool   `yaml:"debug"`
	//
	Rate   float64 `yaml:"rate"`
	Bursts int     `yaml:"bursts"`
	//
	ImageCache time.Duration `yaml:"image_cache"`
	//
	StorePath  string `yaml:"store"`
	LogPath    string `yaml:"log"`
	ChromePath string `yaml:"chrome_data"`
}

// YAMLConfig
func YAMLConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, err
}
