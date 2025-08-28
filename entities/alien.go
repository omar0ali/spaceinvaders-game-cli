package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type alien struct {
	health   int
	speed    int
	origin   core.Point
	triangle core.Triangle // this will be used to draw the shape of the alien
}

func (a *AlienProducer) AddAlien(health, speed int, origin core.Point) {
	alien := alien{
		health: health,
		speed:  speed,
		origin: origin,
		triangle: core.Triangle{
			A: core.Point{
				X: origin.X,
				Y: origin.Y + 2,
			},
			B: core.Point{
				X: origin.X + 3,
				Y: origin.Y,
			},
			C: core.Point{
				X: origin.X - 3,
				Y: origin.Y,
			},
		},
	}
	a.aliens = append(a.aliens, &alien)
}

func (a *alien) IsHit(point core.Point) bool {
	if point == a.triangle.A ||
		point == a.triangle.B ||
		point == a.triangle.C ||
		point.X == a.origin.X+1 && point.Y == a.origin.Y ||
		point.X == a.origin.X-1 && point.Y == a.origin.Y {
		window.SetContent(point.X-1, point.Y+1, 'X')
		window.SetContent(point.X-1, point.Y, 'X')
		window.SetContent(point.X+1, point.Y, 'X')
		window.SetContent(point.X, point.Y+1, 'X')
		window.SetContent(point.X, point.Y-1, 'X')
		window.SetContent(point.X+1, point.Y+1, 'X')
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
	var beams []*Beam
	// look for the spaceship since it has the gun and the number of beams
	for _, entity := range gc.GetEntities() {
		if entity.GetType() == "spaceship" {
			if spaceship, ok := entity.(*SpaceShip); ok {
				beams = spaceship.Gun.Beams
			}
			break
		}
	}
	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.aliens {
		for _, beam := range beams {
			if alien.IsHit(beam.position) {
				alien.health -= beam.power
			}
		}

		// check the health of each alien
		if alien.health > 0 {
			activeAliens = append(activeAliens, alien)
		}
	}

	a.aliens = activeAliens
}

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the aliens.
	// should create alines and place them in random positions on screen.
	a.CheckAliensHealth(gc) // this will ensure to clean up dead aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	for _, alien := range a.aliens {
		// drawing the points
		window.SetContent(alien.triangle.A.X, alien.triangle.A.Y, '*')
		window.SetContent(alien.triangle.B.X, alien.triangle.B.Y, '*')
		window.SetContent(alien.triangle.C.X, alien.triangle.C.Y, '*')
		// lines bellow
		window.SetContent(alien.triangle.C.X+1, alien.triangle.C.Y, '-')
		window.SetContent(alien.triangle.B.X-1, alien.triangle.C.Y, '-')
		window.SetContent(alien.triangle.C.X+2, alien.triangle.C.Y, '-')
		window.SetContent(alien.triangle.B.X-2, alien.triangle.C.Y, '-')
		// lines left
		window.SetContent(alien.triangle.B.X-1, alien.triangle.B.Y+1, '/')
		window.SetContent(alien.triangle.B.X-2, alien.triangle.B.Y+2, '/')
		// lines right
		window.SetContent(alien.triangle.C.X+1, alien.triangle.B.Y+1, '\\')
		window.SetContent(alien.triangle.C.X+2, alien.triangle.B.Y+2, '\\')
	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			// testing spawn an alinen
			w, h := gc.Screen.Size()
			a.AddAlien(300, 10, core.Point{
				X: w / 2,
				Y: (h / 2) - 5,
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
