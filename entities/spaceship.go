// Package entities
package entities

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

var previousLevel = 0

type SpaceshipOpts struct {
	SpaceShipHealth int
	GunPower        int
	GunCapacity     int
	GunSpeed        int
}

type SpaceShip struct {
	triangle core.Triangle
	origin   core.Point
	Gun      Gun // a spaceship has a gun
	Health   int
	Score    int
	Kills    int
	Level    int
}

// init in the bottom center of the secreen by default

func InitSpaceShip(opts SpaceshipOpts) *SpaceShip {
	w, h := window.GetSize()
	origin := core.Point{
		X: w / 2,
		Y: h - 3,
	}

	return &SpaceShip{
		origin: origin,
		triangle: core.Triangle{
			A: &core.Point{X: origin.X, Y: origin.Y - 1},
			B: &core.Point{X: origin.X - 2, Y: origin.Y + 1}, // left
			C: &core.Point{X: origin.X + 2, Y: origin.Y + 1}, // right
		},
		Gun: Gun{
			Beams: []*Beam{},
			Cap:   opts.GunCapacity,
			Power: opts.GunPower,
			Speed: opts.GunSpeed,
		},
		Health: opts.SpaceShipHealth,
	}
}

func (s *SpaceShip) Update(gc *core.GameContext, delta float64) {
	defer s.Gun.Update(gc, delta)
	if s.Health <= 0 {
		if ui, ok := gc.FindEntity("ui").(*UI); ok {
			ui.GameOverScreen = true
		}
	}
	s.Level = int(math.Pow(float64(s.Score+1), 0.1))
}

func (s *SpaceShip) Draw(gc *core.GameContext) {
	color := window.StyleIt(tcell.ColorReset, tcell.ColorRoyalBlue)

	defer s.Gun.Draw(gc)

	window.SetContentWithStyle(
		int(s.triangle.A.GetX()), int(s.triangle.A.GetY()), '^', color)
	window.SetContentWithStyle(
		int(s.triangle.B.GetX()), int(s.triangle.B.GetY()), 'O', color) // right
	window.SetContentWithStyle(
		int(s.triangle.C.GetX()), int(s.triangle.C.GetY()), 'O', color) // left
	// left line
	window.SetContentWithStyle(
		int(s.triangle.A.GetX())-1, int(s.triangle.A.GetY())+1, '/', color)
	// right line
	window.SetContentWithStyle(
		int(s.triangle.A.GetX())+1, int(s.triangle.A.GetY())+1, '\\', color)
	// bottom line
	window.SetContentWithStyle(
		int(s.triangle.A.GetX()), int(s.triangle.A.GetY())+2, '=', color)
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

func (s *SpaceShip) UISpaceshipData(gc *core.GameContext) {
	const startX, startY = 2, 2
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	for i := 0; i < s.Health+2; i++ {
		var ch rune
		switch i {
		case 0:
			ch = '*'
		case 1:
			ch = ' '
		case s.Health:
			ch = tcell.RuneCkBoard
		case s.Health + 1:
			ch = tcell.RuneBoard
		default:
			ch = tcell.RuneBlock
		}
		window.SetContentWithStyle(startX+i, startY, ch, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Score: %d", s.Score)) {
		window.SetContentWithStyle(startX+i, startY+1, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Kills: %d", s.Kills)) {
		window.SetContentWithStyle(startX+i, startY+2, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun CAP: %d/%d", len(s.Gun.Beams), s.Gun.Cap+1)) {
		window.SetContentWithStyle(startX+i, startY+3, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun POW: %d", s.Gun.Power)) {
		window.SetContentWithStyle(startX+i, startY+4, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun SPD: %d", s.Gun.Speed)) {
		window.SetContentWithStyle(startX+i, startY+5, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Level: %d", s.Level)) {
		window.SetContentWithStyle(startX+i, startY+6, r, whiteColor)
	}
}

func (s *SpaceShip) LevelUp(levelit func()) {
	if s.Level > previousLevel {
		levelit()
		if s.Level%2 == 0 {
			s.Gun.Power += 1
			s.Gun.Cap += 1
			s.Gun.Speed += 1
		}
		previousLevel = s.Level
	}
}

func (s *SpaceShip) ScoreKill() {
	s.Kills += 1
	s.Score += s.Kills * s.Kills
}

func (s *SpaceShip) ScoreHit() {
	s.Score += s.Gun.Power
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
