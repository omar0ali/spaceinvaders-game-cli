// Package entities
package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type SpaceShip struct {
	Gun
	OriginPoint     core.Point
	health          int
	Kills           int
	Level           int
	Score           int
	NextLevelScore  int
	previousLevel   int
	Width, Height   int
	cfg             core.GameConfig
	OnLevelUp       []func(newLevel int)
	SpaceshipDesign core.SpaceshipDesign
}

func (s *SpaceShip) IncreaseGunPower(i int) bool {
	s.Power += i
	return true
}

func (s *SpaceShip) IncreaseGunSpeed(i int) bool {
	if s.Speed < s.cfg.SpaceShipConfig.GunMaxSpeed {
		s.Speed += i
		s.Speed = min(s.Speed, s.cfg.SpaceShipConfig.GunMaxSpeed)
		return true
	}
	return false
}

func (s *SpaceShip) IncreaseGunCap(i int) bool {
	if s.Cap < s.cfg.SpaceShipConfig.GunMaxCap {
		s.Cap += i
		s.Cap = min(s.Cap, s.cfg.SpaceShipConfig.GunMaxCap)
		return true
	}
	return false
}

func (s *SpaceShip) RestoreFullHealth() bool {
	if s.health >= s.SpaceshipDesign.EntityHealth {
		return false
	}
	s.health = s.SpaceshipDesign.EntityHealth
	return true
}

func (s *SpaceShip) IncreaseHealth(i int) bool {
	if s.health < s.SpaceshipDesign.EntityHealth {
		s.health += i
		s.health = min(s.health, s.SpaceshipDesign.EntityHealth)
		return true
	}
	return false
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
	designs, err := core.LoadListOfAssets[core.SpaceshipDesign]("assets/spaceship.json")
	design := designs[rand.Intn(len(designs))]
	if err != nil {
		panic(err)
	}
	width := len(design.Shape[0])
	height := len(design.Shape)

	return &SpaceShip{
		OriginPoint:     origin,
		SpaceshipDesign: design,
		Gun: Gun{
			Beams: []*Beam{},
			Cap:   design.GunCap,
			Power: design.GunPower,
			Speed: design.GunSpeed,
		},
		health:         design.EntityHealth,
		Width:          width,
		Height:         height,
		cfg:            cfg,
		NextLevelScore: cfg.SpaceShipConfig.NextLevelScore,
	}
}

func (s *SpaceShip) Update(gc *core.GameContext, delta float64) {
	defer s.Gun.Update(gc, delta)
	if s.health <= 0 {
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

	for rowIndex, line := range s.SpaceshipDesign.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(s.OriginPoint.GetX()) + colIndex
				y := int(s.OriginPoint.GetY()) + rowIndex
				window.SetContentWithStyle(x, y, char, color)
			}
		}
	}
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	defer s.Gun.InputEvents(event, gc)
	moveMouse := func(x int, y int) {
		s.OriginPoint.X = x - (s.Width / 2)
		s.OriginPoint.Y = y - (s.Height / 2)
	}
	switch ev := event.(type) {
	case *tcell.EventMouse:
		x, y := ev.Position()
		moveMouse(x, y)

		if ev.Buttons() == tcell.Button1 {
			x := s.OriginPoint.X + s.Width/2
			y := s.OriginPoint.Y
			s.initBeam(core.Point{X: x, Y: y}, Up)
		}
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			s.initBeam(s.OriginPoint, Up)
		}
	}
}

func (s *SpaceShip) UISpaceshipData(gc *core.GameContext) {
	const padding, startY = 2, 2
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	for i, r := range []rune(fmt.Sprintf("* Score: %d/%d", s.Score, s.NextLevelScore)) {
		window.SetContentWithStyle(padding+i, startY, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Kills: %d", s.Kills)) {
		window.SetContentWithStyle(padding+i, startY+1, r, whiteColor)
	}

	endPositionOfHealth := 4
	for i := padding; i < s.health+(padding*2); i++ {
		var ch rune
		switch i {
		case padding:
			ch = '*'
		case padding + 1:
			ch = ' '
		default:
			endPositionOfHealth = i + (padding * 2) // more padding
			ch = tcell.RuneBlock
		}
		window.SetContentWithStyle(i, startY+3, ch, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("%d/%d", s.health, s.SpaceshipDesign.EntityHealth)) {
		window.SetContentWithStyle(endPositionOfHealth+i, startY+3, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun CAP: %d/%d", len(s.Beams), s.Cap)) {
		window.SetContentWithStyle(padding+i, startY+4, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun POW: %d", s.Power)) {
		window.SetContentWithStyle(padding+i, startY+5, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Gun SPD: %d", s.Speed)) {
		window.SetContentWithStyle(padding+i, startY+6, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Level: %d", s.Level)) {
		window.SetContentWithStyle(padding+i, startY+7, r, whiteColor)
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

	if int(pointBeam.GetX()) >= s.OriginPoint.X &&
		int(pointBeam.GetX()) <= s.OriginPoint.X+s.Width &&
		int(pointBeam.GetY()) >= s.OriginPoint.Y &&
		int(pointBeam.GetY()) <= s.OriginPoint.Y+s.Height {

		s.health -= power // update health of the falling object

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
	s.Score += s.Kills
}

func (s *SpaceShip) ScoreHit() {
	s.Score += s.Power * 2
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
