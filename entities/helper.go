package entities

import (
	"math"

	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type PointableInt interface {
	GetPosition() *base.Point
}

type PointableFloat interface {
	GetPosition() *base.PointFloat
}

type Sizeable interface {
	GetWidth() int
	GetHeight() int
}

type Movable interface {
	Sizeable
	PointableFloat
	AppendPositionY(float64)
	GetSpeed() float64
}

func Move(m Movable, delta float64) {
	distance := m.GetSpeed() * delta
	m.AppendPositionY(distance)
}

func MoveTo(from, to Movable, delta float64, gc *game.GameContext) {
	distance := from.GetSpeed() * delta

	const toleranceX = 2
	const toleranceY = 5

	bossCenterX := from.GetPosition().X + float64(from.GetWidth())/2
	shipCenterX := to.GetPosition().X + float64(to.GetWidth())/2 + 2

	if math.Abs(bossCenterX-shipCenterX) > toleranceX {
		if bossCenterX > shipCenterX {
			from.GetPosition().AppendX(-distance)
		} else {
			from.GetPosition().AppendX(distance)
		}
	}

	targetY := to.GetPosition().Y - float64(to.GetHeight()) - 10
	if targetY < -5 {
		targetY = -5
	}

	if math.Abs(from.GetPosition().Y-targetY) > toleranceY {
		if from.GetPosition().Y > targetY {
			from.GetPosition().AppendY(-distance)
		} else {
			from.GetPosition().AppendY(distance)
		}
	}
}

func GettingHit(m Movable, beam PointableInt, gc *game.GameContext) bool {
	px := int(math.Round(beam.GetPosition().GetX()))
	py := int(math.Round(beam.GetPosition().GetY()))
	ox := int(math.Round(m.GetPosition().GetX()))
	oy := int(math.Round(m.GetPosition().GetY()))

	if px >= ox && px < ox+m.GetWidth() &&
		py >= oy && py < oy+m.GetHeight() {

		if p, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			p.AddParticles(
				particles.InitExplosion(3,
					particles.WithDimensions(
						float64(beam.GetPosition().X),
						float64(beam.GetPosition().Y),
						0,
						0,
					),
					particles.WithSymbols([]rune("Oo;.")),
				),
			)
			gc.Sounds.PlaySound("8-bit-explosion.mp3", -1)
		}
		return true
	}
	return false
}

func Crash(c1, c2 Movable, gc *game.GameContext) bool {
	x1 := int(math.Round(c1.GetPosition().GetX()))
	y1 := int(math.Round(c1.GetPosition().GetY()))
	w1 := c1.GetWidth()
	h1 := c1.GetHeight()

	x2 := int(math.Round(c2.GetPosition().GetX()))
	y2 := int(math.Round(c2.GetPosition().GetY()))
	w2 := c2.GetWidth()
	h2 := c2.GetHeight()

	if x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2 {

		if p, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			p.AddParticles(
				particles.InitExplosion(3,
					particles.WithDimensions(
						c1.GetPosition().X,
						c1.GetPosition().Y,
						c1.GetWidth(),
						c1.GetHeight(),
					),
					particles.WithSymbols([]rune(".oO0*;.")),
				),
			)
			p.AddParticles(
				particles.InitExplosion(3,
					particles.WithDimensions(
						c2.GetPosition().X,
						c2.GetPosition().Y,
						c2.GetWidth(),
						c2.GetHeight(),
					),
					particles.WithSymbols([]rune(".oO0*;.")),
				),
			)

		}

		return true
	}

	return false
}
