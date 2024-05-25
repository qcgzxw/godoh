package server

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Listen   []string `toml:"listen"`
	Upstream string   `toml:"upstream"`
}

func LoadConfig(path string) (config *Config) {
	config = &Config{}
	_, _ = toml.DecodeFile(path, &config)
	fmt.Println(config)
	return
}
