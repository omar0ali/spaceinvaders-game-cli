// Package base
package base

import (
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type ObjectBase struct {
	Health        int
	MaxHealth     int
	Position      PointFloat
	Width, Height int
	Speed         float64
}

type FallingObjectBase struct {
	ObjectBase
}

func (f *ObjectBase) GetCurrent() int {
	return f.Health
}

func (f *ObjectBase) GetMax() int {
	return f.MaxHealth
}

func (f *ObjectBase) IsOffScreen(h int) bool {
	return int(f.Position.Y) > h-2
}

func (f *ObjectBase) IsDead() bool {
	return f.Health <= 0
}

func Crash(c1, c2 Movable) bool {
	grayColor := StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := StyleIt(tcell.ColorReset, tcell.ColorYellow)

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

	x1 := int(math.Round(c1.GetPosition().GetX()))
	y1 := int(math.Round(c1.GetPosition().GetY()))
	w1 := c1.GetWidth()
	h1 := c1.GetHeight()

	x2 := int(math.Round(c2.GetPosition().GetX()))
	y2 := int(math.Round(c2.GetPosition().GetY()))
	w2 := c2.GetWidth()
	h2 := c2.GetHeight()

	if x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2 {

		for _, p := range pattern {
			SetContentWithStyle(x1+(w1/2)+3+p.dx, y1+p.dy, p.r, p.style)
			SetContentWithStyle(x2+(w2/2)+3+p.dx, y2+p.dy, p.r, p.style)
		}
		return true
	}

	return false
}

func (f *ObjectBase) TakeDamage(by int) {
	f.Health -= by
}

func GettingHit(m Movable, beam PointableInt) bool {
	grayColor := StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := StyleIt(tcell.ColorReset, tcell.ColorYellow)

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

	px := int(math.Round(beam.GetPosition().GetX()))
	py := int(math.Round(beam.GetPosition().GetY()))
	ox := int(math.Round(m.GetPosition().GetX()))
	oy := int(math.Round(m.GetPosition().GetY()))

	if px >= ox && px < ox+m.GetWidth() &&
		py >= oy && py < oy+m.GetHeight() {

		for _, p := range pattern {
			SetContentWithStyle(
				int(beam.GetPosition().GetX())+p.dx,
				int(beam.GetPosition().GetY())+p.dy,
				p.r, p.style,
			)
		}
		return true
	}
	return false
}

func (f *ObjectBase) GetWidth() int {
	return f.Width
}

func (f *ObjectBase) GetHeight() int {
	return f.Height
}

type PointableInt interface {
	GetPosition() *Point
}

type PointableFloat interface {
	GetPosition() *PointFloat
}

type Sizeable interface {
	GetWidth() int
	GetHeight() int
}

type Movable interface {
	Sizeable
	PointableFloat
	AppendPositionY(float64)
	GetSpeed() float64
}

func (f *ObjectBase) GetSpeed() float64 {
	return f.Speed
}

func (f *ObjectBase) GetPosition() *PointFloat {
	return &f.Position
}

func (f *ObjectBase) AppendPositionY(y float64) {
	f.Position.AppendY(y)
}

func Move(m Movable, delta float64) {
	distance := m.GetSpeed() * delta
	m.AppendPositionY(distance)
}

func MoveTo(from, to Movable, delta float64, gc *game.GameContext) {
	distance := from.GetSpeed() * delta

	const toleranceX = 2
	const toleranceY = 5

	bossCenterX := from.GetPosition().X + float64(from.GetWidth())/2
	shipCenterX := to.GetPosition().X + float64(to.GetWidth())/2 + 2

	if math.Abs(bossCenterX-shipCenterX) > toleranceX {
		if bossCenterX > shipCenterX {
			from.GetPosition().AppendX(-distance)
		} else {
			from.GetPosition().AppendX(distance)
		}
	}

	targetY := to.GetPosition().Y - float64(to.GetHeight()) - 10
	if targetY < -5 {
		targetY = -5
	}

	if math.Abs(from.GetPosition().Y-targetY) > toleranceY {
		if from.GetPosition().Y > targetY {
			from.GetPosition().AppendY(-distance)
		} else {
			from.GetPosition().AppendY(distance)
		}
	}
}

func (f *ObjectBase) DisplayHealth(barSize int, showPercentage bool, style tcell.Style, gun Gunner) {
	DisplayBar(
		f,
		WithPosition(
			int(f.Position.GetX())+(f.Width/2)-(barSize/2)-1,
			int(f.Position.GetY()-1),
		),
		WithBarSize(barSize),
		WithStatus(showPercentage),
		WithStyle(style),
		WithGun(gun),
	)
}

//
// func DisplayBar(xPos, yPos, barSize int, h Meter, showStats bool, inPercent bool, style tcell.Style, gun Gunner) {
// 	reloadAnimation := []rune("•○")
// 	if gun != nil && gun.IsReloading() {
// 		frame := int(time.Now().UnixNano()/300_000_000) % len(reloadAnimation)
// 		SetContentWithStyle(xPos-2, yPos, reloadAnimation[frame], style)
// 	}
//
// 	trackXPossition := xPos
// 	// pre draw health
// 	for _, r := range string("[") {
// 		SetContentWithStyle(trackXPossition, yPos, r, style)
// 		trackXPossition++
// 	}
// 	// draw health
// 	ratio := float64(h.GetCurrent()) / float64(h.GetMax())
// 	filled := int(ratio * float64(barSize))
//
// 	for i := range barSize {
// 		if i < filled {
// 			SetContentWithStyle(trackXPossition+i, yPos, HealthBoxStyle, style)
// 		} else {
// 			SetContentWithStyle(trackXPossition+i, yPos, HealthBoxEmptyStyle, style)
// 		}
// 	}
//
// 	SetContentWithStyle(trackXPossition+barSize, yPos, ']', style)
// 	if showStats {
// 		// or end with showing stats (total health)
// 		trackXPossition += barSize
//
// 		if inPercent {
// 			for i, r := range []rune(fmt.Sprintf(" %2.f%%", (float64(h.GetCurrent())/float64(h.GetMax()))*100)) {
// 				SetContentWithStyle(trackXPossition+i+1, yPos, r, style)
// 			}
//
// 			return
// 		}
// 		for i, r := range []rune(fmt.Sprintf(" %d/%d", h.GetCurrent(), h.GetMax())) {
// 			SetContentWithStyle(trackXPossition+i+1, yPos, r, style)
// 		}
// 	}
// }

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

func DoOnce(delay time.Duration, fn func(), done <-chan struct{}) {
	select {
	case <-time.After(delay):
		fn() // run the function after the delay
	case <-done:
		// stop early if done signal received
		return
	}
}
