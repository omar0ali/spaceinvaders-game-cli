package core

import (
	"log"

	"github.com/BurntSushi/toml"
)

const (
	// Aliens Objects
	MaxLimitAliens  = 1
	MaxSpeedAliens  = 7
	MaxHealthAliens = 8
	// Health Objects
	MaxLimitHealthDrops = 1
	MaxSpeedHealthDrops = 4
	MaxHealthPack       = 5
	// Stars
	MaxLimitStars = 15
	MaxSpeedStars = 50
)

type GameConfig struct {
	SpaceShipConfig struct {
		Health    int `toml:"health"`
		MaxHealth int `toml:"max_health"`
		MaxLevel  int `toml:"max_level"`
	} `toml:"spaceship"`
	GunConfig struct {
		Speed int `toml:"speed"`
		Limit int `toml:"limit"`
	} `toml:"gun"`
	AliensConfig struct {
		Limit int `toml:"limit"`
	} `toml:"aliens"`
	StarsConfig struct {
		Speed int `toml:"speed"`
		Limit int `toml:"limit"`
	} `toml:"stars"`
	HealthDropConfig struct {
		Limit int `toml:"limit"`
		Speed int `toml:"speed"`
		Max   int `toml:"max"`
		Start int `toml:"start"`
	} `toml:"health_drop"`
}

func LoadConfig() GameConfig {
	var cfg GameConfig
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
