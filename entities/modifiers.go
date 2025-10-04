package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

const MaxHealthKitsToOwn = 5

var (
	ModifierHealth   int     = 5
	IncreaseHealthBy float64 = 3
)

type Health struct {
	base.FallingObjectBase
	game.Design
}

type Modifier struct {
	base.FallingObjectBase
	game.ModifierDesign
}

type Producer struct {
	Modifiers *Modifier
	HealthKit *Health
}

func (p *Producer) Update(gc *game.GameContext, delta float64) {
	if p.Modifiers != nil {
		p.Modifiers.Move(delta)
	}

	if p.HealthKit != nil {
		p.HealthKit.Move(delta)
	}
	p.MovementAndCollision(delta, gc)
}

func (p *Producer) Draw(gc *game.GameContext) {
	if p.HealthKit != nil {
		color := window.StyleIt(tcell.ColorReset, p.HealthKit.GetColor())
		for rowIndex, line := range p.HealthKit.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(p.HealthKit.Position.GetX()) + colIndex
					y := int(p.HealthKit.Position.GetY()) + rowIndex
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
					x := int(p.Modifiers.Position.GetX()) + colIndex
					y := int(p.Modifiers.Position.GetY()) + rowIndex
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

	designs, err := game.LoadListOfAssets[game.ModifierDesign]("modifiers.json")
	if err != nil {
		panic(err)
	}
	design := designs[rand.Intn(len(designs))]

	width := len(design.Shape[0])
	height := len(design.Shape)

	randSpeed := rand.Float64()*float64(4) + 2

	m := &Modifier{
		FallingObjectBase: base.FallingObjectBase{
			ObjectBase: base.ObjectBase{
				Health:    design.EntityHealth + ModifierHealth,
				MaxHealth: design.EntityHealth + ModifierHealth,
				Position:  game.PointFloat{X: float64(xPos), Y: -5},
				Width:     width,
				Height:    height,
				Speed:     randSpeed,
			},
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

	design, err := game.LoadAsset[game.Design]("health_kit.json")
	if err != nil {
		panic(err)
	}

	width := len(design.Shape[0])
	height := len(design.Shape)

	randSpeed := rand.Float64()*float64(4) + 2

	p.HealthKit = &Health{
		FallingObjectBase: base.FallingObjectBase{
			ObjectBase: base.ObjectBase{
				Health:    design.EntityHealth + ModifierHealth,
				MaxHealth: design.EntityHealth + ModifierHealth,
				Position:  game.PointFloat{X: float64(xPos), Y: -5},
				Width:     width,
				Height:    height,
				Speed:     randSpeed,
			},
		},
		Design: design,
	}
}

func (p *Producer) MovementAndCollision(delta float64, gc *game.GameContext) {
	var spaceship *SpaceShip
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	if p.HealthKit != nil {
		if p.HealthKit.movementAndCollision(delta, spaceship) {
			p.HealthKit = nil
		}
	}
	if p.Modifiers != nil {
		if p.Modifiers.movementAndCollision(delta, gc, spaceship) {
			p.Modifiers = nil
		}
	}
}

func (m *Modifier) movementAndCollision(delta float64, gc *game.GameContext, spaceship *SpaceShip) bool {
	_, hight := window.GetSize()
	m.Move(delta)
	for _, beam := range spaceship.GetBeams() {
		if m.IsHit(beam.GetPosition(), spaceship.GetPower()) {
			spaceship.ScoreHit()
			spaceship.RemoveBeam(beam)
		}
	}

	var u *UI
	if ui, ok := gc.FindEntity("ui").(*UI); ok {
		u = ui
	}

	if m.IsDead() {
		spaceship.IncreaseHealth(m.ModifyHealth)
		spaceship.IncreaseGunCap(m.ModifyGunCap, spaceship.cfg.SpaceShipConfig.GunMaxCap)
		spaceship.IncreaseGunPower(m.ModifyGunPower)
		spaceship.IncreaseGunSpeed(m.ModifyGunSpeed, spaceship.cfg.SpaceShipConfig.GunMaxSpeed)
		if m.ModifyLevel {
			SetStatus("Free Upgrade!")
			u.LevelUpScreen = true
			return true
		}
		SetStatus(fmt.Sprintf("Modifier %s Applied!", m.Name))
		return true
	}

	if m.IsOffScreen(hight) {
		return true
	}
	return false
}

func (h *Health) movementAndCollision(delta float64, spaceship *SpaceShip) bool {
	_, hight := window.GetSize()

	h.Move(delta)
	for _, beam := range spaceship.GetBeams() {
		if h.IsHit(beam.GetPosition(), spaceship.GetPower()) {
			spaceship.ScoreHit()
			spaceship.RemoveBeam(beam)
		}
	}

	if h.IsDead() {
		spaceship.healthKitsOwned += 1
		SetStatus("Health: Health kit +1")
		return true
	}

	if h.IsOffScreen(hight) {
		return true
	}
	return false
}

func (p *Producer) InputEvents(event tcell.Event, gc *game.GameContext) {
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
