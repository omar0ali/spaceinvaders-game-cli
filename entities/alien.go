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

type AlienProducer struct {
	aliens []*alien
}

func InitAlienProducer() AlienProducer {
	return AlienProducer{
		aliens: []*alien{},
	}
}

func (a *AlienProducer) CheckAliensHealth() {
	var activeAliens []*alien
	for _, a := range a.aliens {
		if a.health > 0 { // only add aliens with high health
			// this should ensure dead aliens are untracked can be gone with GC
			activeAliens = append(activeAliens, a)
		}
	}
	a.aliens = activeAliens
}

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the aliens.
	// should create alines and place them in random positions on screen.
	a.CheckAliensHealth() // this will ensure to clean up dead aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	for _, alien := range a.aliens {
		window.SetContent(alien.triangle.A.X, alien.triangle.A.Y, '*')
		window.SetContent(alien.triangle.B.X, alien.triangle.B.Y, '*')
		window.SetContent(alien.triangle.C.X, alien.triangle.C.Y, '*')
	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			// testing spawn an alinen
			w, h := gc.Screen.Size()
			a.AddAlien(10, 10, core.Point{
				X: w / 2,
				Y: (h / 2) - 5,
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
