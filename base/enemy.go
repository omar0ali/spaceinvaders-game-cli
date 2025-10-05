package base

import (
	"math/rand"

	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Enemy struct {
	FallingObjectBase
	Gun
	game.AlienshipDesign
}

func Deploy(fileDesigns string, level int) *Enemy {
	w, _ := window.GetSize()
	const padding = 20
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding // starting from 18

	designs, err := game.LoadListOfAssets[game.AlienshipDesign](fileDesigns)
	if err != nil {
		panic(err)
	}

	// pick random design: based on the current health level. The higher the stronger the ships.
	design := designs[rand.Intn(min(level+1, len(designs)))]
	width := len(design.Shape[0])
	height := len(design.Shape)

	// TODO: I need to know if its boss fight or alien ship. cuz the speed differ. The boss should be always start from 50

	// will pick the first alienship as the min or starting point.
	lowest := designs[0].Speed
	randSpeed := rand.Float64()*float64(design.Speed) + float64(lowest) - 1
	alien := &Enemy{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:    design.EntityHealth * level,
				MaxHealth: design.EntityHealth * level,
				Width:     width,
				Height:    height,
				Position:  game.PointFloat{X: float64(xPos), Y: -5},
				Speed:     randSpeed,
			},
		},
		Gun:             NewGun(design.GunCap, design.GunPower, design.GunSpeed, design.GunCooldown),
		AlienshipDesign: design,
	}

	return alien
}
