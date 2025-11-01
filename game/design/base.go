// Package design
package design

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/game/loader"
)

type Design struct {
	Name         string   `json:"name"`
	Shape        []string `json:"shape"`
	EntityHealth int      `json:"health"`
	Color        string   `json:"color"`
	Speed        int      `json:"speed"`
}

type Designable interface {
	GetColor() tcell.Color
	GetHealth() int
	GetName() string
	GetShape() []string
	GetMaxSpeed() int
}

func (d *Design) GetColor() tcell.Color { return HexToColor(d.Color) }
func (d *Design) GetHealth() int        { return d.EntityHealth }
func (d *Design) GetName() string       { return d.Name }
func (d *Design) GetShape() []string    { return d.Shape }
func (d *Design) GetMaxSpeed() int      { return d.Speed }

type LoadedDesigns struct {
	HealthKitDesign  Design
	ModifierDesign   []ModifierDesign
	ListOfSpaceships []SpaceshipDesign
	ListOfAbilities  []AbilityDesign
	ListOfBossShips  []AlienshipDesign
	ListOfAsteroids  AsteroidDesign
	ListOfAlienships []AlienshipDesign
}

func LoadDesigns() *LoadedDesigns {
	healthKitDesign, err := loader.LoadAsset[Design]("health_kit.json")
	if err != nil {
		panic(err)
	}
	modifierDesigns, err := loader.LoadListOfAssets[ModifierDesign]("modifiers.json")
	if err != nil {
		panic(err)
	}
	listOfSpaceships, err := loader.LoadListOfAssets[SpaceshipDesign]("spaceships.json")
	if err != nil {
		panic(err)
	}
	listOfAbilities, err := loader.LoadListOfAssets[AbilityDesign]("abilities.json")
	if err != nil {
		panic(err)
	}
	listOfBossShips, err := loader.LoadListOfAssets[AlienshipDesign]("bossships.json")
	if err != nil {
		panic(err)
	}
	listOfAsteroids, err := loader.LoadAsset[AsteroidDesign]("asteroids.json")
	if err != nil {
		panic(err)
	}
	listOfAlienships, err := loader.LoadListOfAssets[AlienshipDesign]("alienships.json")
	if err != nil {
		panic(err)
	}

	return &LoadedDesigns{
		HealthKitDesign:  healthKitDesign,
		ModifierDesign:   modifierDesigns,
		ListOfSpaceships: listOfSpaceships,
		ListOfAbilities:  listOfAbilities,
		ListOfBossShips:  listOfBossShips,
		ListOfAsteroids:  listOfAsteroids,
		ListOfAlienships: listOfAlienships,
	}
}
