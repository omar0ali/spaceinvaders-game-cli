// Package base
package base

import (
	"fmt"
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

const (
	HealthBoxStyle      = '■'
	HealthBoxEmptyStyle = '□'
)

type HealthBar interface {
	GetHealth() int
	GetMaxHealth() int
}

type ObjectBase struct {
	Health        int
	MaxHealth     int
	Position      game.PointFloat
	Width, Height int
	Speed         float64
}

type FallingObjectBase struct {
	ObjectBase
}

func (f *ObjectBase) GetHealth() int {
	return f.Health
}

func (f *ObjectBase) GetMaxHealth() int {
	return f.MaxHealth
}

func (f *ObjectBase) IsOffScreen(h int) bool {
	return int(f.Position.Y) > h-2
}

func (f *ObjectBase) IsDead() bool {
	return f.Health <= 0
}

func (f *ObjectBase) IsHit(pointBeam game.PointInterface, power int) bool {
	grayColor := window.StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)

	// draw flash when hitting
	pattern := []struct {
		dx, dy int
		r      rune
		style  tcell.Style
	}{
		{0, 0, '*', yellowColor},
		{-1, 0, '-', yellowColor},
		{1, 0, '-', yellowColor},
		{0, -1, '|', grayColor},
		{0, 1, '|', grayColor},
		{-1, -1, '\\', grayColor},
		{1, -1, '/', grayColor},
		{-1, 1, '/', redColor},
		{1, 1, '\\', grayColor},
	}

	px := int(math.Round(pointBeam.GetX()))
	py := int(math.Round(pointBeam.GetY()))
	ox := int(math.Round(f.Position.X))
	oy := int(math.Round(f.Position.Y))

	if px >= ox && px < ox+f.Width &&
		py >= oy && py < oy+f.Height {

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

func (f *FallingObjectBase) Move(delta float64) {
	distance := f.Speed * delta
	f.Position.AppendY(distance)
}

func MoveTo(from, to *ObjectBase, delta float64, gc *game.GameContext) {
	distance := from.Speed * delta

	const toleranceX = 2
	const toleranceY = 5

	bossCenterX := from.Position.X + float64(from.Width)/2
	shipCenterX := to.Position.X + float64(to.Width)/2 + 2

	if math.Abs(bossCenterX-shipCenterX) > toleranceX {
		if bossCenterX > shipCenterX {
			from.Position.AppendX(-distance)
		} else {
			from.Position.AppendX(distance)
		}
	}
	targetY := to.Position.Y - float64(to.Height) - 18
	if targetY < -5 {
		targetY = -5
	}

	if math.Abs(from.Position.Y-targetY) > toleranceY {
		if from.Position.Y > targetY {
			from.Position.AppendY(-distance)
		} else {
			from.Position.AppendY(distance)
		}
	}
}

func (f *ObjectBase) DisplayHealth(barSize int, showStats bool, style tcell.Style) {
	DisplayHealth(
		int(f.Position.GetX())+(f.Width/2)-(barSize/2)-1,
		int(f.Position.GetY()-1),
		barSize,
		f,
		showStats,
		style,
	)
}

func DisplayHealth(xPos, yPos, barSize int, h HealthBar, showStats bool, style tcell.Style) {
	trackXPossition := xPos
	// pre draw health
	for _, r := range string("[") {
		window.SetContentWithStyle(trackXPossition, yPos, r, style)
		trackXPossition++
	}
	// draw health
	ratio := float64(h.GetHealth()) / float64(h.GetMaxHealth())
	filled := int(ratio * float64(barSize))

	for i := range barSize {
		if i < filled {
			window.SetContentWithStyle(trackXPossition+i, yPos, HealthBoxStyle, style)
		} else {
			window.SetContentWithStyle(trackXPossition+i, yPos, HealthBoxEmptyStyle, style)
		}
	}
	if !showStats {
		// end with a bracket
		window.SetContentWithStyle(trackXPossition+barSize, yPos, ']', style)
		return
	}
	// or end with showing stats (total health)
	trackXPossition += barSize
	// last
	for i, r := range []rune(fmt.Sprintf("] %d/%d", h.GetHealth(), h.GetMaxHealth())) {
		window.SetContentWithStyle(trackXPossition+i, yPos, r, style)
	}
}

// Will use this for i.e alien ship shooting every # seconds

func DoEvery(interval time.Duration, fn func(), done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fn()
		case <-done:
			return
		}
	}
}
