package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Health struct {
	FallingObjectBase
}

type HealthProducer struct {
	HealthPacks     []*Health
	health          int
	totalHealthPack int
	Cfg             core.GameConfig
}

func NewHealthProducer(cfg core.GameConfig) *HealthProducer {
	return &HealthProducer{
		HealthPacks:     []*Health{},
		health:          cfg.HealthDropConfig.Health,
		totalHealthPack: cfg.HealthDropConfig.Start,
		Cfg:             cfg,
	}
}

func (h *HealthProducer) GenerateHealthPack() {
	if h.totalHealthPack < h.Cfg.HealthDropConfig.MaxDrop {
		h.totalHealthPack++
	}
}

func (h *HealthProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the stars.
	for _, health := range h.HealthPacks {
		health.move(delta)
	}

	spacehsip := h.MovementAndCollision(delta, gc)

	spacehsip.LevelUp(func() {
		h.health += 1
	})
}

func (h *HealthProducer) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	for _, health := range h.HealthPacks {
		healthDropPattern := []struct {
			dx, dy int
			symbol rune
			color  tcell.Style
		}{
			// corners
			{int(health.OriginPoint.X) - health.Width, int(health.OriginPoint.Y) - health.Height, tcell.RuneULCorner, whiteColor},
			{int(health.OriginPoint.X) + health.Width, int(health.OriginPoint.Y) - health.Height, tcell.RuneURCorner, whiteColor},
			{int(health.OriginPoint.X) + health.Width, int(health.OriginPoint.Y) + health.Height, tcell.RuneLRCorner, whiteColor},
			{int(health.OriginPoint.X) - health.Width, int(health.OriginPoint.Y) + health.Height, tcell.RuneLLCorner, whiteColor},
		}
		for _, corner := range healthDropPattern {
			window.SetContentWithStyle(corner.dx, corner.dy, corner.symbol, corner.color)
		}
		// lines
		// top line
		for i := range (health.Width * 2) - 1 {
			window.SetContentWithStyle((int(health.OriginPoint.X)-health.Width)+i+1, int(health.OriginPoint.Y)-health.Height, tcell.RuneHLine, whiteColor) // left top
		}
		// sides | health.Height will be changing
		for j := range (health.Height * 2) - 1 {
			for i := range (health.Width * 2) + 1 {
				switch i {
				case 0, (health.Width * 2):
					window.SetContentWithStyle(int(health.OriginPoint.X)-health.Width+i, int(health.OriginPoint.Y)-health.Height+j+1, tcell.RuneVLine, whiteColor) // left top
				default:
					window.SetContentWithStyle(int(health.OriginPoint.X)-health.Width+i, int(health.OriginPoint.Y)-health.Height+j+1, ' ', whiteColor) // left top
				}
			}
		}
		// bottom line
		for i := range (health.Width * 2) - 1 {
			window.SetContentWithStyle((int(health.OriginPoint.X)-health.Width)+i+1, int(health.OriginPoint.Y)+health.Height, tcell.RuneHLine, whiteColor) // left top
		}

		// writing text in the middle of the box
		hpStr := []rune("Health++")
		for i, r := range hpStr {
			window.SetContentWithStyle(int(health.OriginPoint.X)-health.Width+i+1, int(health.OriginPoint.Y)-(health.Height/2), r, whiteColor) // left top
		}
	}
}

func (h *HealthProducer) DeployHealthPack() {
	w, _ := window.GetSize()
	distance := (w - (15 * 2))
	xPos := rand.Intn(distance) + 15
	randSpeed := rand.Intn(max(h.Cfg.HealthDropConfig.Speed, 3)) + 2
	h.HealthPacks = append(h.HealthPacks, &Health{
		FallingObjectBase: FallingObjectBase{
			Speed:       randSpeed,
			Health:      h.health,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       5,
			Height:      1,
		},
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
			spaceship.Health += 1 // increase spaceship health by one.
		}
		if !health.isDead() && !health.isOffScreen(h) {
			activeHealthPacks = append(activeHealthPacks, health)
		}
	}
	h.HealthPacks = activeHealthPacks
	return spaceship
}

func (h *HealthProducer) UIHealthPackData(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	healthStr := []rune(fmt.Sprintf("* Health Packs: %d/%d", h.totalHealthPack, h.Cfg.HealthDropConfig.MaxDrop))
	for i, r := range healthStr {
		window.SetContentWithStyle(2+i, 10, r, whiteColor)
	}
}

func (h *HealthProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'h' || ev.Rune() == 'H' {
			if h.totalHealthPack > 0 {
				h.DeployHealthPack()
				h.totalHealthPack--
			}
		}
	}
}

func (h *HealthProducer) GetType() string {
	return "health"
}
