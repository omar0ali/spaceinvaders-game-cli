package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Direction int

const (
	Up Direction = iota
	Down
)

type Beam struct {
	position  core.Point
	Symbol    rune
	Direction Direction
}

type Gun struct {
	Beams []*Beam
	Cap   int
	Power int
	Speed int
}

func (g *Gun) initBeam(pos core.Point, dir Direction) {
	beam := Beam{
		core.Point{
			X: pos.X,
			Y: pos.Y - 2,
		},
		tcell.RuneVLine,
		dir,
	}
	g.Beams = append(g.Beams, &beam)
}

func (g *Gun) RemoveBeam(beam *Beam) {
	for i, b := range g.Beams {
		if beam == b {
			g.Beams = append(g.Beams[:i], g.Beams[i+1:]...)
			break
		}
	}
}

func (g *Gun) Update(gc *core.GameContext, delta float64) {
	// update the coordinates of the beam
	if len(g.Beams) < 1 {
		return
	}

	_, h := window.GetSize()

	var activeBeams []*Beam

	for _, beam := range g.Beams {
		distance := int(float64(g.Speed) * delta)
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

	g.Beams = activeBeams
}

func (g *Gun) Draw(gc *core.GameContext) {
	// draw the beam new position
	color := window.StyleIt(tcell.ColorReset, tcell.ColorBlueViolet)

	if len(g.Beams) < 1 {
		return
	}

	for _, beam := range g.Beams {
		if beam.Direction == Down {
			window.SetContentWithStyle(beam.position.X, beam.position.Y, beam.Symbol, color)
		} else {
			window.SetContent(beam.position.X, beam.position.Y, beam.Symbol)
		}
	}
}

func (g *Gun) InputEvents(event tcell.Event, gc *core.GameContext) {}

func (g *Gun) GetType() string {
	return "gun"
}
