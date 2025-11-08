package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
)

type BossProducer struct {
	BossAlien       *base.Enemy
	Level           float64
	LoadedDesigns   *design.LoadedDesigns
	deploymentTimer int
}

func (b *BossProducer) GetType() string {
	return "boss"
}

func NewBossAlienProducer(gc *game.GameContext, designs *design.LoadedDesigns) *BossProducer {
	b := &BossProducer{
		Level:           1.0,
		deploymentTimer: 2,
		LoadedDesigns:   designs,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			b.Level += 0.1
		})
	}

	return b
}

func (b *BossProducer) Update(gc *game.GameContext, delta float64) {
	if b.BossAlien == nil && b.deploymentTimer == minutes {
		SetStatus("Warning: Massive energy spike detected.", gc)
		gc.Sounds.PlaySound("sfx-alarm.mp3", -1)
		b.BossAlien = base.Deploy(b.LoadedDesigns.ListOfBossShips, b.Level)
		b.deploymentTimer += 3
	}

	if b.BossAlien != nil {
		b.BossAlien.Update(gc, delta)
		b.BossAlien.InitBeam(base.Point{
			X: int(b.BossAlien.Position.X) + (b.BossAlien.Width / 2),
			Y: int(b.BossAlien.Position.Y) + (b.BossAlien.Height) + 1,
		}, base.Down, gc.Sounds)

		b.MovementAndCollision(delta, gc)
	}
}

func (b *BossProducer) Draw(gc *game.GameContext) {
	if b.BossAlien == nil {
		return
	}

	color := base.StyleIt(b.BossAlien.GetColor())

	base.DisplayHealthTop(&b.BossAlien.ObjectBase, b.BossAlien.Name, 26, true, color, &b.BossAlien.Gun)

	b.BossAlien.Draw(gc, b.BossAlien.GetColor())

	// draw shape
	for rowIndex, line := range b.BossAlien.Shape {
		for colIndex, char := range line {
			if char != ' ' {
				x := int(b.BossAlien.Position.GetX()) + colIndex
				y := int(b.BossAlien.Position.GetY()) + rowIndex
				base.SetContentWithStyle(x, y, char, color)
			}
		}
	}
}

func (b *BossProducer) InputEvents(event tcell.Event, gc *game.GameContext) {
	// testing code

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'm' { // dev mode
	// 		b.BossAlien = base.Deploy(b.BossShipDesigns, b.Level)
	// 	}
	// }
}

func (b *BossProducer) MovementAndCollision(delta float64, gc *game.GameContext) {
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {

		MoveTo(&b.BossAlien.ObjectBase, &spaceship.ObjectBase, delta, gc)

		for _, beam := range spaceship.GetBeams() {
			if GettingHit(&b.BossAlien.ObjectBase, beam, gc) {
				b.BossAlien.TakeDamage(spaceship.GetPower())
				spaceship.ScoreHit()
				spaceship.RemoveBeam(beam)
			}
		}

		// can collid with a asteroid
		if a, ok := gc.FindEntity("asteroid").(*AsteroidProducer); ok {
			for _, asteroid := range a.Asteroids {
				if Crash(&b.BossAlien.ObjectBase, &asteroid.ObjectBase, gc) {
					b.BossAlien.TakeDamage(1)
					asteroid.TakeDamage(100)
				}
			}
		}

		// can collid with a meteroid
		if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			for _, p := range ps.ParticleProducable {
				switch p.(type) {
				case *particles.MeteroidProducer:
					for _, m := range p.GetParticles() {
						if Crash(&b.BossAlien.ObjectBase, &m.ObjectEntity, gc) {
							b.BossAlien.TakeDamage(1)
							p.RemoveParticle(m)
						}
					}
				}
			}
		}

		if b.BossAlien.IsDead() {
			if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
				ps.AddParticles(
					particles.InitExplosion(15,
						particles.WithDimensions(
							b.BossAlien.Position.X,
							b.BossAlien.Position.Y,
							b.BossAlien.Width,
							b.BossAlien.Height,
						),
					),
				)
			}
			gc.Sounds.PlaySound("8-bit-explosion-low-resonant.mp3", -1)

			spaceship.ScoreKill(b.BossAlien.Health)
			SetStatus("Threat neutralized. Returning to standby.", gc)
			b.BossAlien = nil
		}
	}
}
