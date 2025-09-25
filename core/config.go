package core

import (
	"log"

	"github.com/BurntSushi/toml"
)

type GameConfig struct {
	SpaceShipConfig struct {
		Health         int `toml:"health"`
		MaxLevel       int `toml:"max_level"`
		NextLevelScore int `toml:"next_level_score"`
		GunMaxCap      int `toml:"gun_max_cap"`
		GunCap         int `toml:"gun_cap"`
		GunPower       int `toml:"gun_power"`
		GunMaxSpeed    int `toml:"gun_max_speed"`
		GunSpeed       int `toml:"gun_speed"`
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
