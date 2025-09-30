package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

const (
	HealthBoxStyle      = '■'
	HealthBoxEmptyStyle = '□'
)

type FallingObjectBase struct {
	Health        int
	Speed         int
	OriginPoint   core.PointFloat
	Width, Height int
}

func (f *FallingObjectBase) isOffScreen(h int) bool {
	return int(f.OriginPoint.Y) > h-2
}

func (f *FallingObjectBase) isDead() bool {
	return f.Health <= 0
}

func (f *FallingObjectBase) isHit(pointBeam core.PointInterface, power int) bool {
	grayColor := window.StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)

	// draw flash when hitting
	pattern := []struct {
		dx, dy int
		r      rune
		style  tcell.Style
	}{
		{-1, 0, tcell.RuneBoard, yellowColor},
		{1, 0, tcell.RuneBoard, yellowColor},
		{0, -1, tcell.RuneBoard, grayColor},
		{0, 1, tcell.RuneBoard, grayColor},
		{-1, -1, tcell.RuneCkBoard, grayColor},
		{1, -1, tcell.RuneCkBoard, grayColor},
		{-1, 1, tcell.RuneCkBoard, redColor},
		{1, 1, tcell.RuneBoard, grayColor},
	}

	if pointBeam.GetX() >= f.OriginPoint.X &&
		pointBeam.GetX() <= f.OriginPoint.X+float64(f.Width) &&
		pointBeam.GetY() >= f.OriginPoint.Y &&
		pointBeam.GetY() <= f.OriginPoint.Y+float64(f.Height-2) {

		f.Health -= power // update health of the falling object

		for _, p := range pattern {
			window.SetContentWithStyle(
				int(pointBeam.GetX())+p.dx,
				int(pointBeam.GetY())+p.dy,
				p.r, p.style,
			)
		}

		return true
	}
	return false
}

func (f *FallingObjectBase) move(delta float64) {
	distance := float64(f.Speed) * delta
	f.OriginPoint.AppendY(distance)
}

type ObjectOpts struct {
	Health        int
	Speed         int
	OriginPoint   core.PointFloat
	Width, Height int
}

type FallingObjects interface {
	move(distance float64)
	isHit(point core.PointInterface)
	NewObject(health, speed int, origin core.PointFloat)
}

type HealthBar interface {
	GetHealth() int
	GetMaxHealth() int
}

func DisplayHealth(xPos, yPos int, h HealthBar) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	trackXPossition := xPos
	// pre draw health
	for _, r := range string("  * [") {
		window.SetContentWithStyle(trackXPossition, yPos, r, whiteColor)
		trackXPossition++
	}
	// draw health
	barSize := 10
	ratio := float64(h.GetHealth()) / float64(h.GetMaxHealth())
	filled := int(ratio * float64(barSize))

	for i := range barSize {
		if i < filled {
			window.SetContentWithStyle(trackXPossition+i, yPos, HealthBoxStyle, whiteColor)
		} else {
			window.SetContentWithStyle(trackXPossition+i, yPos, HealthBoxEmptyStyle, whiteColor)
		}
	}
	trackXPossition += barSize
	// last
	for i, r := range []rune(fmt.Sprintf("] %d/%d", h.GetHealth(), h.GetMaxHealth())) {
		window.SetContentWithStyle(trackXPossition+i, yPos, r, whiteColor)
	}
}
