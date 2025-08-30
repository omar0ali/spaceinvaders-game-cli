package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type Beam struct {
	speed    int
	position core.Point
	symbol   rune
	power    int
}

type Gun struct {
	Beams []*Beam
}

func (g *Gun) initBeam(power, speed int, pos core.Point) {
	beam := Beam{
		speed,
		core.Point{
			X: pos.X,
			Y: pos.Y - 2,
		},
		tcell.RuneVLine,
		power,
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
		distance := int(float64(beam.speed) * delta)
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
		window.SetContent(beam.position.X, beam.position.Y, beam.symbol)
	}
}

func (g *Gun) InputEvents(event tcell.Event, gc *core.GameContext) {
	// on click, will create a new beam
	switch ev := event.(type) {
	case *tcell.EventMouse:
		if ev.Buttons() == tcell.Button1 {
			// limit how many beams shot
			// number of beams can't exceed 10
			// TODO: this can be changed later when implementing config file.
			if len(g.Beams) > 10 {
				return
			}

			if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
				g.initBeam(20, 40, spaceship.origin)
			}
		}
	}
}

func (g *Gun) GetType() string {
	return "gun"
}
