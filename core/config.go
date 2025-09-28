package core

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
[spaceship]
health = 20
max_level = 25
next_level_score = 300
gun_max_cap = 6
gun_cap = 3
gun_speed = 40
gun_max_speed = 80
gun_power = 2

[aliens]
limit = 1
health = 5
speed = 3
gun_speed = 35
gun_power = 1

[stars] 
limit = 15
speed = 50

[health_drop]
health = 6
speed = 3
limit = 1
max_drop = 5
start = 1
`

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
	if _, err := toml.DecodeFile("config.toml", &cfg); err == nil {
		return cfg
	}
	if _, err := toml.Decode(defaultConfig, &cfg); err == nil {
		return cfg
	}
	log.Fatal("Failed to load configuration or invalid defaultConfig")
	return GameConfig{}
}

type Design interface {
	GetShape() []string
	GetName() string
	GetHealth() int
}

func LoadSingleAssetDesign[T Design](filePath string) (T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		var zero T
		return zero, err
	}
	defer file.Close()

	var rawDesign T
	if err := json.NewDecoder(file).Decode(&rawDesign); err != nil {
		var zero T
		return zero, err
	}

	return rawDesign, nil
}

func LoadAssetDesign[T Design](filePath string) ([]T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rawDesigns []T
	if err := json.NewDecoder(file).Decode(&rawDesigns); err != nil {
		return nil, err
	}

	return rawDesigns, nil
}
