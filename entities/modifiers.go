package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

const (
	MaxHealthKitsToOwn = 5
)

type ModifierProducer struct {
	Modifiers *base.DropDown
	HealthKit *base.DropDown
	Level     float64
}

func NewModifierProducer(gc *game.GameContext) *ModifierProducer {
	p := &ModifierProducer{
		Level: 3.0,
	}
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship.OnLevelUp = append(spaceship.OnLevelUp, func(newLevel int) {
			p.Level += 0.1
		})
	}
	return p
}

func (p *ModifierProducer) Update(gc *game.GameContext, delta float64) {
	if nextMinute < minutes {
		if p.HealthKit != nil {
			return
		}
		design, err := game.LoadAsset[game.Design]("health_kit.json")
		if err != nil {
			panic(err)
		}
		p.HealthKit = base.DeployDropDown(&design, p.Level)
		nextMinute++
	}

	if seconds == 20 || seconds == 50 {
		if p.Modifiers != nil {
			return
		}
		designs, err := game.LoadListOfAssets[game.ModifierDesign]("modifiers.json")
		if err != nil {
			panic(err)
		}

		design := designs[rand.Intn(len(designs))]

		p.Modifiers = base.DeployDropDown(&design, p.Level)
	}

	if p.Modifiers != nil {
		Move(&p.Modifiers.ObjectBase, delta)
	}

	if p.HealthKit != nil {
		Move(&p.HealthKit.ObjectBase, delta)
	}

	p.MovementAndCollision(delta, gc)
}

func (p *ModifierProducer) Draw(gc *game.GameContext) {
	if p.HealthKit != nil {
		color := base.StyleIt(p.HealthKit.Design.GetColor())
		for rowIndex, line := range p.HealthKit.Design.GetShape() {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(p.HealthKit.Position.GetX()) + colIndex
					y := int(p.HealthKit.Position.GetY()) + rowIndex
					base.SetContentWithStyle(x, y, char, color)
				}
			}
		}
		p.HealthKit.DisplayHealth(6, true, color, nil)
	}
	if p.Modifiers != nil {
		color := base.StyleIt(p.Modifiers.Design.GetColor())
		for rowIndex, line := range p.Modifiers.Design.GetShape() {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(p.Modifiers.Position.GetX()) + colIndex
					y := int(p.Modifiers.Position.GetY()) + rowIndex
					base.SetContentWithStyle(x, y, char, color)
				}
			}
		}
		p.Modifiers.DisplayHealth(8, true, color, nil)
	}
}

func (p *ModifierProducer) MovementAndCollision(delta float64, gc *game.GameContext) {
	var spaceship *SpaceShip
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	if p.HealthKit != nil {
		Move(&p.HealthKit.ObjectBase, delta)
		for _, beam := range spaceship.GetBeams() {
			if GettingHit(&p.HealthKit.ObjectBase, beam, gc) {
				p.HealthKit.TakeDamage(spaceship.GetPower())
				spaceship.RemoveBeam(beam)
			}
		}

		p.HealthKit.MovementAndColision(delta, func(isDead bool) {
			if isDead {
				if spaceship.healthKitsOwned >= MaxHealthKitsToOwn {
					SetStatus("Health: Health kits maxed out!")
					p.HealthKit = nil
					return
				}
				spaceship.healthKitsOwned += 1
				spaceship.ScoreHit()
				SetStatus("Health: Health kit +1")
			}
			p.HealthKit = nil
		})
	}
	if p.Modifiers != nil {
		Move(&p.Modifiers.ObjectBase, delta)
		for _, beam := range spaceship.GetBeams() {
			if GettingHit(&p.Modifiers.ObjectBase, beam, gc) {
				p.Modifiers.TakeDamage(spaceship.GetPower())
				spaceship.RemoveBeam(beam)
			}
		}

		p.Modifiers.MovementAndColision(delta, func(isDead bool) {
			if isDead {
				spaceship.ScoreHit()
				if m, ok := p.Modifiers.Design.(*game.ModifierDesign); ok {
					spaceship.IncreaseHealth(m.ModifyHealth)
					spaceship.IncreaseGunCap(m.ModifyGunCap, spaceship.cfg.SpaceShipConfig.GunMaxCap)
					spaceship.IncreaseGunPower(m.ModifyGunPower)
					spaceship.IncreaseGunSpeed(m.ModifyGunSpeed, spaceship.cfg.SpaceShipConfig.GunMaxSpeed)
					spaceship.DecreaseCooldown(m.ModifyGunCoolDown)
					spaceship.DecreaseGunReloadCooldown(m.ModifyGunReloadCoolDown)
					if m.ModifyLevel {
						SetStatus("Free Level Up!")
						if u, ok := gc.FindEntity("ui").(*UI); ok {
							u.LevelUpScreen = true
							LevelUpPopUp(gc, u, spaceship)
						}
					} else {
						SetStatus(fmt.Sprintf("Modifier %s Applied!", m.Name))
					}
				}
			}
			p.Modifiers = nil
		})
	}
}

func (p *ModifierProducer) InputEvents(event tcell.Event, gc *game.GameContext) {
	// This code used for testing

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'y' {
	// 		design, err := game.LoadListOfAssets[game.ModifierDesign]("modifiers.json")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		p.Modifiers = base.DeployDropDown(&design[1], p.Level)
	// 	}
	// }
}

func (p *ModifierProducer) GetType() string {
	return "producer"
}
