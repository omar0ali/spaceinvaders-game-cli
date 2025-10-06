package game

import (
	"log"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
[spaceship]
max_level = 50
next_level_score = 100
gun_max_cap = 40
gun_max_speed = 70

[stars] 
limit = 10
speed = 50
`

type GameConfig struct {
	SpaceShipConfig struct {
		MaxLevel       int `toml:"max_level"`
		NextLevelScore int `toml:"next_level_score"`
		GunMaxCap      int `toml:"gun_max_cap"`
		GunMaxSpeed    int `toml:"gun_max_speed"`
	} `toml:"spaceship"`
	StarsConfig struct {
		Limit int `toml:"limit"`
		Speed int `toml:"speed"`
	} `toml:"stars"`
}

func LoadConfig() GameConfig {
	var cfg GameConfig
	if _, err := toml.DecodeFile("config.toml", &cfg); err == nil {
		return cfg
	}
	if _, err := toml.Decode(defaultConfig, &cfg); err == nil {
		return cfg
	}
	log.Fatal("Failed to load configuration or invalid defaultConfig")
	return GameConfig{}
}
