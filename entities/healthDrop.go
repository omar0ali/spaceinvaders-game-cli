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

	var activeHealthPacks []*Health

	for _, health := range h.HealthPacks {
		// check the star height position
		// clear
		_, h := window.GetSize()
		if !health.isOffScreen(h) {
			activeHealthPacks = append(activeHealthPacks, health)
		}
	}
	h.HealthPacks = activeHealthPacks
}

func (h *HealthProducer) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	width, height := 10, 5
	for _, health := range h.HealthPacks {
		// corners
		window.SetContentWithStyle(int(health.OriginPoint.X)-width, int(health.OriginPoint.Y)-height, '*', whiteColor) // left top
		window.SetContentWithStyle(int(health.OriginPoint.X)+width, int(health.OriginPoint.Y)-height, '*', whiteColor) // right top
		window.SetContentWithStyle(int(health.OriginPoint.X)+width, int(health.OriginPoint.Y)+height, '*', whiteColor) // bottom right
		window.SetContentWithStyle(int(health.OriginPoint.X)-width, int(health.OriginPoint.Y)+height, '*', whiteColor) // bottom left

		// lines
		// left and right lines

	}
}

func (h *HealthProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'h' { // dev mode
			// testing spawn an alinen
			w, _ := window.GetSize()
			// pick a random X position to place the alien ship on screen
			// ----from----------------------------to----// example
			//      15                             85
			// distance = from 15 to width-15 = high - low (85 - 15) = 70
			distance := (w - (15 * 2))
			xPos := rand.Intn(distance) + 15 // starting from 15
			randSpeed := rand.Intn(10) + 2   // start at 2

			// create alien
			h.HealthPacks = append(h.HealthPacks, &Health{
				FallingObjectBase: *NewObject(10, randSpeed, core.PointFloat{X: float64(xPos), Y: -5}),
			})
		}
	}
}
