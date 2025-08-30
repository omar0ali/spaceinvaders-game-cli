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

type AlienProducer struct {
	Aliens []*Alien
}

// TODO: should create alines and place them in random positions on screen.

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// limit amount of ships Falling
	if len(a.Aliens) < 3 {
		// testing spawn an alinen
		w, _ := window.GetSize()
		// pick a random X position to place the alien ship on screen
		// ----from----------------------------to----// example
		//      15                             85
		// distance = from 15 to width-15 = high - low (85 - 15) = 70
		distance := (w - (15 * 2))
		xPos := rand.Intn(distance) + 15 // starting from 15
		randSpeed := rand.Intn(8) + 2    // start at 2

		// create alien
		a.Aliens = append(a.Aliens, &Alien{
			FallingObjectBase: *NewObject(100, randSpeed, core.PointFloat{X: float64(xPos), Y: -5}),
		})

	}
	// Update the coordinates of the aliens.
	for _, alien := range a.Aliens {
		distance := float64(alien.Speed) * delta
		alien.move(distance)
	}

	// -------- this will ensure to clean up dead aliens and beams --------

	var activeAliens []*Alien
	var gun *Gun

	// look for the spaceship since it has the gun and the number of beams
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		gun = &spaceship.Gun
	}
	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.Aliens {
		for _, beam := range gun.Beams {
			if alien.isHit(&beam.position, beam.power) {
				gun.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien
		// clear
		_, h := window.GetSize()
		if !alien.isDead() && alien.isOffScreen(h) {
			activeAliens = append(activeAliens, alien)
		}
	}

	a.Aliens = activeAliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	brownColor := window.StyleIt(tcell.ColorReset, tcell.ColorBrown)
	color := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)
	for _, alien := range a.Aliens {
		// drawing the points
		// header
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()), int(alien.TrianglePoint.A.GetY()+1), 'v', brownColor) // top
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()+1), int(alien.TrianglePoint.A.GetY()+1), '>', brownColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()+2), int(alien.TrianglePoint.A.GetY()+1), ']', brownColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()-1), int(alien.TrianglePoint.A.GetY()+1), '<', brownColor)
		window.SetContentWithStyle(
			int(alien.TrianglePoint.A.GetX()-2), int(alien.TrianglePoint.A.GetY()+1), '[', brownColor)

		window.SetContentWithStyle(
			int(alien.TrianglePoint.B.GetX()), int(alien.TrianglePoint.B.GetY()), ']', brownColor) // right
		window.SetContentWithStyle(
			int(alien.TrianglePoint.C.GetX()), int(alien.TrianglePoint.C.GetY()), '[', brownColor) // left

		// lines bellow
		window.SetContentWithStyle(int(alien.TrianglePoint.C.GetX()+1), int(alien.TrianglePoint.C.GetY()), '-', brownColor)
		window.SetContentWithStyle(int(alien.TrianglePoint.C.GetX()+2), int(alien.TrianglePoint.C.GetY()), '-', brownColor)
		window.SetContentWithStyle(int(alien.TrianglePoint.C.GetX()+3), int(alien.TrianglePoint.C.GetY()), '^', brownColor)
		window.SetContentWithStyle(int(alien.TrianglePoint.B.GetX()-1), int(alien.TrianglePoint.C.GetY()), '-', brownColor)
		window.SetContentWithStyle(int(alien.TrianglePoint.B.GetX()-2), int(alien.TrianglePoint.C.GetY()), '-', brownColor)
		// lines left
		window.SetContentWithStyle(int(alien.TrianglePoint.B.GetX()-1), int(alien.TrianglePoint.B.GetY()+1), '/', color)
		window.SetContentWithStyle(int(alien.TrianglePoint.B.GetX()-2), int(alien.TrianglePoint.B.GetY()+2), '/', color)
		// lines right
		window.SetContentWithStyle(int(alien.TrianglePoint.C.GetX()+1), int(alien.TrianglePoint.B.GetY()+1), '\\', color)
		window.SetContentWithStyle(int(alien.TrianglePoint.C.GetX()+2), int(alien.TrianglePoint.B.GetY()+2), '\\', color)
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
			randSpeed := rand.Intn(10) + 2   // start at 2

			// create alien
			a.Aliens = append(a.Aliens, &Alien{
				FallingObjectBase: *NewObject(100, randSpeed, core.PointFloat{X: float64(xPos), Y: -5}),
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "alien"
}
