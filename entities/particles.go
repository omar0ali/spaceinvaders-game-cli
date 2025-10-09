package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

var ExplosiuonSymbols = []rune{'.', '*', 'o', '0', 'O', 'o', '*', ';', '.'}

type Direction = int

const (
	Up Direction = iota
	Down
	Left
	Right
	UpRight
	UpLeft
	DownLeft
	DownRight
)

type Particle struct {
	Speed     float64
	Position  base.PointFloat
	Symbol    []rune
	Direction Direction
	Style     tcell.Style
}

type ParticleProducer struct {
	Particles []*Particle
}
type ParticleSystem struct {
	ParticleProducer *ParticleProducer
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		ParticleProducer: &ParticleProducer{
			Particles: []*Particle{},
		},
	}
}

func (p *ParticleProducer) NewExplosion(scale int, x, y int, width, height int, style tcell.Style) {
	centerX := x + width/2
	centerY := y + height/2

	for dir := Up; dir <= DownRight; dir++ {
		for i := range scale {
			particle := &Particle{
				Speed:     10 + float64(i)*2,
				Position:  base.PointFloat{X: float64(centerX), Y: float64(centerY)},
				Symbol:    ExplosiuonSymbols,
				Direction: dir,
				Style:     style,
			}
			p.Particles = append(p.Particles, particle)
		}
	}
}

func (ps *ParticleSystem) Update(gc *game.GameContext, delta float64) {
	alive := ps.ParticleProducer.Particles[:0] // reset
	for _, p := range ps.ParticleProducer.Particles {
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
			alive = append(alive, p)
		}
	}

	ps.ParticleProducer.Particles = alive
}

// updating all the particles from all the sides

func (ps *ParticleSystem) Draw(gc *game.GameContext) {
	for _, p := range ps.ParticleProducer.Particles {
		currentSymbol := p.Symbol[0] // use the first symbol (it always updates)
		base.SetContentWithStyle(int(p.Position.X), int(p.Position.Y), currentSymbol, p.Style)
	}
}

func (ps *ParticleSystem) InputEvents(event tcell.Event, gc *game.GameContext) {}

func (ps *ParticleSystem) GetType() string {
	return "particles"
}
