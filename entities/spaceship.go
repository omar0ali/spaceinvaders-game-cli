// Package entities
package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type SpaceShip struct {
	triangle core.Triangle
	origin   core.Point
	Gun      Gun // a spaceship has a gun
}

// init in the bottom center of the secreen by default

func InitSpaceShip() SpaceShip {
	w, h := window.GetSize()
	origin := core.Point{
		X: w / 2,
		Y: h - 2,
	}

	return SpaceShip{
		origin: origin,
		triangle: core.Triangle{
			A: &core.Point{X: origin.X, Y: origin.Y - 1},
			B: &core.Point{X: origin.X - 2, Y: origin.Y + 1}, // left
			C: &core.Point{X: origin.X + 2, Y: origin.Y + 1}, // right
		},
		Gun: Gun{
			Beams: []*Beam{},
		},
	}
}

func (s *SpaceShip) Update(gc *core.GameContext, delta float64) {
	defer s.Gun.Update(gc, delta)
}

func (s *SpaceShip) Draw(gc *core.GameContext) {
	defer s.Gun.Draw(gc)
	gc.Screen.SetContent(int(s.triangle.A.GetX()), int(s.triangle.A.GetY()), '^', nil, window.GetStyle())
	gc.Screen.SetContent(int(s.triangle.B.GetX()), int(s.triangle.B.GetY()), '*', nil, window.GetStyle()) // left
	gc.Screen.SetContent(int(s.triangle.C.GetX()), int(s.triangle.C.GetY()), '*', nil, window.GetStyle()) // right
	// left line
	gc.Screen.SetContent(int(s.triangle.A.GetX())-1, int(s.triangle.A.GetY())+1, '/', nil, window.GetStyle())
	// right line
	gc.Screen.SetContent(int(s.triangle.A.GetX())+1, int(s.triangle.A.GetY())+1, '\\', nil, window.GetStyle())
	// bottom line
	gc.Screen.SetContent(int(s.triangle.A.GetX()), int(s.triangle.A.GetY())+2, tcell.RuneS7, nil, window.GetStyle())
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	defer s.Gun.InputEvents(event, gc)
	moveMouse := func(x int) {
		s.origin.X = x
		s.triangle.A.SetX(float64(x))
		s.triangle.B.SetX(float64(x + 2))
		s.triangle.C.SetX(float64(x - 2))
	}
	switch ev := event.(type) {
	case *tcell.EventMouse:
		x, _ := ev.Position()
		moveMouse(x)
	}
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
