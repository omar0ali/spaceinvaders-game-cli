package base

import (
	"math/rand"

	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
)

type Enemy struct {
	FallingObjectBase
	Gun
	design.AlienshipDesign
}

func Deploy(designs []design.AlienshipDesign, level float64, currentShips ...*Enemy) *Enemy {
	w, _ := GetSize()
	const padding = 30

	distance := (w - (padding * 2))

	// choosing the position to place the ship
	var xPos int
	const tolerance = 25 // how much space does it need each ship

	for {
		xPos = rand.Intn(distance) + padding
		overlap := false

		for _, ship := range currentShips {
			start := int(ship.Position.X) - tolerance/2
			end := int(ship.Position.X) + tolerance/2
			if xPos > start && xPos < end {
				overlap = true
				break
			}
		}

		if !overlap {
			break
		}
	}

	// pick random design: based on the current level. The higher the stronger the ships.
	design := designs[rand.Intn(min(int(level)+1, len(designs)))]
	width := len(design.Shape[0])
	height := len(design.Shape)

	// will pick the first alienship as the min or starting point.
	lowest := designs[0].Speed
	randSpeed := rand.Float64()*float64(design.Speed) + float64(lowest) - 1
	enemy := &Enemy{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:    design.EntityHealth * int(level),
				MaxHealth: design.EntityHealth * int(level),
				ObjectEntity: ObjectEntity{
					Width:    width,
					Height:   height,
					Position: PointFloat{X: float64(xPos), Y: -5},
					Speed:    randSpeed,
				},
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

	return enemy
}
