package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
	"github.com/omar0ali/spaceinvaders-game-cli/game/loader"
)

type ModifierProducer struct {
	Modifiers        *base.DropDown
	HealthKit        *base.DropDown
	Level            float64
	SelectedDropDown *base.DropDown
}

func NewModifierProducer(gc *game.GameContext) *ModifierProducer {
	p := &ModifierProducer{
		Level: 2,
	}
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship.OnLevelUp = append(spaceship.OnLevelUp, func(newLevel int) {
			p.Level += 0.5
		})
	}
	return p
}

func (p *ModifierProducer) Update(gc *game.GameContext, delta float64) {
	if nextMinute < minutes {
		if p.HealthKit != nil {
			return
		}

		design, err := loader.LoadAsset[design.Design]("health_kit.json")
		if err != nil {
			panic(err)
		}

		p.HealthKit = base.DeployDropDown(&design, int(p.Level))
		nextMinute++
	}

	if seconds == 20 || seconds == 50 {
		if p.Modifiers != nil {
			return
		}
		designs, err := loader.LoadListOfAssets[design.ModifierDesign]("modifiers.json")
		if err != nil {
			panic(err)
		}

		design := designs[rand.Intn(len(designs))]

		p.Modifiers = base.DeployDropDown(&design, int(p.Level))
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
	// display the last dropdown that was hit
	if p.SelectedDropDown != nil {
		base.DisplayHealthLeft(
			&p.SelectedDropDown.ObjectBase,
			11,
			p.SelectedDropDown.Design.GetName(),
			15,
			true,
			base.StyleIt(p.SelectedDropDown.Design.GetColor()),
			nil,
		)
	}

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
		p.HealthKit.DisplayHealth(11, color, nil)
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
		p.Modifiers.DisplayHealth(13, color, nil)
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
				p.SelectedDropDown = p.HealthKit
				p.HealthKit.TakeDamage(spaceship.GetPower())
				spaceship.RemoveBeam(beam)
			}
		}

		p.HealthKit.MovementAndColision(delta, func(isDead bool) {
			if isDead {
				if spaceship.HealthKit.HealthKitsOwned >= spaceship.HealthKit.HealthKitLimit {
					SetStatus("Health kits maxed out!")
					p.HealthKit = nil
					return
				}
				spaceship.HealthKit.HealthKitsOwned += 1
				spaceship.ScoreHit()
				SetStatus("Health kit +1")
			}
			if p.SelectedDropDown == p.HealthKit {
				p.SelectedDropDown = nil
			}

			p.HealthKit = nil
		})
	}
	if p.Modifiers != nil {
		Move(&p.Modifiers.ObjectBase, delta)
		for _, beam := range spaceship.GetBeams() {
			if GettingHit(&p.Modifiers.ObjectBase, beam, gc) {
				p.SelectedDropDown = p.Modifiers
				p.Modifiers.TakeDamage(spaceship.GetPower())
				spaceship.RemoveBeam(beam)
			}
		}

		p.Modifiers.MovementAndColision(delta, func(isDead bool) {
			if isDead {
				spaceship.ScoreHit()
				if m, ok := p.Modifiers.Design.(*design.ModifierDesign); ok {
					spaceship.IncreaseHealth(m.ModifyHealth)
					spaceship.IncreaseGunCap(m.ModifyGunCap, m.MaxValue)
					spaceship.IncreaseGunPower(m.ModifyGunPower)
					spaceship.IncreaseGunSpeed(m.ModifyGunSpeed, m.MaxValue)
					spaceship.DecreaseCooldown(m.ModifyGunCoolDown)
					spaceship.DecreaseGunReloadCooldown(m.ModifyGunReloadCoolDown)
					if m.ModifyLevel {
						SetStatus("Free Level Up!")
						if u, ok := gc.FindEntity("ui").(*UI); ok {
							u.LevelUpScreen = true
							spaceship.LevelUpMenu(gc)
						}
					} else {
						SetStatus(fmt.Sprintf("Modifier %s Applied!", m.Name))
					}
				}
			}
			if p.SelectedDropDown == p.Modifiers {
				p.SelectedDropDown = nil
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
