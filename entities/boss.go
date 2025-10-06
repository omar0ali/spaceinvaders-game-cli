package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

var everyThreeMinutes = 2

type BossProducer struct {
	BossAlien *base.Enemy
	Level     float64
}

func (b *BossProducer) GetType() string {
	return "boss"
}

func NewBossAlienProducer(gc *game.GameContext) *BossProducer {
	b := &BossProducer{
		Level: 1.0,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			b.Level += 0.1
		})
	}

	return b
}

func (b *BossProducer) Update(gc *game.GameContext, delta float64) {
	if b.BossAlien == nil && everyThreeMinutes == minutes {
		SetStatus("Warning: Massive energy spike detected.")
		b.BossAlien = base.Deploy("bossships.json", b.Level)
		everyThreeMinutes += 3
	}

	if b.BossAlien != nil {
		b.BossAlien.Update(gc, delta)
		b.BossAlien.InitBeam(base.Point{
			X: int(b.BossAlien.Position.X) + (b.BossAlien.Width / 2),
			Y: int(b.BossAlien.Position.Y) + (b.BossAlien.Height) + 1,
		}, base.Down)

		b.MovementAndCollision(delta, gc)
	}
}

func (b *BossProducer) Draw(gc *game.GameContext) {
	if b.BossAlien == nil {
		return
	}
	colorHealth := base.StyleIt(tcell.ColorReset, tcell.ColorIndianRed)
	color := base.StyleIt(tcell.ColorReset, b.BossAlien.GetColor())

	b.BossAlien.DisplayHealth(6, true, colorHealth, &b.BossAlien.Gun)

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
	// 		b.BossAlien = base.Deploy("bossships.json", int(b.Level))
	// 	}
	// }
}

func (b *BossProducer) MovementAndCollision(delta float64, gc *game.GameContext) {
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {

		base.MoveTo(&b.BossAlien.ObjectBase, &spaceship.ObjectBase, delta, gc)

		for _, beam := range spaceship.GetBeams() {
			if base.GettingHit(&b.BossAlien.ObjectBase, beam) {
				b.BossAlien.TakeDamage(spaceship.GetPower())
				spaceship.ScoreHit()
				spaceship.RemoveBeam(beam)
			}
		}
	}

	if b.BossAlien.IsDead() {
		ScoreKill()
		SetStatus("Threat neutralized. Returning to standby.")
		b.BossAlien = nil
	}
}
