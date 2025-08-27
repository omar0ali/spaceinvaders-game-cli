package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
)

type Beam struct {
	speed    int
	position core.Point
	symbol   rune
}

type Gun struct {
	Beams []*Beam
}

func (g *Gun) initBeam(speed int, pos core.Point) {
	beam := Beam{
		speed, pos, 'X',
	}
	g.Beams = append(g.Beams, &beam)
}

func (g *Gun) Update(gc *core.GameContext, delta float64) {
	// update the coordinates of the beam
}

func (g *Gun) Draw(gc *core.GameContext) {
	// draw the beam new position
}

func (g *Gun) InputEvents(event tcell.Event, gc *core.GameContext) {
	// on click, will create a new beam
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Key() == ' ' {
			for _, entity := range gc.GetEntities() {
				if entity.GetType() == "spaceship" {
					spaceShip, ok := entity.(*SpaceShip)
					if ok {
						g.initBeam(10, spaceShip.origin)
						break
					}
				}
			}
		}
	}
}

func (g *Gun) GetType() string {
	return "Gun"
}
