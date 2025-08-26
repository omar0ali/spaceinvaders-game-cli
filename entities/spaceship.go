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
			A: core.Point{X: origin.X, Y: origin.Y - 1},
			B: core.Point{X: origin.X - 2, Y: origin.Y + 1}, // left
			C: core.Point{X: origin.X + 2, Y: origin.Y + 1}, // right
		},
	}
}

// TODO: will need to use delta to ensure smooth movements
// delta can be used with window.GetDelta()

func (s *SpaceShip) Update(gs *core.GameContext) {}

func (s *SpaceShip) Draw(gs *core.GameContext) {
	gs.Screen.SetContent(s.triangle.A.X, s.triangle.A.Y, '^', nil, window.GetStyle())
	gs.Screen.SetContent(s.triangle.B.X, s.triangle.B.Y, '*', nil, window.GetStyle()) // left
	gs.Screen.SetContent(s.triangle.C.X, s.triangle.C.Y, '*', nil, window.GetStyle()) // right
	// left line
	gs.Screen.SetContent(s.triangle.A.X-1, s.triangle.A.Y+1, '/', nil, window.GetStyle())
	// right line
	gs.Screen.SetContent(s.triangle.A.X+1, s.triangle.A.Y+1, '\\', nil, window.GetStyle())
	// bottom line
	gs.Screen.SetContent(s.triangle.A.X, s.triangle.A.Y+2, tcell.RuneS7, nil, window.GetStyle())
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	moveMouse := func(x int) {
		s.origin.X = x
		s.triangle.A.X = x
		s.triangle.B.X = x + 2
		s.triangle.C.X = x - 2
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
