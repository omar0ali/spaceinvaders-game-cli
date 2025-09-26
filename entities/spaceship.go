// Package entities
package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type SpaceShip struct {
	Gun            Gun // a spaceship has a gun
	OriginPoint    core.Point
	Health         int
	Kills          int
	Level          int
	Score          int
	NextLevelScore int
	previousLevel  int
	Width, Height  int
	cfg            core.GameConfig
	OnLevelUp      []func(newLevel int)
}

func (s *SpaceShip) AddOnLevelUp(fn func(newLevel int)) {
	s.OnLevelUp = append(s.OnLevelUp, fn)
}

// player initialized in the bottom center of the secreen by default

func NewSpaceShip(cfg core.GameConfig) *SpaceShip {
	w, h := window.GetSize()
	origin := core.Point{
		X: w / 2,
		Y: h - 3,
	}

	return &SpaceShip{
		OriginPoint: origin,

		Gun: Gun{
			Beams: []*Beam{},
			Cap:   cfg.SpaceShipConfig.GunCap,
			Power: cfg.SpaceShipConfig.GunPower,
			Speed: cfg.SpaceShipConfig.GunSpeed,
		},
		Health:         cfg.SpaceShipConfig.Health,
		Width:          3,
		Height:         1,
		cfg:            cfg,
		NextLevelScore: cfg.SpaceShipConfig.NextLevelScore,
	}
}

func (s *SpaceShip) Update(gc *core.GameContext, delta float64) {
	defer s.Gun.Update(gc, delta)
	if s.Health <= 0 {
		if ui, ok := gc.FindEntity("ui").(*UI); ok {
			ui.GameOverScreen = true
		}
	}
	if s.Score >= s.NextLevelScore {
		s.Level++
		s.NextLevelScore *= 2
	}
	s.LevelUp()
}

func (s *SpaceShip) Draw(gc *core.GameContext) {
	color := window.StyleIt(tcell.ColorReset, tcell.ColorRoyalBlue)
	defer s.Gun.Draw(gc)

	spaceshipPattern := []struct {
		dx, dy int
		symbol rune
		color  tcell.Style
	}{
		{s.OriginPoint.X, s.OriginPoint.Y - 1, '^', color},                      // top corner
		{s.OriginPoint.X - s.Width + 1, s.OriginPoint.Y + s.Height, 'O', color}, // left corner
		{s.OriginPoint.X + s.Width - 1, s.OriginPoint.Y + s.Height, 'O', color}, // right corner
		// lines
		{int(s.OriginPoint.GetX()) - 1, int(s.OriginPoint.GetY()), '/', color},     // left line
		{int(s.OriginPoint.GetX()) + 1, int(s.OriginPoint.GetY()), '\\', color},    // right line
		{int(s.OriginPoint.GetX() - 1), int(s.OriginPoint.GetY() + 1), ')', color}, // bottom right
		{int(s.OriginPoint.GetX() + 1), int(s.OriginPoint.GetY() + 1), '(', color}, // bottom right
		{int(s.OriginPoint.GetX()), int(s.OriginPoint.GetY() + 1), '=', color},     // bottom middle
	}
	for _, line := range spaceshipPattern {
		window.SetContentWithStyle(line.dx, line.dy, line.symbol, line.color)
	}
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	defer s.Gun.InputEvents(event, gc)
	moveMouse := func(x int, y int) {
		s.OriginPoint.X = x
		s.OriginPoint.Y = y
	}
	switch ev := event.(type) {
	case *tcell.EventMouse:
		x, y := ev.Position()
		moveMouse(x, y)

		if ev.Buttons() == tcell.Button1 {
			s.Gun.initBeam(s.OriginPoint, Up)
		}
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			s.Gun.initBeam(s.OriginPoint, Up)
		}
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

	for i, r := range []rune(fmt.Sprintf("* Score: %d/%d", s.Score, s.NextLevelScore)) {
		window.SetContentWithStyle(startX+i, startY+1, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Kills: %d", s.Kills)) {
		window.SetContentWithStyle(startX+i, startY+2, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun CAP: %d/%d", len(s.Gun.Beams), s.Gun.Cap)) {
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

func (s *SpaceShip) isHit(pointBeam core.PointInterface, power int) bool {
	grayColor := window.StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)

	// draw flash when hitting
	pattern := []struct {
		dx, dy int
		r      rune
		style  tcell.Style
	}{
		{-1, 0, tcell.RuneBoard, yellowColor},
		{1, 0, tcell.RuneBoard, yellowColor},
		{0, -1, tcell.RuneBoard, grayColor},
		{0, 1, tcell.RuneBoard, grayColor},
		{-1, -1, tcell.RuneCkBoard, grayColor},
		{1, -1, tcell.RuneCkBoard, grayColor},
		{-1, 1, tcell.RuneCkBoard, redColor},
		{1, 1, tcell.RuneBoard, grayColor},
	}

	if int(pointBeam.GetX()) >= s.OriginPoint.X-s.Width &&
		int(pointBeam.GetX()) <= s.OriginPoint.X+s.Width &&
		int(pointBeam.GetY()) >= s.OriginPoint.Y-s.Height &&
		int(pointBeam.GetY()) <= s.OriginPoint.Y+s.Height {

		s.Health -= power // update health of the falling object

		for _, p := range pattern {
			window.SetContentWithStyle(
				int(pointBeam.GetX())+p.dx,
				int(pointBeam.GetY())+p.dy,
				p.r, p.style,
			)
		}

		return true
	}
	return false
}

func (s *SpaceShip) LevelUp() {
	if s.Level > s.previousLevel {
		if s.cfg.SpaceShipConfig.MaxLevel <= s.Level {
			return // skip when reaching max level, will not increase any elements of other objects
		}
		for _, fn := range s.OnLevelUp {
			fn(s.Level)
		}
		s.previousLevel = s.Level
	}
}

func (s *SpaceShip) ScoreKill() {
	s.Kills += 1
	s.Score += 50
}

func (s *SpaceShip) ScoreHit() {
	s.Score += (s.Gun.Power * 2)
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
