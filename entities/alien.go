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

func (a *AlienProducer) initAlien(health, speed int, point core.Point) {
	alien := alien{
		health: health,
		speed:  speed,
		origin: point,
		triangle: core.Triangle{
			A: core.Point{
				X: point.X - 1,
				Y: point.Y,
			},
			B: core.Point{
				X: point.X + 1,
				Y: point.Y + 1,
			},
			C: core.Point{
				X: point.X - 1,
				Y: point.Y - 1,
			},
		},
	}
	a.aliens = append(a.aliens, &alien)
}

type AlienProducer struct {
	aliens []*alien
}

func InitAlienProducer() AlienProducer {
	return AlienProducer{
		aliens: []*alien{},
	}
}

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	for _, alien := range a.aliens {
		window.SetContent(alien.triangle.A.X, alien.triangle.A.Y, '*')
		window.SetContent(alien.triangle.B.X, alien.triangle.B.Y, '*')
		window.SetContent(alien.triangle.C.X, alien.triangle.C.Y, '*')
	}
}
func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {}
func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
