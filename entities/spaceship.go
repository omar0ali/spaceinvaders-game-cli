// Package entities
package entities

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
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
	base.Gun
	base.ObjectBase
	cfg               game.GameConfig
	OnLevelUp         []func(newLevel int)
	SelectedSpaceship *game.SpaceshipDesign
	ListOfSpaceships  []game.SpaceshipDesign
	healthKitsOwned   int
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

func NewSpaceShip(cfg game.GameConfig, gc *game.GameContext) *SpaceShip {
	w, h := window.GetSize()
	origin := game.PointFloat{
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

	designs, err := game.LoadListOfAssets[game.SpaceshipDesign]("spaceships.json")
	if err != nil {
		panic(err)
	}

	nextLevelScore = cfg.SpaceShipConfig.NextLevelScore

	return &SpaceShip{
		ObjectBase: base.ObjectBase{
			Position: origin,
		},
		ListOfSpaceships: designs,
		cfg:              cfg,
		healthKitsOwned:  1,
	}
}

func (s *SpaceShip) SpaceshipSelection(id int) string {
	s.Gun = base.NewGun(
		s.ListOfSpaceships[id].GunCap,
		s.ListOfSpaceships[id].GunPower,
		s.ListOfSpaceships[id].GunSpeed,
		s.ListOfSpaceships[id].GunCooldown,
		s.ListOfSpaceships[id].GunReloadCooldown,
	)
	s.SelectedSpaceship = &s.ListOfSpaceships[id]
	s.Health = s.ListOfSpaceships[id].EntityHealth
	s.MaxHealth = s.ListOfSpaceships[id].EntityHealth
	s.Width = len(s.ListOfSpaceships[id].Shape[0])
	s.Height = len(s.ListOfSpaceships[id].Shape)
	return s.SelectedSpaceship.Name
}

func (s *SpaceShip) Update(gc *game.GameContext, delta float64) {
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

func (s *SpaceShip) Draw(gc *game.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	color := window.StyleIt(tcell.ColorReset, s.SelectedSpaceship.GetColor())

	defer s.Gun.Draw(gc, s.SelectedSpaceship.GetColor())

	for rowIndex, line := range s.SelectedSpaceship.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(s.Position.GetX()) + colIndex
				y := int(s.Position.GetY()) + rowIndex
				window.SetContentWithStyle(x, y, char, color)
			}
		}
	}
	barSize := 5
	base.DisplayBar(int(s.Position.GetX())+(s.Width/2)-(barSize/2)-1, int(s.Position.GetY())+(s.Height), barSize, s, false, color, &s.Gun)
	// -1 because there are the brackets []. So the barSize+[] which is + 2.
}

func (s *SpaceShip) InputEvents(event tcell.Event, gc *game.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	moveMouse := func(x int, y int) {
		s.Position.X = float64(x - (s.Width / 2))
		s.Position.Y = float64(y - (s.Height / 2))
	}

	shootBeam := func() {
		x := int(s.Position.GetX()) + s.Width/2
		y := int(s.Position.Y)
		s.InitBeam(game.Point{X: x, Y: y}, base.Up)
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
		if ev.Rune() == 'E' || ev.Rune() == 'e' {
			if s.healthKitsOwned > 0 {
				SetStatus(fmt.Sprintf("[E] Health: Consumed +%d", int(IncreaseHealthBy)))
				if s.IncreaseHealth(int(IncreaseHealthBy)) {
					s.healthKitsOwned--
					return
				}
				SetStatus("[E] Health: Can't use right now")
			} else {
				SetStatus("[E] Health: N/A")
			}
		}
		if ev.Rune() == 'R' || ev.Rune() == 'r' {
			if s.GetLoaded() != s.GetCapacity() {
				s.ReloadGun()
			}
		}

	}
}

func (s *SpaceShip) UISpaceshipData(gc *game.GameContext) {
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

	healthStr := []rune(fmt.Sprintf("[HP Kit: %d/%d]", s.healthKitsOwned, MaxHealthKitsToOwn))
	for i, r := range healthStr {
		window.SetContentWithStyle(i, h-10, r, whiteColor)
	}

	base.DisplayBar(0, h-9, 10, s, true, whiteColor, &s.Gun)

	for i, r := range []rune(fmt.Sprintf("[Level:     %d", level)) {
		window.SetContentWithStyle(i, h-8, r, whiteColor)
	}

	reloadAnimation := []rune("•○")
	str := fmt.Sprintf("[CAP:    %d/%d", s.GetLoaded(), s.GetCapacity())

	if s.IsReloading() {
		frame := int(time.Now().UnixNano()/300_000_000) % len(reloadAnimation)
		str += " " + string(reloadAnimation[frame]) + " RELOADING"
	}
	for i, r := range []rune(str) {
		window.SetContentWithStyle(i, h-7, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[POW:    %d", s.GetPower())) {
		window.SetContentWithStyle(i, h-6, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[SPD:    %d", int(s.Gun.GetSpeed()))) {
		window.SetContentWithStyle(i, h-5, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[CD:     %d ms", int(s.GetCooldown()))) {
		window.SetContentWithStyle(i, h-4, r, whiteColor)
	}
	for i, r := range []rune(fmt.Sprintf("[RLD:    %d ms", int(s.GetReloadCooldown()))) {
		window.SetContentWithStyle(i, h-3, r, whiteColor)
	}
}

func (s *SpaceShip) MovementAndCollision(delta float64, gc *game.GameContext) {
	if a, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		for _, alien := range a.Aliens {
			// check alien shooting the spaceship
			for _, alienBeam := range alien.GetBeams() {
				if s.isHit(alienBeam.GetPosition(), alien.GetPower()) {
					s.TakeDamage(alien.GetPower())
					alien.RemoveBeam(alienBeam)
				}
			}
			if base.Crash(&s.ObjectBase, &alien.ObjectBase) {
				s.TakeDamage(1)
				alien.TakeDamage(5)
			}
		}
	}
	if b, ok := gc.FindEntity("boss").(*BossProducer); ok {
		if b.BossAlien != nil {
			for _, bossBeam := range b.BossAlien.GetBeams() {
				if s.isHit(bossBeam.GetPosition(), b.BossAlien.GetPower()) {
					s.TakeDamage(b.BossAlien.GetPower())
					b.BossAlien.RemoveBeam(bossBeam)
				}
			}
			if base.Crash(&s.ObjectBase, &b.BossAlien.ObjectBase) {
				s.TakeDamage(1)
				b.BossAlien.TakeDamage(5)
			}

		}
	}
}

func (s *SpaceShip) isHit(pointBeam game.PointInterface, power int) bool {
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

	if int(pointBeam.GetX()) >= int(s.Position.GetX()) &&
		int(pointBeam.GetX()) <= int(s.Position.GetX())+s.Width &&
		int(pointBeam.GetY()) >= int(s.Position.GetY()) &&
		int(pointBeam.GetY()) <= int(s.Position.GetY())+s.Height {

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
		ModifierHealth += 5
		IncreaseHealthBy += 0.3

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
	score += s.GetPower() * 2
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
