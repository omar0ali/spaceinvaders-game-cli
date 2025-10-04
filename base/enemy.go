package base

import (
	"math/rand"
	"time"

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

	randSpeed := rand.Float64()*float64(design.Speed) + 2
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
		Gun:             NewGun(design.GunCap*level, design.GunPower, design.GunSpeed),
		AlienshipDesign: design,
	}

	done := make(chan struct{})
	go DoEvery(2*time.Second,
		func() {
			alien.InitBeam(game.Point{
				X: int(alien.Position.X) + (alien.Width / 2),
				Y: int(alien.Position.Y) + (alien.Height) + 1,
			}, Down)
		},
		done,
	)

	return alien
}
