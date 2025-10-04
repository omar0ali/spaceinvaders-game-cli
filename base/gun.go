package base

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Direction = int

const (
	Up Direction = iota
	Down
)

type beam struct {
	position  game.Point
	Symbol    rune
	Direction Direction
}

func (b beam) GetPosition() *game.Point {
	return &b.position
}

type Gun struct {
	beams []*beam
	cap   int
	power int
	speed int
}

func NewGun(cap, power, speed int) Gun {
	return Gun{
		beams: []*beam{},
		cap:   cap,
		power: power,
		speed: speed,
	}
}

func (g Gun) GetPower() int {
	return g.power
}

func (g Gun) GetSpeed() int {
	return g.speed
}

func (g *Gun) IncreaseGunSpeed(i, limit int) bool {
	if g.speed < limit {
		g.speed += i
		g.speed = min(g.speed, limit)
		return true
	}
	return false
}

func (g *Gun) IncreaseGunPower(i int) bool {
	g.power += i
	return true
}

func (g *Gun) IncreaseGunCap(i, limit int) bool {
	if g.cap < limit {
		g.cap += i
		g.cap = min(g.cap, limit)
		return true
	}
	return false
}

func (g Gun) GetCapacity() int {
	return g.cap
}

func (g Gun) GetBeams() []*beam {
	return g.beams
}

func (g *Gun) InitBeam(pos game.Point, dir Direction) {
	if len(g.beams) >= g.cap {
		return
	}

	symbol := '↑'
	if dir == Down {
		symbol = '↓'
	}

	beam := beam{
		game.Point{
			X: pos.X,
			Y: pos.Y,
		},
		symbol,
		dir,
	}
	g.beams = append(g.beams, &beam)
}

func (g *Gun) RemoveBeam(beam *beam) {
	for i, b := range g.beams {
		if beam == b {
			g.beams = append(g.beams[:i], g.beams[i+1:]...)
			break
		}
	}
}

func (g *Gun) Update(gc *game.GameContext, delta float64) {
	// update the coordinates of the beam
	_, h := window.GetSize()
	var activeBeams []*beam
	for _, beam := range g.beams {
		distance := int(float64(g.speed) * delta)
		switch beam.Direction {
		case Up:
			beam.position.Y -= distance
		case Down:
			beam.position.Y += distance
		}
		if beam.position.Y >= 0 && beam.position.Y <= h {
			activeBeams = append(activeBeams, beam)
		}
	}

	g.beams = activeBeams
}

func (g *Gun) Draw(gc *game.GameContext, color tcell.Color) {
	// draw the beam new position
	style := window.StyleIt(tcell.ColorReset, color)

	for _, beam := range g.beams {
		window.SetContentWithStyle(beam.position.X, beam.position.Y, beam.Symbol, style)
	}
}

func (g *Gun) InputEvents(event tcell.Event, gc *game.GameContext) {}

func (g *Gun) GetType() string {
	return "gun"
}
