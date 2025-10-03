package core

import (
	"log"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
[spaceship]
max_level = 50
next_level_score = 300
gun_max_cap = 15
gun_max_speed = 75

[aliens]
limit = 1
health = 1
speed = 3
gun_speed = 35
gun_power = 1

[stars] 
limit = 10
speed = 45
`

type GameConfig struct {
	SpaceShipConfig struct {
		MaxLevel       int `toml:"max_level"`
		NextLevelScore int `toml:"next_level_score"`
		GunMaxCap      int `toml:"gun_max_cap"`
		GunMaxSpeed    int `toml:"gun_max_speed"`
	} `toml:"spaceship"`
	AliensConfig struct {
		Limit    int `toml:"limit"`
		Health   int `toml:"health"`
		Speed    int `toml:"speed"`
		GunSpeed int `toml:"gun_speed"`
		GunPower int `toml:"gun_power"`
	} `toml:"aliens"`
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
