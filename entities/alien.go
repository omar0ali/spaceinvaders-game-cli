package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/ui"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
)

type AlienProducer struct {
	Aliens        []*base.Enemy
	Level         float64
	SelectedAlien *base.Enemy // used to display healthbar for selected ones
	LoadedDesigns *design.LoadedDesigns
}

func NewAlienProducer(gc *game.GameContext, designs *design.LoadedDesigns) *AlienProducer {
	a := &AlienProducer{
		Level:         1.0,
		LoadedDesigns: designs,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			a.Level += 0.1
			// clear screen from aliens when the player levels up
			a.Aliens = nil
		})
	}
	return a
}

func (a *AlienProducer) Update(gc *game.GameContext, delta float64) {
	// temporary fix, stop deploying ships until game starts
	if ui, ok := gc.FindEntity("ui").(*UI); ok {
		if ui.MenuScreen {
			return
		}
	} else {
		return
	}
	// start deploying
	if boss, ok := gc.FindEntity("boss").(*BossProducer); ok {
		// saying if there is a boss alien ship deployed. It should stop alien ships.
		if boss.BossAlien != nil {
			if len(a.Aliens) > -1 {
				a.Aliens = nil
				a.SelectedAlien = nil
			}
			return
		}
	}

	if len(a.Aliens) < int(a.Level) {
		a.Aliens = append(a.Aliens, base.Deploy(a.LoadedDesigns.ListOfAlienships, a.Level, a.Aliens...))
	}

	// go through each alien's gun and shoot
	for _, alien := range a.Aliens {
		alien.Update(gc, delta)
		alien.InitBeam(base.Point{
			X: int(alien.Position.X) + (alien.Width / 2),
			Y: int(alien.Position.Y) + (alien.Height) + 1,
		}, base.Down, gc.Sounds)
	}

	// -------- this will ensure to clean up dead aliens and beams --------
	a.MovementAndCollision(delta, gc)
}

func (a *AlienProducer) Draw(gc *game.GameContext) {
	// display the last alien ship that was shot

	if a.SelectedAlien != nil {
		base.DisplayHealthLeft(
			&a.SelectedAlien.ObjectBase,
			5,
			a.SelectedAlien.Name,
			15,
			true,
			base.StyleIt(a.SelectedAlien.GetColor()),
			&a.SelectedAlien.Gun,
		)
	}

	for _, alien := range a.Aliens {
		color := base.StyleIt(alien.GetColor())
		alien.Draw(gc, alien.GetColor())

		alien.DisplayHealth(11, color, &alien.Gun)

		// draw shape
		for rowIndex, line := range alien.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(alien.Position.GetX()) + colIndex
					y := int(alien.Position.GetY()) + rowIndex
					base.SetContentWithStyle(x, y, char, color)
				}
			}
		}

	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *game.GameContext) {
	// testing code

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'z' { // dev mode
	// 		base.Deploy("alienships.json", int(a.Level))
	// 	}
	// }
}

func (a *AlienProducer) UIAlienShipData(gc *game.GameContext) {
	whiteColor := base.StyleIt(tcell.ColorWhite)
	greenColor := base.StyleIt(tcell.ColorYellowGreen)
	aliensStr := []rune(fmt.Sprintf("Difficulty Level: %d", int(a.Level)))
	w, h := base.GetSize()
	ui.DrawBoxOverlap(base.Point{
		X: w - (len(aliensStr) + 4), Y: h - 5,
	}, len(aliensStr)+4,
		3,
		func(x, y int) {
			for i, r := range aliensStr {
				base.SetContentWithStyle(x+i+2, y+1, r, whiteColor)
			}
		}, greenColor)
}

func (a *AlienProducer) MovementAndCollision(delta float64, gc *game.GameContext) {
	activeAliens := a.Aliens[:0]
	var spaceship *SpaceShip

	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.Aliens {
		// Update the coordinates of the aliens.
		Move(&alien.ObjectBase, delta)
		for _, beam := range spaceship.GetBeams() {
			if GettingHit(&alien.ObjectBase, beam, gc) {
				a.SelectedAlien = alien
				alien.TakeDamage(spaceship.GetPower())
				spaceship.ScoreHit()
				spaceship.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// can collid with a asteroid
		if a, ok := gc.FindEntity("asteroid").(*AsteroidProducer); ok {
			for _, asteroid := range a.Asteroids {
				if Crash(&alien.ObjectBase, &asteroid.ObjectBase, gc) {
					alien.TakeDamage(1)
					asteroid.TakeDamage(3)
				}
			}
		}

		// can collid with a meteroid
		if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			for _, p := range ps.ParticleProducable {
				switch p.(type) {
				case *particles.MeteroidProducer:
					for _, m := range p.GetParticles() {
						if Crash(&alien.ObjectBase, &m.ObjectEntity, gc) {
							alien.TakeDamage(1)
							p.RemoveParticle(m)
						}
					}
				}
			}
		}

		// only if destroyed by the spaceship (player) not an asteroid.
		if alien.IsDead() {
			if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
				ps.AddParticles(
					particles.InitExplosion(8,
						particles.WithDimensions(
							alien.Position.X,
							alien.Position.Y,
							alien.Width,
							alien.Height,
						),
					),
				)
			}
			gc.Sounds.PlaySound("8-bit-explosion-2.mp3", -1)

			a.SelectedAlien = nil
			spaceship.ScoreKill(alien.EntityHealth)
		}

		// check the alien ship height position
		// check the health of each alien

		_, h := base.GetSize()
		if alien.IsOffScreen(h) {
			a.SelectedAlien = nil
			spaceship.TakeDamage(1)
		}
		if !alien.IsDead() && !alien.IsOffScreen(h) { // still flying
			activeAliens = append(activeAliens, alien)
		}
	}
	a.Aliens = activeAliens
}

func (a *AlienProducer) GetType() string {
	return "alien"
}
