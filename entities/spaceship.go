// Package entities
package entities

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/ui"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
)

type Score struct {
	Score          int
	Level          int
	Kills          int
	PreviousLevel  int
	NextLevelScore int
}

type HealthKit struct {
	HealthKitsOwned int
	HealthKitLimit  int
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
	HealthKit         HealthKit
	cfg               game.GameConfig
	OnLevelUp         []func(newLevel int)
	SelectedSpaceship *design.SpaceshipDesign
	LoadedDesigns     *design.LoadedDesigns
	mouseDown         bool
}

func (s *SpaceShip) IncreaseHealthCapacity() bool {
	s.MaxHealth++
	s.Health = s.MaxHealth
	return true
}

func (s *SpaceShip) IncreaseHealth(i int) bool {
	if s.Health <= 0 {
		return false
	}
	if s.Health < s.MaxHealth {
		s.Health += i
		s.Health = min(s.Health, s.MaxHealth)
		return true
	}
	return false
}

func (s *SpaceShip) AddOnLevelUp(fn func(newLevel int)) {
	s.OnLevelUp = append(s.OnLevelUp, fn)
}

// player initialized in the bottom center of the secreen by default

func NewSpaceShip(cfg game.GameConfig, gc *game.GameContext, designs *design.LoadedDesigns) *SpaceShip {
	w, h := base.GetSize()
	origin := base.PointFloat{
		X: float64(w / 2),
		Y: float64(h - 3),
	}

	return &SpaceShip{
		ObjectBase: base.ObjectBase{
			ObjectEntity: base.ObjectEntity{
				Position: origin,
			},
		},
		LoadedDesigns: designs,
		cfg:           cfg,
		HealthKit: HealthKit{
			HealthKitsOwned: 1,
			HealthKitLimit:  5,
		},
		Score: Score{
			NextLevelScore: cfg.SpaceShipConfig.NextLevelScore,
		},
	}
}

func (s *SpaceShip) SpaceshipSelection(id int) string {
	s.Gun = base.NewGun(
		s.LoadedDesigns.ListOfSpaceships[id].GunCap,
		s.LoadedDesigns.ListOfSpaceships[id].GunPower,
		s.LoadedDesigns.ListOfSpaceships[id].GunSpeed,
		s.LoadedDesigns.ListOfSpaceships[id].GunCooldown,
		s.LoadedDesigns.ListOfSpaceships[id].GunReloadCooldown,
	)
	s.SelectedSpaceship = &s.LoadedDesigns.ListOfSpaceships[id]
	s.Health = s.LoadedDesigns.ListOfSpaceships[id].EntityHealth
	s.MaxHealth = s.LoadedDesigns.ListOfSpaceships[id].EntityHealth
	s.Width = len(s.LoadedDesigns.ListOfSpaceships[id].Shape[0])
	s.Height = len(s.LoadedDesigns.ListOfSpaceships[id].Shape)
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
		s.shootBeam(gc)
	}

	s.LevelUp(gc)

	s.MovementAndCollision(delta, gc)
}

func (s *SpaceShip) Draw(gc *game.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	color := base.StyleIt(s.SelectedSpaceship.GetColor())

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
	barSize := 7
	base.DisplayBar(
		s,
		base.WithPosition(
			int(s.Position.GetX())+(s.Width/2)-(barSize/2)-1,
			int(s.Position.GetY())+(s.Height),
		),
		base.WithBarSize(barSize),
		base.WithStyle(base.StyleIt(tcell.ColorGreenYellow)),
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

		if ev.Buttons() == tcell.Button2 {
			if s.GetLoaded() != s.GetCapacity() {
				s.ReloadGun(gc.Sounds)
			}
		}

	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			s.mouseDown = true
		}
		if ev.Rune() == 'E' || ev.Rune() == 'e' {
			if s.HealthKit.HealthKitsOwned > 0 {
				if p, ok := gc.FindEntity("producer").(*ModifierProducer); ok {
					if s.IncreaseHealth(int(p.Level)) {
						SetStatus(fmt.Sprintf("[E] Health: Consumed +%d", int(p.Level)), gc)
						s.HealthKit.HealthKitsOwned--
						return
					}
				}
				SetStatus("[E] Health: Can't use right now", gc)
			} else {
				SetStatus("[E] Health: N/A", gc)
			}
		}
		if ev.Rune() == 'R' || ev.Rune() == 'r' {
			if s.GetLoaded() != s.GetCapacity() {
				s.ReloadGun(gc.Sounds)
			}
		}

	}
}

func (s *SpaceShip) UISpaceshipData(gc *game.GameContext) {
	if s.SelectedSpaceship == nil {
		return
	}

	whiteColor := base.StyleIt(tcell.ColorWhite)
	greenColor := base.StyleIt(tcell.ColorGreenYellow)
	// display health at the bottome left
	_, h := base.GetSize()
	ui.DrawBoxOverlap(
		base.Point{
			X: 0, Y: h - 8,
		}, 23, 6, func(x int, y int) {
			// display health bar of the spaceship at bottom left of the screen
			base.DisplayBar(s, base.WithPosition(x+2, y+1),
				base.WithBarSize(17),
				base.WithStatus(false),
				base.WithStyle(whiteColor),
			)

			for i, r := range fmt.Sprintf("Level:  %d", s.Level) {
				base.SetContentWithStyle(x+i+2, h-6, r, whiteColor)
			}

			str := fmt.Sprintf("CAP:    %d/%d", s.GetLoaded(), s.GetCapacity())

			if s.IsReloading() {
				reloadAnimation := []rune{'·', '•', '●', '○', '●', '•', '·'}
				frame := int(time.Now().UnixNano()/100_000_000) % len(reloadAnimation)
				str += " " + string(reloadAnimation[frame])
			}

			for i, r := range string(str) {
				base.SetContentWithStyle(x+i+2, h-5, r, whiteColor)
			}
			for i, r := range fmt.Sprintf("HP Kit: %d/%d", s.HealthKit.HealthKitsOwned, s.HealthKit.HealthKitLimit) {
				base.SetContentWithStyle(x+i+2, h-4, r, whiteColor)
			}
		}, greenColor)
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
			gc.Sounds.PlaySound("8-bit-explosion.mp3", -1)
		}

		return true
	}
	return false
}

