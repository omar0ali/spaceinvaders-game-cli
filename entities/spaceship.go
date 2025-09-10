// Package entities
package entities

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type SpaceshipOpts struct {
	SpaceShipHealth int
	GunPower        int
	GunCapacity     int
	GunSpeed        int
}

type SpaceShip struct {
	origin        core.Point
	Gun           Gun // a spaceship has a gun
	Health        int
	Score         int
	Kills         int
	Level         int
	previousLevel int
	Width, Height int
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

		Gun: Gun{
			Beams: []*Beam{},
			Cap:   opts.GunCapacity,
			Power: opts.GunPower,
			Speed: opts.GunSpeed,
		},
		Health: opts.SpaceShipHealth,
		Width:  3,
		Height: 1,
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

	spaceshipPattern := []struct {
		dx, dy int
		symbol rune
		color  tcell.Style
	}{
		{s.origin.X, s.origin.Y - 1, '^', color},                      // top corner
		{s.origin.X - s.Width + 1, s.origin.Y + s.Height, 'O', color}, // left corner
		{s.origin.X + s.Width - 1, s.origin.Y + s.Height, 'O', color}, // right corner
		// lines
		{int(s.origin.GetX()) - 1, int(s.origin.GetY()), '/', color},     // left line
		{int(s.origin.GetX()) + 1, int(s.origin.GetY()), '\\', color},    // right line
		{int(s.origin.GetX() - 1), int(s.origin.GetY() + 1), ')', color}, // bottom right
		{int(s.origin.GetX() + 1), int(s.origin.GetY() + 1), '(', color}, // bottom right
		{int(s.origin.GetX()), int(s.origin.GetY() + 1), '=', color},     // bottom middle
	}
	for _, line := range spaceshipPattern {
		window.SetContentWithStyle(line.dx, line.dy, line.symbol, line.color)
	}
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	defer s.Gun.InputEvents(event, gc)
	moveMouse := func(x int) {
		s.origin.X = x
		s.origin.SetX(float64(x))
		s.origin.SetX(float64(x + 2))
		s.origin.SetX(float64(x - 2))
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
	if s.Level > s.previousLevel {
		levelit()
		if s.Level%2 == 0 {
			s.Gun.Power += 1
			s.Gun.Cap += 1
			s.Gun.Speed += 1
		}
		s.previousLevel = s.Level
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
