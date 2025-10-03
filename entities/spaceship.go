// Package entities
package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

var (
	previousLevel  int = 0
	nextLevelScore int = 0
	kills          int = 0
	level          int = 0
	score          int = 0
)

type SpaceShip struct {
	Gun
	ObjectBase
	Width, Height     int
	cfg               core.GameConfig
	OnLevelUp         []func(newLevel int)
	SelectedSpaceship *core.SpaceshipDesign
	ListOfSpaceships  []core.SpaceshipDesign
	healthKitsOwned   int
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
	if s.Health >= s.SelectedSpaceship.EntityHealth {
		return false
	}
	s.Health = s.SelectedSpaceship.EntityHealth
	return true
}

func (s *SpaceShip) IncreaseHealth(i int) bool {
	if s.Health < s.SelectedSpaceship.EntityHealth {
		s.Health += i
		s.Health = min(s.Health, s.SelectedSpaceship.EntityHealth)
		return true
	}
	return false
}

func (s *SpaceShip) AddOnLevelUp(fn func(newLevel int)) {
	s.OnLevelUp = append(s.OnLevelUp, fn)
}

// player initialized in the bottom center of the secreen by default

func NewSpaceShip(cfg core.GameConfig, gc *core.GameContext) *SpaceShip {
	w, h := window.GetSize()
	origin := core.PointFloat{
		X: float64(w / 2),
		Y: float64(h - 3),
	}

	// reset values
	previousLevel = 0
	nextLevelScore = 0
	kills = 0
	level = 0
	score = 0
	ModifierHealth = 5
	IncreaseHealthBy = 3

	designs, err := core.LoadListOfAssets[core.SpaceshipDesign]("spaceships.json")
	if err != nil {
		panic(err)
	}

	nextLevelScore = cfg.SpaceShipConfig.NextLevelScore

	return &SpaceShip{
		ObjectBase: ObjectBase{
			OriginPoint: origin,
		},
		ListOfSpaceships: designs,
		cfg:              cfg,
		healthKitsOwned:  1,
	}
}

func (s *SpaceShip) SpaceshipSelection(id int) string {
	s.Gun = Gun{
		Beams: []*Beam{},
		Cap:   s.ListOfSpaceships[id].GunCap,
		Power: s.ListOfSpaceships[id].GunPower,
		Speed: s.ListOfSpaceships[id].GunSpeed,
	}
	s.SelectedSpaceship = &s.ListOfSpaceships[id]
	s.Health = s.ListOfSpaceships[id].EntityHealth
	s.MaxHealth = s.ListOfSpaceships[id].EntityHealth
	s.Width = len(s.ListOfSpaceships[id].Shape[0])
	s.Height = len(s.ListOfSpaceships[id].Shape)
	return s.SelectedSpaceship.Name
}

func (s *SpaceShip) Update(gc *core.GameContext, delta float64) {
	defer s.Gun.Update(gc, delta)
	if s.Health <= 0 && s.SelectedSpaceship != nil {
		if ui, ok := gc.FindEntity("ui").(*UI); ok {
			ui.GameOverScreen = true
		}
	}
	if score >= nextLevelScore {
		level++
		nextLevelScore *= 2
	}
	s.LevelUp()

	s.MovementAndCollision(delta, gc)
}

