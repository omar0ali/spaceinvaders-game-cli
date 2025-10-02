package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

const MaxHealthKitsToOwn = 5

var (
	ModifierHealth   int     = 5
	IncreaseHealthBy float64 = 3
)

type Health struct {
	FallingObjectBase
	core.Design
}

type Modifier struct {
	FallingObjectBase
	core.ModifierDesign
}

type Producer struct {
	Modifiers *Modifier
	HealthKit *Health
}

func (p *Producer) Update(gc *core.GameContext, delta float64) {
	if p.Modifiers != nil {
		p.Modifiers.move(delta)
	}

	if p.HealthKit != nil {
		p.HealthKit.move(delta)
	}
	p.MovementAndCollision(delta, gc)
}

func (p *Producer) Draw(gc *core.GameContext) {
	if p.HealthKit != nil {
		color := window.StyleIt(tcell.ColorReset, p.HealthKit.GetColor())
		for rowIndex, line := range p.HealthKit.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(p.HealthKit.OriginPoint.GetX()) + colIndex
					y := int(p.HealthKit.OriginPoint.GetY()) + rowIndex
					window.SetContentWithStyle(x, y, char, color)
				}
			}
		}
		p.HealthKit.DisplayHealth(6, true, color)
	}
	if p.Modifiers != nil {
		color := window.StyleIt(tcell.ColorReset, p.Modifiers.GetColor())
		for rowIndex, line := range p.Modifiers.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(p.Modifiers.OriginPoint.GetX()) + colIndex
					y := int(p.Modifiers.OriginPoint.GetY()) + rowIndex
					window.SetContentWithStyle(x, y, char, color)
				}
			}
		}
		p.Modifiers.DisplayHealth(6, true, color)
	}
}

func (p *Producer) DeployModifiers() {
	if p.Modifiers != nil {
		return
	}
	w, _ := window.GetSize()
	distance := (w - (15 * 2))
	xPos := rand.Intn(distance) + 15
	randSpeed := rand.Intn(4) + 2
	designs, err := core.LoadListOfAssets[core.ModifierDesign]("modifiers.json")
	if err != nil {
		panic(err)
	}
	design := designs[rand.Intn(len(designs))]

	width := len(design.Shape[0])
	height := len(design.Shape)

	m := &Modifier{
		FallingObjectBase: FallingObjectBase{
			Speed:       randSpeed,
			Health:      design.EntityHealth + ModifierHealth,
			MaxHealth:   design.EntityHealth + ModifierHealth,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       width,
			Height:      height,
		},
		ModifierDesign: design,
	}

	p.Modifiers = m
}

func (p *Producer) DeployHealthKit() {
	if p.HealthKit != nil {
		return
	}
	w, _ := window.GetSize()
	distance := (w - (15 * 2))
	xPos := rand.Intn(distance) + 15
	randSpeed := rand.Intn(4) + 2
	design, err := core.LoadAsset[core.Design]("health_kit.json")
	if err != nil {
		panic(err)
	}

	width := len(design.Shape[0])
	height := len(design.Shape)

	p.HealthKit = &Health{
		FallingObjectBase: FallingObjectBase{
			Speed:       randSpeed,
			Health:      design.EntityHealth + ModifierHealth,
			MaxHealth:   design.EntityHealth + ModifierHealth,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       width,
			Height:      height,
		},
		Design: design,
	}
}

func (p *Producer) MovementAndCollision(delta float64, gc *core.GameContext) {
	var spaceship *SpaceShip
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	if p.HealthKit != nil {
		if p.HealthKit.movementAndCollision(delta, gc, spaceship) {
			p.HealthKit = nil
		}
	}
	if p.Modifiers != nil {
		if p.Modifiers.movementAndCollision(delta, gc, spaceship) {
			p.Modifiers = nil
		}
	}
}

func (m *Modifier) movementAndCollision(delta float64, gc *core.GameContext, spaceship *SpaceShip) bool {
	_, hight := window.GetSize()
	m.move(delta)
	for _, beam := range spaceship.Beams {
		if m.isHit(&beam.position, spaceship.Power) {
			spaceship.ScoreHit()
			spaceship.RemoveBeam(beam)
		}
	}

	var u *UI
	if ui, ok := gc.FindEntity("ui").(*UI); ok {
		u = ui
	}

	if m.isDead() {
		spaceship.IncreaseHealth(m.ModifyHealth)
		spaceship.IncreaseGunCap(m.ModifyGunCap)
		spaceship.IncreaseGunPower(m.ModifyGunPower)
		spaceship.IncreaseGunSpeed(m.ModifyGunSpeed)
		if m.ModifyLevel {
			u.SetStatus("Level Up +1")
			u.LevelUpScreen = true
			return true
		}
		u.SetStatus("Modifier Applied!")
		return true
	}

	if m.isOffScreen(hight) {
		return true
	}
	return false
}

func (h *Health) movementAndCollision(delta float64, gc *core.GameContext, spaceship *SpaceShip) bool {
	_, hight := window.GetSize()

	h.move(delta)
	for _, beam := range spaceship.Beams {
		if h.isHit(&beam.position, spaceship.Power) {
			spaceship.ScoreHit()
			spaceship.RemoveBeam(beam)
		}
	}

	if h.isDead() {
		spaceship.healthKitsOwned += 1
		if ui, ok := gc.FindEntity("ui").(*UI); ok {
			ui.SetStatus("Health: Health kit +1")
		}
		return true
	}

	if h.isOffScreen(hight) {
		return true
	}
	return false
}

func (p *Producer) InputEvents(event tcell.Event, gc *core.GameContext) {
	// This code used for testing

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'y' {
	// 		p.DeployModifiers()
	// 	}
	// }
}

func (p *Producer) GetType() string {
	return "producer"
}
