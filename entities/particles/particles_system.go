// Package particles
package particles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type ParticleProducable interface {
	Update(gc *game.GameContext, delta float64)
	Draw(gc *game.GameContext)
	GetTotalParticles() int
	GetParticles() []*Particle
	RemoveParticle(particle *Particle)
}

type Particle struct {
	base.ObjectEntity
	Symbol    []rune
	Direction Direction
	Style     tcell.Style
}

type ParticleSystem struct {
	ParticleProducable []ParticleProducable
}

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

type ParticleOption func(*Particle)

func WithDimensions(x, y float64, width, height int) ParticleOption {
	return func(p *Particle) {
		centerX := int(x) + width/2
		centerY := int(y) + height/2

		p.Position = base.PointFloat{X: float64(centerX), Y: float64(centerY)}
	}
}

func WithSpeed(initialSpeed int) ParticleOption {
	return func(p *Particle) {
		p.Speed = float64(initialSpeed)
	}
}

func WithSymbols(symbols []rune) ParticleOption {
	return func(p *Particle) {
		p.Symbol = symbols
	}
}

func WithStyle(style tcell.Style) ParticleOption {
	return func(p *Particle) {
		p.Style = style
	}
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{}
}

func (ps *ParticleSystem) AddParticles(particleProducable ParticleProducable) {
	ps.ParticleProducable = append(ps.ParticleProducable, particleProducable)
}

func (ps *ParticleSystem) Update(gc *game.GameContext, delta float64) {
	activeProducers := ps.ParticleProducable[:0]
	for _, p := range ps.ParticleProducable {
		p.Update(gc, delta)
		if p.GetTotalParticles() > 0 {
			activeProducers = append(activeProducers, p)
		}
	}
	ps.ParticleProducable = activeProducers
}

func (ps *ParticleSystem) Draw(gc *game.GameContext) {
	for _, p := range ps.ParticleProducable {
		p.Draw(gc)
	}
}

func (ps *ParticleSystem) InputEvents(event tcell.Event, gc *game.GameContext) {}

func (ps *ParticleSystem) GetType() string {
	return "particles"
}
