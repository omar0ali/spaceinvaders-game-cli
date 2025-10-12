// Package base
package base

import (
	"time"

	"github.com/gdamore/tcell/v2"
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

func (f *ObjectBase) TakeDamage(by int) {
	f.Health -= by
}

func (f *ObjectBase) GetWidth() int {
	return f.Width
}

func (f *ObjectBase) GetHeight() int {
	return f.Height
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
