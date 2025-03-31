package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type (
	Config struct {
		Smtp Smtp      `toml:"smtp"`
		Api  Apifutbol `toml:"api"`
	}

	Smtp struct {
		Server   string `toml:"server"`
		Port     int    `toml:"port"`
		Email    string `toml:"email"`
		Password string `toml:"password"`
	}

	Apifutbol struct {
		Apikey   string `toml:"apikey"`
		Leagues  []int  `toml:"leagues"`
		Teams    []int  `toml:"teams"`
		Timezone string `toml:"timezone"`
	}
)

func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err = toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
