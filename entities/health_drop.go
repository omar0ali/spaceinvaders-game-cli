package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type HealthDesign struct {
	Name   string   `json:"name"`
	Shape  []string `json:"shape"`
	Health int      `json:"health"`
	Color  string   `json:"color"`
}

func (hd *HealthDesign) GetName() string {
	return hd.Name
}

func (hd *HealthDesign) GetShape() []string {
	return hd.Shape
}

func (hd *HealthDesign) GetHealth() int {
	return hd.Health
}

func (hd *HealthDesign) GetColor() tcell.Color {
	return window.HexToColor(hd.Color)
}

type Health struct {
	FallingObjectBase
	Shape core.Design
}

type HealthProducer struct {
	HealthPacks      []*Health
	health           int
	totalHealthPack  int
	Cfg              core.GameConfig
	healthIncreaseBy int
}

func NewHealthProducer(cfg core.GameConfig, gc *core.GameContext, healthIncreaseBy int) *HealthProducer {
	h := &HealthProducer{
		HealthPacks:      []*Health{},
		health:           cfg.HealthDropConfig.Health,
		totalHealthPack:  cfg.HealthDropConfig.Start,
		Cfg:              cfg,
		healthIncreaseBy: healthIncreaseBy,
	}

	if s, ok := gc.FindEntity("spacehsip").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			h.health += 1
		})
	}

	return h
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

	h.MovementAndCollision(delta, gc)
}

func (h *HealthProducer) Draw(gc *core.GameContext) {
	for _, health := range h.HealthPacks {
		color := window.StyleIt(tcell.ColorReset, health.Shape.GetColor())
		for rowIndex, line := range health.Shape.GetShape() {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(health.OriginPoint.GetX()) + colIndex
					y := int(health.OriginPoint.GetY()) + rowIndex
					window.SetContentWithStyle(x, y, char, color)
				}
			}
		}
	}
}

func (h *HealthProducer) DeployHealthPack() {
	w, _ := window.GetSize()
	distance := (w - (15 * 2))
	xPos := rand.Intn(distance) + 15
	randSpeed := rand.Intn(max(h.Cfg.HealthDropConfig.Speed, 3)) + 2
	design, err := core.LoadAsset[*HealthDesign]("assets/health_pack.json")
	if err != nil {
		panic(err)
	}
	width := len(design.GetShape()[0])
	height := len(design.GetShape())

	h.HealthPacks = append(h.HealthPacks, &Health{
		FallingObjectBase: FallingObjectBase{
			Speed:       randSpeed,
			Health:      design.Health + h.health,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       width,
			Height:      height,
		},
		Shape: design,
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

		_, hight := window.GetSize()
		if health.isDead() {
			if !spaceship.IncreaseHealth(h.healthIncreaseBy) {
				// if spacehsip already full, will return the health pack
				h.totalHealthPack += 1
			}
		}
		if !health.isDead() && !health.isOffScreen(hight) {
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
		if ev.Rune() == 'f' || ev.Rune() == 'F' {
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
