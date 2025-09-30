package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Health struct {
	FallingObjectBase
	core.Design
}

type HealthProducer struct {
	HealthKit        *Health
	totalHealthKits  int
	Cfg              core.GameConfig
	health           float64
	increaseHealthBy float64
}

func NewHealthProducer(cfg core.GameConfig, gc *core.GameContext) *HealthProducer {
	h := &HealthProducer{
		health:           float64(cfg.HealthDropConfig.Health),
		totalHealthKits:  cfg.HealthDropConfig.StartWith,
		Cfg:              cfg,
		increaseHealthBy: 3,
	}

	if s, ok := gc.FindEntity("spacehsip").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			h.health += 0.3
			h.increaseHealthBy += 0.3
		})
	}

	return h
}

func (h *HealthProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the stars.
	if h.HealthKit == nil {
		return
	}
	h.HealthKit.move(delta)

	h.MovementAndCollision(delta, gc)
}

func (h *HealthProducer) Draw(gc *core.GameContext) {
	if h.HealthKit == nil {
		return
	}

	color := window.StyleIt(tcell.ColorReset, h.HealthKit.GetColor())
	for rowIndex, line := range h.HealthKit.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(h.HealthKit.OriginPoint.GetX()) + colIndex
				y := int(h.HealthKit.OriginPoint.GetY()) + rowIndex
				window.SetContentWithStyle(x, y, char, color)
			}
		}
	}
	h.HealthKit.DisplayHealth(6, true, color)
}

func (h *HealthProducer) DeployHealthKit() {
	w, _ := window.GetSize()
	distance := (w - (15 * 2))
	xPos := rand.Intn(distance) + 15
	randSpeed := rand.Intn(max(h.Cfg.HealthDropConfig.Speed, 3)) + 2
	design, err := core.LoadAsset[core.Design]("health_kit.json")
	if err != nil {
		panic(err)
	}

	width := len(design.Shape[0])
	height := len(design.Shape)

	h.HealthKit = &Health{
		FallingObjectBase: FallingObjectBase{
			Speed:       randSpeed,
			Health:      design.EntityHealth + int(h.health),
			MaxHealth:   design.EntityHealth + int(h.health),
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       width,
			Height:      height,
		},
		Design: design,
	}
}

func (h *HealthProducer) MovementAndCollision(delta float64, gc *core.GameContext) {
	if h.HealthKit == nil {
		return
	}

	var gun *Gun
	var spaceship *SpaceShip
	// similar example in alien.go
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
		gun = &spaceship.Gun
	}

	h.HealthKit.move(delta)
	for _, beam := range gun.Beams {
		if h.HealthKit.isHit(&beam.position, gun.Power) {
			spaceship.ScoreHit()
			gun.RemoveBeam(beam)
		}
	}

	_, hight := window.GetSize()

	if h.HealthKit.isDead() {
		h.totalHealthKits += 1
		h.HealthKit = nil
		return
	}

	if h.HealthKit.isOffScreen(hight) {
		h.HealthKit = nil
	}
}

func (h *HealthProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	// This code used for testing

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'y' {
	// 		h.DeployHealthKit()
	// 	}
	// }
}

func (h *HealthProducer) GetType() string {
	return "health"
}