func (s *SpaceShip) ApplyAbility(eff design.AbilityEffect, max int) bool {
	if eff.PowerIncrease != 0 {
		return s.IncreaseGunPower(eff.PowerIncrease)
	}
	if eff.SpeedIncrease != 0 {
		return s.IncreaseGunSpeed(eff.SpeedIncrease, max)
	}
	if eff.CapacityIncrease != 0 {
		return s.IncreaseGunCap(eff.CapacityIncrease, max)
	}
	if eff.CooldownDecrease != 0 {
		return s.DecreaseCooldown(eff.CooldownDecrease)
	}
	if eff.ReloadCooldownDecrease != 0 {
		return s.DecreaseGunReloadCooldown(eff.ReloadCooldownDecrease)
	}
	if eff.HealthCpacity != 0 {
		return s.IncreaseHealthCapacity()
	}
	return false
}

func (s *SpaceShip) LevelUpMenu(gc *game.GameContext) {
	gc.Sounds.PlaySound("8-bit-game-sfx-levelup-menu.mp3", 0)
	if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {
		if u, ok := gc.FindEntity("ui").(*UI); ok {

			SetStatus("Level Up", gc)
			u.LevelUpScreen = true

			var boxes []*ui.Box

			upgrade := func(up func() bool) {
				if up() {
					u.LevelUpScreen = false
					layout.SetLayout(nil)
				}
			}

			displayUpgrade := func(v int, max string) string {
				if v > 0 {
					return fmt.Sprintf("+ %1.f %s", math.Abs(float64(v)), max)
				} else if v < 0 {
					return fmt.Sprintf("- %1.f", math.Abs(float64(v)))
				}
				return ""
			}

			for _, design := range s.LoadedDesigns.ListOfAbilities {
				var displayMax string
				if design.Effect.MaxValue > 0 {
					displayMax = fmt.Sprintf("(Max: %d)", design.Effect.MaxValue)
				}

				increaseSpeed := design.Effect.SpeedIncrease
				increaseCap := design.Effect.CapacityIncrease
				decreaseCD := design.Effect.CooldownDecrease
				decreaseRDCD := design.Effect.ReloadCooldownDecrease
				increasePower := design.Effect.PowerIncrease
				increaseHealthCap := design.Effect.HealthCpacity

				boxes = append(
					boxes,
					ui.NewUIBox(
						design.Shape,
						[]string{
							"(*) " + design.Name,
							"Details: " + design.Description,
							fmt.Sprintf("Gun Power: (%d) %s", s.GetPower(), displayUpgrade(increasePower, displayMax)),
							fmt.Sprintf("Gun Capacity: (%d) %s", s.GetCapacity(), displayUpgrade(increaseCap, displayMax)),
							fmt.Sprintf("Gun Speed: (%d) %s", s.GetSpeed(), displayUpgrade(increaseSpeed, displayMax)),
							fmt.Sprintf("Gun Cooldown: (%d) %s", s.GetCooldown(), displayUpgrade(decreaseCD, displayMax)),
							fmt.Sprintf("Gun Reload Cooldown: (%d) %s", s.GetReloadCooldown(), displayUpgrade(decreaseRDCD, displayMax)),
							fmt.Sprintf("Health Capacity: (%d) %s", s.MaxHealth, displayUpgrade(increaseHealthCap, displayMax)),
						},

						func() {
							upgrade(func() bool {
								if design.Status != "" {
									SetStatus(design.Status, gc)
								}
								return s.ApplyAbility(design.Effect, design.Effect.MaxValue)
							})
						},
					),
				)
			}

			// shuffle the list
			rand.Shuffle(len(boxes), func(i, j int) {
				boxes[i], boxes[j] = boxes[j], boxes[i]
			})

			// pick the first 3 boxes
			pickedBoxes := boxes[:3]

			layout.SetLayout(
				ui.InitLayout(21, 10, pickedBoxes...),
			)
		}
	}
}

func (s *SpaceShip) LevelUp(gc *game.GameContext) {
	if s.Level > s.PreviousLevel {
		if s.cfg.SpaceShipConfig.MaxLevel <= s.Level {
			return // skip when reaching max level, will not increase any elements of other objects
		}
		for _, fn := range s.OnLevelUp {
			fn(s.Level)
		}

		// pop up level up
		s.LevelUpMenu(gc)

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

func (s *SpaceShip) shootBeam(gc *game.GameContext) {
	x := int(s.Position.GetX()) + s.Width/2
	y := int(s.Position.Y)
	s.InitBeam(base.Point{X: x, Y: y}, base.Up, gc.Sounds)
}

func (s *SpaceShip) GetMax() int {
	return s.MaxHealth
}
