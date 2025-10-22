package game

import (
	"log"

	"github.com/BurntSushi/toml"
)

var IsDebug = false

const defaultConfig = `
[dev]
debug = true
fps_counter = false
asteroids = true

[spaceship]
max_level = 50
next_level_score = 100

[stars] 
limit = 10
speed = 50
`

type GameConfig struct {
	SpaceShipConfig struct {
		MaxLevel       int `toml:"max_level"`
		NextLevelScore int `toml:"next_level_score"`
	} `toml:"spaceship"`
	StarsConfig struct {
		Limit int `toml:"limit"`
		Speed int `toml:"speed"`
	} `toml:"stars"`
	Dev struct {
		Debug      bool `toml:"debug"`
		FPSCounter bool `toml:"fps_counter"`
		Asteroids  bool `toml:"asteroids"`
	} `toml:"dev"`
}

func LoadConfig() GameConfig {
	var cfg GameConfig
	if _, err := toml.DecodeFile("config.toml", &cfg); err == nil {
		IsDebug = cfg.Dev.Debug
		return cfg
	}
	if _, err := toml.Decode(defaultConfig, &cfg); err == nil {
		IsDebug = cfg.Dev.Debug
		return cfg
	}

	log.Fatal("Failed to load configuration or invalid defaultConfig")
	return GameConfig{}
}
