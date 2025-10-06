package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type AlienProducer struct {
	Aliens []*base.Enemy
	Level  float64
}

func NewAlienProducer(gc *game.GameContext) *AlienProducer {
	a := &AlienProducer{
		Level: 1.0,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			a.Level += 0.1
		})
	}

	return a
}

func (a *AlienProducer) Update(gc *game.GameContext, delta float64) {
	if boss, ok := gc.FindEntity("boss").(*BossProducer); ok {
		// saying if there is a boss alien ship deployed. It should stop alien ships.
		if boss.BossAlien != nil {
			if len(a.Aliens) > -1 {
				a.Aliens = nil
			}
			return
		}
	}
	if len(a.Aliens) < int(a.Level) {
		a.Aliens = append(a.Aliens, base.Deploy("alienships.json", a.Level))
	}

	// go through each alien's gun and shoot
	for _, alien := range a.Aliens {
		alien.Update(gc, delta)
		alien.InitBeam(base.Point{
			X: int(alien.Position.X) + (alien.Width / 2),
			Y: int(alien.Position.Y) + (alien.Height) + 1,
		}, base.Down)
	}

	// -------- this will ensure to clean up dead aliens and beams --------
	a.MovementAndCollision(delta, gc)
}

func (a *AlienProducer) Draw(gc *game.GameContext) {
	for _, alien := range a.Aliens {
		color := base.StyleIt(tcell.ColorReset, alien.GetColor())
		alien.Draw(gc, alien.GetColor())

		alien.DisplayHealth(6, true, color, &alien.Gun)

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
	w, _ := base.GetSize()
	whiteColor := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	aliensStr := []rune(fmt.Sprintf("Enemy Level: %d * ", int(a.Level)))
	for i, r := range aliensStr {
		base.SetContentWithStyle(w+i-len(aliensStr), 2, r, whiteColor)
	}
}

func (a *AlienProducer) MovementAndCollision(delta float64, gc *game.GameContext) *SpaceShip {
	var activeAliens []*base.Enemy
	var spaceship *SpaceShip

	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.Aliens {
		// Update the coordinates of the aliens.
		base.Move(&alien.ObjectBase, delta)
		for _, beam := range spaceship.GetBeams() {
			if base.GettingHit(&alien.ObjectBase, beam) {
				alien.TakeDamage(spaceship.GetPower())
				spaceship.ScoreHit()
				spaceship.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien

		_, h := base.GetSize()
		if alien.IsOffScreen(h) {
			spaceship.TakeDamage(1)
		}
		if alien.IsDead() {
			ScoreKill()
		}
		if !alien.IsDead() && !alien.IsOffScreen(h) { // still flying
			activeAliens = append(activeAliens, alien)
		}
	}
	a.Aliens = activeAliens
	return spaceship
}

func (a *AlienProducer) GetType() string {
	return "alien"
}
