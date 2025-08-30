package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type Alien struct {
	FallingObjectBase
}

func NewAlien(health, speed int, origin core.PointFloat) *Alien {
	return &Alien{
		FallingObjectBase: FallingObjectBase{
			Health:      health,
			Speed:       speed,
			OriginPoint: origin,
			TrianglePoint: core.Triangle{
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
		},
	}
}

type AlienProducer struct {
	aliens []*Alien
}

func (a *AlienProducer) AddAlien(health, speed int, origin core.PointFloat) {
	alien := NewAlien(health, speed, origin)
	a.aliens = append(a.aliens, alien)
}

func InitAlienProducer() AlienProducer {
	return AlienProducer{
		aliens: []*Alien{},
	}
}

func (a *AlienProducer) CheckAliensState(gc *core.GameContext) {
	var activeAliens []*Alien
	var gun *Gun
	// look for the spaceship since it has the gun and the number of beams
	for _, entity := range gc.GetEntities() {
		if spaceship, ok := entity.(*SpaceShip); ok {
			gun = &spaceship.Gun
			break
		}
	}
	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.aliens {
		for _, beam := range gun.Beams {
			if alien.isHit(&beam.position) {
				alien.Health -= beam.power
				gun.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien
		_, h := window.GetSize()
		if alien.Health > 0 && int(alien.OriginPoint.Y) < h-1 {
			activeAliens = append(activeAliens, alien)
		}
	}

	a.aliens = activeAliens
}

// TODO: should create alines and place them in random positions on screen.

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// Update the coordinates of the aliens.
	for _, alien := range a.aliens {
		distance := float64(alien.Speed) * delta
		alien.move(distance)
	}
	a.CheckAliensState(gc) // this will ensure to clean up dead aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	for _, alien := range a.aliens {
		// drawing the points
		// header
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()), int(alien.TrianglePoint.A.GetY()+1), 'v', redColor) // top
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()+1), int(alien.TrianglePoint.A.GetY()+1), '>', redColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()+2), int(alien.TrianglePoint.A.GetY()+1), ']', redColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()-1), int(alien.TrianglePoint.A.GetY()+1), '<', redColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()-2), int(alien.TrianglePoint.A.GetY()+1), '[', redColor)

		window.SetContentWithStyle(
			int(alien.TrianglePoint.B.GetX()), int(alien.TrianglePoint.B.GetY()), ']', redColor) // right
		window.SetContentWithStyle(
			int(alien.TrianglePoint.C.GetX()), int(alien.TrianglePoint.C.GetY()), '[', redColor) // left

		// lines bellow
		window.SetContent(int(alien.TrianglePoint.C.GetX()+1), int(alien.TrianglePoint.C.GetY()), '-')
		window.SetContent(int(alien.TrianglePoint.C.GetX()+2), int(alien.TrianglePoint.C.GetY()), '-')
		window.SetContent(int(alien.TrianglePoint.C.GetX()+3), int(alien.TrianglePoint.C.GetY()), '^')
		window.SetContent(int(alien.TrianglePoint.B.GetX()-1), int(alien.TrianglePoint.C.GetY()), '-')
		window.SetContent(int(alien.TrianglePoint.B.GetX()-2), int(alien.TrianglePoint.C.GetY()), '-')
		// lines left
		window.SetContent(int(alien.TrianglePoint.B.GetX()-1), int(alien.TrianglePoint.B.GetY()+1), '/')
		window.SetContent(int(alien.TrianglePoint.B.GetX()-2), int(alien.TrianglePoint.B.GetY()+2), '/')
		// lines right
		window.SetContent(int(alien.TrianglePoint.C.GetX()+1), int(alien.TrianglePoint.B.GetY()+1), '\\')
		window.SetContent(int(alien.TrianglePoint.C.GetX()+2), int(alien.TrianglePoint.B.GetY()+2), '\\')
	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == ' ' { // dev mode
			// testing spawn an alinen
			w, _ := window.GetSize()
			// pick a random X position to place the alien ship on screen
			// ----from----------------------------to----// example
			//      15                             85
			// distance = from 15 to width-15 = high - low (85 - 15) = 70
			distance := (w - (15 * 2))
			xPos := rand.Intn(distance) + 15 // starting from 15
			a.AddAlien(150, 3, core.PointFloat{
				X: float64(xPos),
				Y: (0) - 5,
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