func (s *SpaceShip) Draw(gc *core.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	color := window.StyleIt(tcell.ColorReset, s.SelectedSpaceship.GetColor())

	defer s.Gun.Draw(gc, s.SelectedSpaceship.GetColor())

	for rowIndex, line := range s.SelectedSpaceship.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(s.OriginPoint.GetX()) + colIndex
				y := int(s.OriginPoint.GetY()) + rowIndex
				window.SetContentWithStyle(x, y, char, color)
			}
		}
	}
	barSize := 5
	DisplayHealth(int(s.OriginPoint.GetX())+(s.Width/2)-(barSize/2)-1, int(s.OriginPoint.GetY())+(s.Height), barSize, s, false, color)
	// -1 because there are the brackets []. So the barSize+[] which is + 2.
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	defer s.Gun.InputEvents(event, gc)

	moveMouse := func(x int, y int) {
		s.OriginPoint.X = float64(x - (s.Width / 2))
		s.OriginPoint.Y = float64(y - (s.Height / 2))
	}

	shootBeam := func() {
		x := int(s.OriginPoint.GetX()) + s.Width/2
		y := int(s.OriginPoint.Y)
		s.initBeam(core.Point{X: x, Y: y}, Up)
	}

	switch ev := event.(type) {
	case *tcell.EventMouse:
		x, y := ev.Position()
		moveMouse(x, y)

		if ev.Buttons() == tcell.Button1 {
			shootBeam()
		}
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			shootBeam()
		}
		if ev.Rune() == 'f' || ev.Rune() == 'F' {
			if s.healthKitsOwned > 0 {
				SetStatus(fmt.Sprintf("[F] Health: Consumed +%d", int(IncreaseHealthBy)))
				if s.IncreaseHealth(int(IncreaseHealthBy)) {
					s.healthKitsOwned--
					return
				}
				SetStatus("[F] Health: Can't use right now")
			} else {
				SetStatus("[F] Health: N/A")
			}
		}

	}
}

func (s *SpaceShip) UISpaceshipData(gc *core.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	const padding, startY = 2, 2
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	for i, r := range []rune(fmt.Sprintf("* Score: %d/%d", score, nextLevelScore)) {
		window.SetContentWithStyle(padding+i, startY, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("* Kills: %d", kills)) {
		window.SetContentWithStyle(padding+i, startY+1, r, whiteColor)
	}

	// display health at the bottome left
	_, h := window.GetSize()
	DisplayHealth(0, h-7, 10, s, true, whiteColor)

	healthStr := []rune(fmt.Sprintf("[HP Kit: %d/%d]", s.healthKitsOwned, MaxHealthKitsToOwn))
	for i, r := range healthStr {
		window.SetContentWithStyle(i, h-8, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[Level: %d", level)) {
		window.SetContentWithStyle(i, h-6, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[CAP:   %d/%d", len(s.Beams), s.Cap)) {
		window.SetContentWithStyle(i, h-5, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[POW:   %d", s.Power)) {
		window.SetContentWithStyle(i, h-4, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[SPD:   %d", s.Speed)) {
		window.SetContentWithStyle(i, h-3, r, whiteColor)
	}
}

func (s *SpaceShip) MovementAndCollision(delta float64, gc *core.GameContext) {
	if a, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		for _, alien := range a.Aliens {
			// check alien shooting the spaceship
			for _, alienBeam := range alien.Beams {
				if s.isHit(&alienBeam.position, alien.Power) {
					alien.RemoveBeam(alienBeam) // removing the beam hitting spaceship
				}
			}
		}
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

	if int(pointBeam.GetX()) >= int(s.OriginPoint.GetX()) &&
		int(pointBeam.GetX()) <= int(s.OriginPoint.GetX())+s.Width &&
		int(pointBeam.GetY()) >= int(s.OriginPoint.GetY()) &&
		int(pointBeam.GetY()) <= int(s.OriginPoint.GetY())+s.Height {

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
	if level > previousLevel {

		// TODO: Refactor
		ModifierHealth += 4
		IncreaseHealthBy += 0.2

		if s.cfg.SpaceShipConfig.MaxLevel <= level {
			return // skip when reaching max level, will not increase any elements of other objects
		}
		for _, fn := range s.OnLevelUp {
			fn(level)
		}
		previousLevel = level
	}
}

func ScoreKill() {
	kills += 1
	score += kills
}

func (s *SpaceShip) ScoreHit() {
	score += s.Power * 2
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}

func (s *SpaceShip) GetHealth() int {
	return s.Health
}

func (s *SpaceShip) GetMaxHealth() int {
	return s.MaxHealth
}
