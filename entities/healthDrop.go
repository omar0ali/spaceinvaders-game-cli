package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type Health struct {
	FallingObjectBase
	HealSize int
}

type HealthProducer struct {
	HealthPacks []*Health
	Health      int
}

func (h *HealthProducer) GetType() string {
	return "health"
}

func (h *HealthProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the stars.
	for _, health := range h.HealthPacks {
		health.move(delta)
	}

	spaceship := h.MovementAndCollision(delta, gc)

	spaceship.LevelUp(func() {
		h.Health += 1
	})
}

func (h *HealthProducer) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	width, height := 6, 1
	for _, health := range h.HealthPacks {
		// corners
		window.SetContentWithStyle(int(health.OriginPoint.X)-width, int(health.OriginPoint.Y)-height, tcell.RuneULCorner, whiteColor) // left top
		window.SetContentWithStyle(int(health.OriginPoint.X)+width, int(health.OriginPoint.Y)-height, tcell.RuneURCorner, whiteColor) // right top
		window.SetContentWithStyle(int(health.OriginPoint.X)+width, int(health.OriginPoint.Y)+height, tcell.RuneLRCorner, whiteColor) // bottom right
		window.SetContentWithStyle(int(health.OriginPoint.X)-width, int(health.OriginPoint.Y)+height, tcell.RuneLLCorner, whiteColor) // bottom left
		// lines
		// top line
		for i := range (width * 2) - 1 {
			window.SetContentWithStyle((int(health.OriginPoint.X)-width)+i+1, int(health.OriginPoint.Y)-height, tcell.RuneHLine, whiteColor) // left top
		}
		// sides | height will be changing
		for j := range (height * 2) - 1 {
			for i := range (width * 2) + 1 {
				switch i {
				case 0, (width * 2):
					window.SetContentWithStyle(int(health.OriginPoint.X)-width+i, int(health.OriginPoint.Y)-height+j+1, tcell.RuneVLine, whiteColor) // left top
				default:
					window.SetContentWithStyle(int(health.OriginPoint.X)-width+i, int(health.OriginPoint.Y)-height+j+1, ' ', whiteColor) // left top
				}
			}
		}
		// bottom line
		for i := range (width * 2) - 1 {
			window.SetContentWithStyle((int(health.OriginPoint.X)-width)+i+1, int(health.OriginPoint.Y)+height, tcell.RuneHLine, whiteColor) // left top
		}

		// writing text in the middle of the box
		hpStr := []rune("Health+1")
		for i, r := range hpStr {
			window.SetContentWithStyle(int(health.OriginPoint.X)-width+i+1, int(health.OriginPoint.Y)-(height/2), r, whiteColor) // left top
		}
	}
}

func (h *HealthProducer) DeployHealthPack() {
	w, _ := window.GetSize()
	const padding = 18
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding
	randSpeed := rand.Intn(10) + 3
	h.HealthPacks = append(h.HealthPacks, &Health{
		FallingObjectBase: *NewObject(ObjectOpts{
			Speed:       randSpeed,
			Health:      10,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       6,
			Height:      1,
		}),
		HealSize: 1,
	})
}

func (h *HealthProducer) MovementAndCollision(delta float64, gc *core.GameContext) *SpaceShip {
	var activeHealthPacks []*Health
	var gun *Gun
	var spaceship *SpaceShip
	// similar example in alien.go
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
		gun = &spaceship.Gun
	}

	for _, health := range h.HealthPacks {
		health.move(delta)
		for _, beam := range gun.Beams {
			if health.isHit(&beam.position, gun.Power) {
				spaceship.ScoreHit()
				gun.RemoveBeam(beam)
			}
		}

		_, h := window.GetSize()
		if health.isDead() {
			spaceship.Health += health.HealSize
		}
		if !health.isDead() && !health.isOffScreen(h) {
			activeHealthPacks = append(activeHealthPacks, health)
		}
	}
	h.HealthPacks = activeHealthPacks
	return spaceship
}

func (h *HealthProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'h' { // dev mode
			w, _ := window.GetSize()
			distance := (w - (15 * 2))
			xPos := rand.Intn(distance) + 15
			randSpeed := rand.Intn(10) + 2
			h.HealthPacks = append(h.HealthPacks, &Health{
				FallingObjectBase: *NewObject(ObjectOpts{
					Speed:       randSpeed,
					Health:      10,
					OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
					Width:       6,
					Height:      1,
				}),
				HealSize: 1,
			})
		}
	}
}
