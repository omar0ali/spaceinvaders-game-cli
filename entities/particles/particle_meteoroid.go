package particles

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type MeteroidProducer struct {
	Particles []*Particle
}

func (m *MeteroidProducer) GetTotalParticles() int {
	return len(m.Particles)
}

func (m *MeteroidProducer) RemoveParticle(particle *Particle) {
	for i, meteroid := range m.Particles {
		if particle == meteroid {
			m.Particles = append(m.Particles[:i], m.Particles[i+1:]...)
			break
		}
	}
}

func (m *MeteroidProducer) GetParticles() []*Particle {
	return m.Particles
}

func InitMeteroids(scale int, opts ...ParticleOption) *MeteroidProducer {
	var listOfParticles []*Particle

	po := &Particle{
		ObjectEntity: base.ObjectEntity{
			Position: base.PointFloat{X: 0, Y: 0},
			Speed:    float64(rand.Intn(10) + 3),
		},
		Style:  base.StyleIt(tcell.ColorWhite),
		Symbol: []rune("O○o○"),
	}

	for _, o := range opts {
		o(po)
	}

	for dir := Up; dir <= DownRight; dir++ { // for each direction
		for i := range scale {
			particle := &Particle{
				ObjectEntity: base.ObjectEntity{
					Speed:    po.Speed + float64(i)*3,
					Position: po.Position,
				},
				Symbol:    po.Symbol,
				Direction: dir,
				Style:     po.Style,
			}
			listOfParticles = append(listOfParticles, particle)
		}
	}

	return &MeteroidProducer{
		Particles: listOfParticles,
	}
}

func (m *MeteroidProducer) Update(gc *game.GameContext, delta float64) {
	activeParticles := m.Particles[:0] // reset
	w, h := base.GetSize()
	for _, p := range m.Particles {
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

		// ensure to not include any particles that aren't on screen dimension.
		if int(p.Position.X) >= 0 && int(p.Position.X) < w &&
			int(p.Position.Y) >= 0 && int(p.Position.Y) < h {
			activeParticles = append(activeParticles, p)
		}

		// move meteroids towards the spaceship
		p.AppendPositionY(float64(p.Speed-2) * delta)
	}

	m.Particles = activeParticles
}

func (m *MeteroidProducer) Draw(gc *game.GameContext) {
	for _, p := range m.Particles {
		t := time.Now().UnixNano() / int64(time.Millisecond)
		idx := int(t/250) % len(p.Symbol)
		symbol := p.Symbol[idx]
		base.SetContentWithStyle(int(p.Position.X), int(p.Position.Y), symbol, p.Style)
	}
}
