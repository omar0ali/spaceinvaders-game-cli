package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type Beam struct {
	position core.Point
	Symbol   rune
}

type Gun struct {
	Beams []*Beam
	Cap   int
	Power int
	Speed int
}

func (g *Gun) initBeam(pos core.Point) {
	beam := Beam{
		core.Point{
			X: pos.X,
			Y: pos.Y - 2,
		},
		tcell.RuneVLine,
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

	var activeBeams []*Beam

	for _, beam := range g.Beams {
		distance := int(float64(g.Speed) * delta)
		beam.position.Y -= distance

		if beam.position.Y >= 0 {
			activeBeams = append(activeBeams, beam)
		}
	}

	g.Beams = activeBeams
}

func (g *Gun) Draw(gc *core.GameContext) {
	// draw the beam new position
	if len(g.Beams) < 1 {
		return
	}
	for _, beam := range g.Beams {
		window.SetContent(beam.position.X, beam.position.Y, beam.Symbol)
	}
}

func (g *Gun) InputEvents(event tcell.Event, gc *core.GameContext) {
	if len(g.Beams) > g.Cap {
		return
	}
	// on click, will create a new beam
	switch ev := event.(type) {
	case *tcell.EventMouse:
		if ev.Buttons() == tcell.Button1 {
			// limit many beams shots
			if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
				g.initBeam(spaceship.origin)
			}
		}
	case *tcell.EventKey:
		if ev.Rune() == ' ' {
			if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
				g.initBeam(spaceship.origin)
			}
		}
	}
}

func (g *Gun) GetType() string {
	return "gun"
}
