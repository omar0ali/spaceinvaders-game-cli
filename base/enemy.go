package base

import (
	"math/rand"

	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type Enemy struct {
	FallingObjectBase
	Gun
	game.AlienshipDesign
}

func Deploy(fileDesigns string, level float64) *Enemy {
	w, _ := GetSize()
	const padding = 23
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding // starting from 18

	designs, err := game.LoadListOfAssets[game.AlienshipDesign](fileDesigns)
	if err != nil {
		panic(err)
	}

	// pick random design: based on the current health level. The higher the stronger the ships.
	design := designs[rand.Intn(min(int(level)+1, len(designs)))]
	width := len(design.Shape[0])
	height := len(design.Shape)

	// will pick the first alienship as the min or starting point.
	lowest := designs[0].Speed
	randSpeed := rand.Float64()*float64(design.Speed) + float64(lowest) - 1
	alien := &Enemy{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:    design.EntityHealth * int(level),
				MaxHealth: design.EntityHealth * int(level),
				Width:     width,
				Height:    height,
				Position:  PointFloat{X: float64(xPos), Y: -5},
				Speed:     randSpeed,
			},
		},
		Gun: NewGun(
			design.GunCap+int(level)-1,
			design.GunPower+int(level)-1,
			design.GunSpeed+int(level),
			design.GunCooldown-int(level),
			design.GunReloadCooldown-int(level)),
		AlienshipDesign: design,
	}

	return alien
}
