// Package entities
package entities

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type Score struct {
	Score          int
	Level          int
	Kills          int
	PreviousLevel  int
	NextLevelScore int
}

func (s *Score) GetCurrent() int {
	return s.Score
}

func (s *Score) GetMax() int {
	return s.NextLevelScore
}

type SpaceShip struct {
	base.Gun
	base.ObjectBase
	Score
	cfg               game.GameConfig
	OnLevelUp         []func(newLevel int)
	SelectedSpaceship *game.SpaceshipDesign
	ListOfSpaceships  []game.SpaceshipDesign
	healthKitsOwned   int
	mouseDown         bool
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
	w, h := base.GetSize()
	origin := base.PointFloat{
		X: float64(w / 2),
		Y: float64(h - 3),
	}

	designs, err := game.LoadListOfAssets[game.SpaceshipDesign]("spaceships.json")
	if err != nil {
		panic(err)
	}

	return &SpaceShip{
		ObjectBase: base.ObjectBase{
			ObjectEntity: base.ObjectEntity{
				Position: origin,
			},
		},
		ListOfSpaceships: designs,
		cfg:              cfg,
		healthKitsOwned:  1,
		Score: Score{
			NextLevelScore: cfg.SpaceShipConfig.NextLevelScore,
		},
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
	if s.Score.Score >= s.NextLevelScore {
		s.Level++
		s.NextLevelScore += s.cfg.SpaceShipConfig.NextLevelScore
	}

	if s.mouseDown {
		s.shootBeam()
	}

	s.LevelUp()

	s.MovementAndCollision(delta, gc)
}

func (s *SpaceShip) Draw(gc *game.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	color := base.StyleIt(tcell.ColorReset, s.SelectedSpaceship.GetColor())

	defer s.Gun.Draw(gc, s.SelectedSpaceship.GetColor())

	for rowIndex, line := range s.SelectedSpaceship.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(s.Position.GetX()) + colIndex
				y := int(s.Position.GetY()) + rowIndex
				base.SetContentWithStyle(x, y, char, color)
			}
		}
	}

	// display health bar at the bottom of the spaceship
	barSize := 5
	base.DisplayBar(
		s,
		base.WithPosition(
			int(s.Position.GetX())+(s.Width/2)-(barSize/2)-1,
			int(s.Position.GetY())+(s.Height),
		),
		base.WithBarSize(barSize),
		base.WithStatus(true),
		base.WithStyle(color),
		base.WithGun(&s.Gun),
	)

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

	switch ev := event.(type) {
	case *tcell.EventMouse:
		x, y := ev.Position()
		moveMouse(x, y)

		// buttons() contains (0000 0001, 0000 0100, 0000 0101)
		// & symbol keeps only bits that are on in both
		// Button1 = 0000 0001 and if that pressed and if any bit remains it will equal to true
		if ev.Buttons()&tcell.Button1 != 0 {
			s.mouseDown = true
		} else {
			s.mouseDown = false
		}

	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			s.mouseDown = true
		}
		if ev.Rune() == 'E' || ev.Rune() == 'e' {
			if s.healthKitsOwned > 0 {
				if p, ok := gc.FindEntity("producer").(*ModifierProducer); ok {
					SetStatus(fmt.Sprintf("[E] Health: Consumed +%d", p.ConsumableHealth))
					if s.IncreaseHealth(p.ConsumableHealth) {
						s.healthKitsOwned--
						return
					}
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
	whiteColor := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	for i, r := range []rune("* Score: ") {
		base.SetContentWithStyle(padding+i, startY, r, whiteColor)
	}

	// display score
	base.DisplayBar(
		&s.Score,
		base.WithPosition(padding+9, startY),
		base.WithBarSize(10),
		base.WithStatus(false),
		base.WithStyle(whiteColor),
	)

	for i, r := range []rune(fmt.Sprintf("* Kills: %d", s.Kills)) {
		base.SetContentWithStyle(padding+i, startY+1, r, whiteColor)
	}

	// display health at the bottome left
	_, h := base.GetSize()

	healthStr := []rune(fmt.Sprintf("[HP Kit: %d/%d]", s.healthKitsOwned, MaxHealthKitsToOwn))
	for i, r := range healthStr {
		base.SetContentWithStyle(i, h-10, r, whiteColor)
	}

	// display health bar of the spaceship at bottom left of the screen
	base.DisplayBar(
		s,
		base.WithPosition(0, h-9),
		base.WithBarSize(10),
		base.WithStatus(false),
		base.WithStyle(whiteColor),
		base.WithGun(&s.Gun),
	)

	for i, r := range []rune(fmt.Sprintf("[Level:  %d", s.Level)) {
		base.SetContentWithStyle(i, h-8, r, whiteColor)
	}

	reloadAnimation := []rune("•○")
	str := fmt.Sprintf("[CAP:    %d/%d", s.GetLoaded(), s.GetCapacity())

	if s.IsReloading() {
		frame := int(time.Now().UnixNano()/300_000_000) % len(reloadAnimation)
		str += " " + string(reloadAnimation[frame]) + " RELOADING"
	}
	for i, r := range []rune(str) {
		base.SetContentWithStyle(i, h-7, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[POW:    %d", s.GetPower())) {
		base.SetContentWithStyle(i, h-6, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[SPD:    %d", int(s.GetSpeed()))) {
		base.SetContentWithStyle(i, h-5, r, whiteColor)
	}

	for i, r := range []rune(fmt.Sprintf("[CD:     %d ms", int(s.GetCooldown()))) {
		base.SetContentWithStyle(i, h-4, r, whiteColor)
	}
	for i, r := range []rune(fmt.Sprintf("[RLD:    %d ms", int(s.GetReloadCooldown()))) {
		base.SetContentWithStyle(i, h-3, r, whiteColor)
	}
}

func (s *SpaceShip) MovementAndCollision(delta float64, gc *game.GameContext) {
	if a, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		for _, alien := range a.Aliens {
			// check alien shooting the spaceship
			for _, alienBeam := range alien.GetBeams() {
				if s.isHit(alienBeam.GetPosition(), gc) {
					s.TakeDamage(alien.GetPower())
					alien.RemoveBeam(alienBeam)
				}
			}
			if Crash(&s.ObjectBase, &alien.ObjectBase, gc) {
				s.TakeDamage(1)
				alien.TakeDamage(5)
			}
		}
	}

	// can collid with a meteroid
	if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
		for _, p := range ps.ParticleProducable {
			switch p.(type) {
			case *particles.MeteroidProducer:
				for _, m := range p.GetParticles() {
					if Crash(&s.ObjectBase, &m.ObjectEntity, gc) {
						s.TakeDamage(2)
						p.RemoveParticle(m)
					}
				}
			}
		}
	}

	if b, ok := gc.FindEntity("boss").(*BossProducer); ok {
		if b.BossAlien != nil {
			for _, bossBeam := range b.BossAlien.GetBeams() {
				if s.isHit(bossBeam.GetPosition(), gc) {
					s.TakeDamage(b.BossAlien.GetPower())
					b.BossAlien.RemoveBeam(bossBeam)
				}
			}

			// can collid with a asteroid

			if Crash(&s.ObjectBase, &b.BossAlien.ObjectBase, gc) {
				s.TakeDamage(1)
				b.BossAlien.TakeDamage(5)
			}

		}
	}
	if a, ok := gc.FindEntity("asteroid").(*AsteroidProducer); ok {
		for _, asteroid := range a.Asteroids {
			if Crash(&s.ObjectBase, &asteroid.ObjectBase, gc) {
				s.TakeDamage(2)
				asteroid.TakeDamage(4)
			}
		}
	}
}

func (s *SpaceShip) isHit(pointBeam base.PointInterface, gc *game.GameContext) bool {
	if int(pointBeam.GetX()) >= int(s.Position.GetX()) &&
		int(pointBeam.GetX()) <= int(s.Position.GetX())+s.Width &&
		int(pointBeam.GetY()) >= int(s.Position.GetY()) &&
		int(pointBeam.GetY()) <= int(s.Position.GetY())+s.Height {

		if p, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			p.AddParticles(
				particles.InitExplosion(3,
					particles.WithDimensions(
						pointBeam.GetX(),
						pointBeam.GetY(),
						0,
						0,
					),
					particles.WithSymbols([]rune("0%*;.")),
				),
			)
		}

		return true
	}
	return false
}

func (s *SpaceShip) LevelUp() {
	if s.Level > s.PreviousLevel {
		if s.cfg.SpaceShipConfig.MaxLevel <= s.Level {
			return // skip when reaching max level, will not increase any elements of other objects
		}
		for _, fn := range s.OnLevelUp {
			fn(s.Level)
		}
		s.PreviousLevel = s.Level
		s.Score.Score = 0
	}
}

func (s *SpaceShip) ScoreKill() {
	s.Kills += 1
	s.Score.Score += s.Kills
}

func (s *SpaceShip) ScoreHit() {
	s.Score.Score += s.GetPower()
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}

func (s *SpaceShip) GetCurrent() int {
	return s.Health
}

func (s *SpaceShip) shootBeam() {
	x := int(s.Position.GetX()) + s.Width/2
	y := int(s.Position.Y)
	s.InitBeam(base.Point{X: x, Y: y}, base.Up)
}

func (s *SpaceShip) GetMax() int {
	return s.MaxHealth
}
