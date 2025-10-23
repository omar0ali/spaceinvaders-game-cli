package base

import (
	"math/rand"

	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
)

type DropDown struct {
	FallingObjectBase
	Design design.Designable
}

func DeployDropDown(design design.Designable, level int) *DropDown {
	w, _ := GetSize()

	const padding = 30

	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding

	width := len(design.GetShape()[0])
	height := len(design.GetShape())

	speed := rand.Float64()*float64(min(design.GetMaxSpeed(), level)) + 1

	dropdown := &DropDown{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:    design.GetHealth() + level*2,
				MaxHealth: design.GetHealth() + level*2,
				ObjectEntity: ObjectEntity{
					Position: PointFloat{X: float64(xPos), Y: -5},
					Width:    width,
					Height:   height,
					Speed:    speed,
				},
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
