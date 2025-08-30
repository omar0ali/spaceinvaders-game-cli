package entities

import (
	"fmt"
	"math/rand"

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
	grayColor := window.StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)
	if a.triangle.A.GetY() > point.GetY() &&
		a.triangle.C.GetY()-2 < point.GetY() &&
		(a.triangle.C.GetX()-1 < point.GetX() && a.triangle.B.GetX()+1 > point.GetX()) {

		window.SetContentWithStyle(
			int(point.GetX()-1), int(point.GetY()+1), tcell.RuneBoard, grayColor)
		window.SetContentWithStyle(
			int(point.GetX()-1), int(point.GetY()), tcell.RuneCkBoard, yellowColor)
		window.SetContentWithStyle(
			int(point.GetX()+1), int(point.GetY()), tcell.RuneBoard, grayColor)
		window.SetContentWithStyle(
			int(point.GetX()), int(point.GetY()+1), tcell.RuneCkBoard, redColor)
		window.SetContentWithStyle(
			int(point.GetX()), int(point.GetY()-1), tcell.RuneBoard, yellowColor)
		window.SetContentWithStyle(
			int(point.GetX()+1), int(point.GetY()+1), tcell.RuneCkBoard, grayColor)
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

func (a *AlienProducer) CheckAliensState(gc *core.GameContext) {
	var activeAliens []*alien
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
			if alien.IsHit(&beam.position) {
				alien.health -= beam.power
				gun.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien
		_, h := window.GetSize()
		if alien.health > 0 && int(alien.origin.Y) < h-1 {
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
	a.CheckAliensState(gc) // this will ensure to clean up dead aliens
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	for _, alien := range a.aliens {
		// drawing the points
		// header
		window.SetContentWithStyle(
			int(alien.triangle.A.GetX()), int(alien.triangle.A.GetY()+1), 'v', redColor) // top
		window.SetContentWithStyle(
			int(alien.triangle.A.GetX()+1), int(alien.triangle.A.GetY()+1), '>', redColor)
		window.SetContentWithStyle(
			int(alien.triangle.A.GetX()+2), int(alien.triangle.A.GetY()+1), ']', redColor)
		window.SetContentWithStyle(
			int(alien.triangle.A.GetX()-1), int(alien.triangle.A.GetY()+1), '<', redColor)
		window.SetContentWithStyle(
			int(alien.triangle.A.GetX()-2), int(alien.triangle.A.GetY()+1), '[', redColor)

		window.SetContentWithStyle(
			int(alien.triangle.B.GetX()), int(alien.triangle.B.GetY()), ']', redColor) // right
		window.SetContentWithStyle(
			int(alien.triangle.C.GetX()), int(alien.triangle.C.GetY()), '[', redColor) // left

		// lines bellow
		window.SetContent(int(alien.triangle.C.GetX()+1), int(alien.triangle.C.GetY()), '-')
		window.SetContent(int(alien.triangle.C.GetX()+2), int(alien.triangle.C.GetY()), '-')
		window.SetContent(int(alien.triangle.C.GetX()+3), int(alien.triangle.C.GetY()), '^')
		window.SetContent(int(alien.triangle.B.GetX()-1), int(alien.triangle.C.GetY()), '-')
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
		if ev.Rune() == ' ' { // dev mode
			// testing spawn an alinen
			w, _ := gc.Screen.Size()
			// pick a random X position to place the alien ship on screen
			// ----from----------------------------to----// example
			//      15                             85
			// distance = from 15 to width-15 = high - low (85 - 15) = 70
			distance := (w - (15 * 2))
			xPos := rand.Intn(distance) + 15 // starting from 15
			fmt.Println(xPos, distance, w)
			a.AddAlien(150, 5, core.PointFloat{
				X: float64(xPos),
				Y: (0) - 5,
			})
		}
	}
}

func (a *AlienProducer) GetType() string {
	return "AlienProducer"
}
