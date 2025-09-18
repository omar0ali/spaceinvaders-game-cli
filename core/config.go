package core

import (
	"log"

	"github.com/BurntSushi/toml"
)

type GameConfig struct {
	SpaceShipConfig struct {
		Health   int `toml:"health"`
		MaxLevel int `toml:"max_level"`
	} `toml:"spaceship"`
	GunConfig struct {
		Speed int `toml:"speed"`
		Limit int `toml:"limit"`
		Power int `toml:"power"`
	} `toml:"gun"`
	AliensConfig struct {
		Limit    int `toml:"limit"`
		Health   int `toml:"health"`
		Speed    int `toml:"speed"`
		GunSpeed int `toml:"gun_speed"`
		GunPower int `toml"gun_power"`
	} `toml:"aliens"`
	StarsConfig struct {
		Limit int `toml:"limit"`
		Speed int `toml:"speed"`
	} `toml:"stars"`
	HealthDropConfig struct {
		Health  int `tom:"health"`
		Limit   int `toml:"limit"`
		Speed   int `toml:"speed"`
		MaxDrop int `toml:"max_drop"`
		Start   int `toml:"start"`
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
