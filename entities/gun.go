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
		tcell.RuneDiamond,
		power,
	}
	g.Beams = append(g.Beams, &beam)
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
			for _, entity := range gc.GetEntities() {
				if entity.GetType() == "spaceship" {
					spaceShip, ok := entity.(*SpaceShip)
					if ok {
						g.initBeam(20, 50, spaceShip.origin)
						break
					}
				}
			}
		}
	}
}

func (g *Gun) GetType() string {
	return "gun"
}
