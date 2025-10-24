// Package base
package base

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type ObjectEntity struct {
	Position      PointFloat
	Width, Height int
	Speed         float64
}

type ObjectBase struct {
	ObjectEntity
	Health    int
	MaxHealth int
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

func (f *ObjectEntity) GetWidth() int {
	return f.Width
}

func (f *ObjectEntity) GetHeight() int {
	return f.Height
}

func (f *ObjectEntity) GetSpeed() float64 {
	return f.Speed
}

func (f *ObjectEntity) GetPosition() *PointFloat {
	return &f.Position
}

func (f *ObjectEntity) AppendPositionY(y float64) {
	f.Position.AppendY(y)
}

func DisplayHealthLeft(base *ObjectBase, y int, name string, barSize int, showPercentage bool, style tcell.Style, gun *Gun) {
	for i, r := range name {
		SetContentWithStyle(2+i, y, r, style)
	}

	DisplayBar(
		base,
		WithGun(gun),
		WithBarSize(barSize),
		WithPosition(2, y+1),
		WithStatus(true),
		WithStyle(style),
	)
}

func DisplayHealthTop(base *ObjectBase, name string, barSize int, showPercentage bool, style tcell.Style, gun *Gun) {
	w, _ := GetSize()
	for i, r := range name {
		SetContentWithStyle((w/2)-(len(name)/2)+i, 0, r, style)
	}

	DisplayBar(
		base,
		WithGun(gun),
		WithBarSize(barSize),
		WithPosition((w/2)-(barSize+1)/2, 1),
		WithStatus(true),
		WithStyle(style),
	)
}

func (f *ObjectBase) DisplayHealth(barSize int, style tcell.Style, gun *Gun) {
	DisplayBar(
		f,
		WithPosition(
			int(f.Position.GetX())+(f.Width/2)-(barSize/2)-1,
			int(f.Position.GetY()-1),
		),
		WithBarSize(barSize),
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
