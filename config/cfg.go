package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"time"
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
		Apiurl   string `toml:"apiurl"`
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

func GetTomorrowsDate(timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	now := time.Now().In(loc)
	tomorrow := now.AddDate(0, 0, 1)
	date := tomorrow.Format("2006-01-02")
	return date, nil
}
