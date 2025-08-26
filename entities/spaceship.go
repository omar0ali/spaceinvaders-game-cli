// Package entities
package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type SpaceShip struct {
	points struct {
		A, B, C core.Point
	}
	origin core.Point
}

func InitSpaceShip(origin core.Point) *SpaceShip {
	return &SpaceShip{
		origin: origin,
		points: struct {
			A, B, C core.Point
		}{
			A: core.Point{X: origin.X, Y: origin.Y - 1},
			B: core.Point{X: origin.X - 1, Y: origin.Y + 1}, // left
			C: core.Point{X: origin.X + 1, Y: origin.Y + 1}, // right
		},
	}
}

// TODO: will need to use delta to ensure smooth movements

func (s *SpaceShip) Update(gs *core.GameContext) {}

func (s *SpaceShip) Draw(gs *core.GameContext) {
	gs.Screen.SetContent(s.points.A.X, s.points.A.Y, '*', nil, window.GetStyle())
	gs.Screen.SetContent(s.points.B.X, s.points.B.Y, '*', nil, window.GetStyle()) // left
	gs.Screen.SetContent(s.points.C.X, s.points.C.Y, '*', nil, window.GetStyle()) // right
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		screen := gc.Screen
		w, _ := screen.Size()

		// TODO: Refactor

		// restriction
		if s.points.B.X < 1 {
			s.origin.X += 2
			s.points.A.X += 2
			s.points.B.X += 2
			s.points.C.X += 2
			return
		}
		if s.points.C.X > w-1 {
			s.origin.X -= 2
			s.points.A.X -= 2
			s.points.B.X -= 2
			s.points.C.X -= 2
			return
		}

		// controls
		if ev.Key() == tcell.KeyLeft {
			s.origin.X -= 2
			s.points.A.X -= 2
			s.points.B.X -= 2
			s.points.C.X -= 2
		}
		if ev.Key() == tcell.KeyRight {
			s.origin.X += 2
			s.points.A.X += 2
			s.points.B.X += 2
			s.points.C.X += 2
		}
	}
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
