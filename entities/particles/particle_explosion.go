package particles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type ExplosionProducer struct {
	Particles []*Particle
}

func (e *ExplosionProducer) GetParticles() int {
	return len(e.Particles)
}

func InitExplosion(scale int, opts ...ParticleOption) *ExplosionProducer {
	var listOfParticles []*Particle

	po := &Particle{
		Position: base.PointFloat{X: 0, Y: 0},
		Style:    base.StyleIt(tcell.ColorReset, tcell.ColorYellow),
		Symbol:   []rune{'0', 'O', 'o', '*', ';', '.'},
		Speed:    10,
	}

	for _, o := range opts {
		o(po)
	}

	for dir := Up; dir <= DownRight; dir++ { // for each direction
		for i := range scale {
			particle := &Particle{
				Speed:     po.Speed + float64(i)*2,
				Position:  po.Position,
				Symbol:    po.Symbol,
				Direction: dir,
				Style:     po.Style,
			}
			listOfParticles = append(listOfParticles, particle)
		}
	}
	return &ExplosionProducer{
		Particles: listOfParticles,
	}
}

func (e *ExplosionProducer) Update(gc *game.GameContext, delta float64) {
	activeParticles := e.Particles[:0] // reset
	for _, p := range e.Particles {
		switch p.Direction {
		case Up:
			p.Position.Y -= (float64(p.Speed) * delta)
		case Down:
			p.Position.Y += (float64(p.Speed) * delta)
		case Left:
			p.Position.X -= (float64(p.Speed) * delta)
		case Right:
			p.Position.X += (float64(p.Speed) * delta)
		case UpRight:
			p.Position.Y -= (float64(p.Speed) * delta)
			p.Position.X += (float64(p.Speed) * delta)
		case UpLeft:
			p.Position.Y -= (float64(p.Speed) * delta)
			p.Position.X -= (float64(p.Speed) * delta)
		case DownLeft:
			p.Position.Y += (float64(p.Speed) * delta)
			p.Position.X -= (float64(p.Speed) * delta)
		case DownRight:
			p.Position.Y += (float64(p.Speed) * delta)
			p.Position.X += (float64(p.Speed) * delta)
		}

		// shrink
		if len(p.Symbol) > 1 {
			p.Symbol = p.Symbol[1:]
			activeParticles = append(activeParticles, p)
		}
	}

	e.Particles = activeParticles
}

// updating all the particles from all the sides

func (e *ExplosionProducer) Draw(gc *game.GameContext) {
	for _, p := range e.Particles {
		currentSymbol := p.Symbol[0] // use the first symbol (it always updates)
		base.SetContentWithStyle(int(p.Position.X), int(p.Position.Y), currentSymbol, p.Style)
	}
}
