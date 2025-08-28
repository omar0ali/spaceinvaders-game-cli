package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type alien struct {
	health   int
	speed    int
	origin   core.PointFloat
	triangle core.Triangle
}

func (a *alien) moveForward(distance float64) {
	a.origin.Y += distance
	a.triangle.A.AppendY(distance)
	a.triangle.B.AppendY(distance)
	a.triangle.C.AppendY(distance)
}

func (a *AlienProducer) AddAlien(health, speed int, origin core.PointFloat) {
	alien := alien{
		health: health,
		speed:  speed,
		origin: origin,
		triangle: core.Triangle{
			A: &core.PointFloat{
				X: origin.X,
				Y: origin.Y + 2,
			},
			B: &core.PointFloat{
				X: origin.X + 3,
				Y: origin.Y,
			},
			C: &core.PointFloat{
				X: origin.X - 3,
				Y: origin.Y,
			},
		},
	}
	a.aliens = append(a.aliens, &alien)
}

func (a *alien) IsHit(point core.PointInterface) bool {
	if a.triangle.A.GetY() > point.GetY() &&
		a.triangle.C.GetY() < point.GetY() &&
		(a.triangle.C.GetX()-1 < point.GetX() && a.triangle.B.GetX()+1 > point.GetX()) {

		window.SetContent(int(point.GetX()-1), int(point.GetY()+1), 'X')
		window.SetContent(int(point.GetY()-1), int(point.GetY()), 'X')
		window.SetContent(int(point.GetX()+1), int(point.GetY()), 'X')
		window.SetContent(int(point.GetX()), int(point.GetY()+1), 'X')
		window.SetContent(int(point.GetX()), int(point.GetY()-1), 'X')
		window.SetContent(int(point.GetX()+1), int(point.GetY()+1), 'X')
		return true
	}

	return false
}

type AlienProducer struct {
	aliens []*alien
}

func InitAlienProducer() AlienProducer {
	return AlienProducer{
		aliens: []*alien{},
	}
}

func (a *AlienProducer) CheckAliensHealth(gc *core.GameContext) {
	var activeAliens []*alien
	var gun *Gun
	// look for the spaceship since it has the gun and the number of beams
	for _, entity := range gc.GetEntities() {
		if entity.GetType() == "spaceship" {
			if spaceship, ok := entity.(*SpaceShip); ok {
				gun = &spaceship.Gun
			}
			break
		}
	}
	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.aliens {
		for _, beam := range gun.Beams {
			if alien.IsHit(&beam.position) {
				alien.health -= beam.power
				gun.RemoveBeam(beam)
			}
		}

		// check the alien ship height position
		_, h := window.GetSize()
		if int(alien.origin.Y) >= h-2 {
			alien.health = 0
		}
		// check the health of each alien
		if alien.health > 0 {
			activeAliens = append(activeAliens, alien)
		}
	}

	a.aliens = activeAliens
}

// TODO: should create alines and place them in random positions on screen.

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the aliens.
	for _, alien := range a.aliens {
		distance := float64(alien.speed) * delta
		alien.moveForward(distance)
	}
	a.CheckAliensHealth(gc) // this will ensure to clean up dead aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	for _, alien := range a.aliens {
		// drawing the points
		window.SetContent(int(alien.triangle.A.GetX()), int(alien.triangle.A.GetY()), '*')
		window.SetContent(int(alien.triangle.B.GetX()), int(alien.triangle.B.GetY()), '*')
		window.SetContent(int(alien.triangle.C.GetX()), int(alien.triangle.C.GetY()), '*')
		// lines bellow
		window.SetContent(int(alien.triangle.C.GetX()+1), int(alien.triangle.C.GetY()), '-')

		window.SetContent(int(alien.triangle.B.GetX()-1), int(alien.triangle.C.GetY()), '-')
		window.SetContent(int(alien.triangle.C.GetX()+2), int(alien.triangle.C.GetY()), '-')
		window.SetContent(int(alien.triangle.B.GetX()-2), int(alien.triangle.C.GetY()), '-')
		// lines left
		window.SetContent(int(alien.triangle.B.GetX()-1), int(alien.triangle.B.GetY()+1), '/')
		window.SetContent(int(alien.triangle.B.GetX()-2), int(alien.triangle.B.GetY()+2), '/')
		// lines right
		window.SetContent(int(alien.triangle.C.GetX()+1), int(alien.triangle.B.GetY()+1), '\\')
		window.SetContent(int(alien.triangle.C.GetX()+2), int(alien.triangle.B.GetY()+2), '\\')
	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			// testing spawn an alinen
			w, _ := gc.Screen.Size()
			a.AddAlien(200, 8, core.PointFloat{
				X: float64(w / 2),
				Y: (0) - 5,
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
