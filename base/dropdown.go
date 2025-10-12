package base

import (
	"math/rand"

	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

const (
	MaxHealthKitsToOwn = 5
	MaxSpeed           = 3
)

type DropDown struct {
	FallingObjectBase
	Design game.Designable
}

func DeployDropDown(design game.Designable, level float64) *DropDown {
	w, _ := GetSize()

	const padding = 20

	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding

	width := len(design.GetShape()[0])
	height := len(design.GetShape())

	speed := rand.Float64()*float64(min(MaxSpeed, int(level))) + 1

	dropdown := &DropDown{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:    design.GetHealth() + int(level),
				MaxHealth: design.GetHealth() + int(level),
				Position:  PointFloat{X: float64(xPos), Y: -5},
				Width:     width,
				Height:    height,
				Speed:     speed,
			},
		},
		Design: design,
	}

	return dropdown
}

func (d *DropDown) MovementAndColision(delta float64, fn func(isDead bool)) {
	_, hight := GetSize()

	if d.IsDead() {
		fn(true)
		return
	}

	if d.IsOffScreen(hight) {
		fn(false)
	}
}
