package game

import (
	"encoding/json"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/assets"
)

type Design struct {
	Name         string   `json:"name"`
	Shape        []string `json:"shape"`
	EntityHealth int      `json:"health"`
	Color        string   `json:"color"`
}

type ModifierDesign struct {
	Design
	ModifyHealth            int  `json:"modify_health"`
	ModifyLevel             bool `json:"modify_level"`
	ModifyGunPower          int  `json:"modify_gun_power"`
	ModifyGunCap            int  `json:"modify_gun_cap"`
	ModifyGunSpeed          int  `json:"modify_gun_speed"`
	ModifyGunCoolDown       int  `json:"modify_gun_cooldown"`
	ModifyGunReloadCoolDown int  `json:"modify_gun_reload_cooldown"`
}

type SpaceshipDesign struct {
	Design
	GunPower          int `json:"gun_power"`
	GunSpeed          int `json:"gun_speed"`
	GunCap            int `json:"gun_cap"`
	GunCooldown       int `json:"gun_cooldown"`
	GunReloadCooldown int `json:"gun_reload_cooldown"`
}

type AlienshipDesign struct {
	SpaceshipDesign
	Speed int `json:"speed"`
}

type Designable interface {
	GetColor() tcell.Color
	GetHealth() int
	GetName() string
	GetShape() []string
}

func (d *Design) GetColor() tcell.Color {
	return HexToColor(d.Color)
}

func (d *Design) GetHealth() int {
	return d.EntityHealth
}

func (d *Design) GetShape() []string {
	return d.Shape
}

func (d *Design) GetName() string {
	return d.Name
}

func LoadAsset[T any](filePath string) (T, error) {
	file, err := assets.Files.Open(filePath)
	var design T
	if err != nil {
		return design, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&design); err != nil {
		return design, err
	}

	return design, nil
}

func LoadListOfAssets[T any](filePath string) ([]T, error) {
	file, err := assets.Files.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []T
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func HexToColor(hex string) tcell.Color {
	if len(hex) != 6 {
		return tcell.ColorWhite
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
