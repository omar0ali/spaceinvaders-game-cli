// Package base
package base

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	HealthBoxStyle      = '■'
	HealthBoxEmptyStyle = '□'
)

type Meter interface {
	GetCurrent() int
	GetMax() int
}

type BarOptions struct {
	X, Y      int
	Size      int
	ShowStats bool
	InPercent bool
	Style     tcell.Style
	Gun       Gunner
}

type BarOption func(opts *BarOptions)

func WithPosition(x, y int) BarOption {
	return func(o *BarOptions) {
		o.X = x
		o.Y = y
	}
}

func WithStatus(inPercent bool) BarOption {
	return func(o *BarOptions) {
		o.ShowStats = true
		o.InPercent = inPercent
	}
}

func WithBarSize(size int) BarOption {
	return func(o *BarOptions) {
		o.Size = size
	}
}

func WithStyle(style tcell.Style) BarOption {
	return func(o *BarOptions) {
		o.Style = style
	}
}

func WithGun(gun Gunner) BarOption {
	return func(o *BarOptions) {
		o.Gun = gun
	}
}

func DisplayBar(h Meter, opts ...BarOption) {
	defStyle := StyleIt(tcell.ColorReset, tcell.ColorWhite)
	o := BarOptions{
		Size:      10,
		X:         0,
		Y:         0,
		ShowStats: false,
		InPercent: false,
		Style:     defStyle,
		Gun:       nil,
	}

	for _, opt := range opts {
		opt(&o)
	}

	reloadAnimation := []rune("•○")
	if o.Gun != nil && o.Gun.IsReloading() {
		frame := int(time.Now().UnixNano()/300_000_000) % len(reloadAnimation)
		SetContentWithStyle(o.X-2, o.Y, reloadAnimation[frame], o.Style)
	}

	trackXPossition := o.X
	// pre draw health
	for _, r := range string("[") {
		SetContentWithStyle(trackXPossition, o.Y, r, o.Style)
		trackXPossition++
	}
	// draw health
	ratio := float64(h.GetCurrent()) / float64(h.GetMax())
	filled := int(ratio * float64(o.Size))

	for i := range o.Size {
		if i < filled {
			SetContentWithStyle(trackXPossition+i, o.Y, HealthBoxStyle, o.Style)
		} else {
			SetContentWithStyle(trackXPossition+i, o.Y, HealthBoxEmptyStyle, o.Style)
		}
	}
	SetContentWithStyle(trackXPossition+o.Size, o.Y, ']', o.Style)
	if o.ShowStats {
		// or end with showing stats (total health)
		trackXPossition += o.Size

		if o.InPercent {
			for i, r := range []rune(fmt.Sprintf(" %2.f%%", (float64(h.GetCurrent())/float64(h.GetMax()))*100)) {
				SetContentWithStyle(trackXPossition+i+1, o.Y, r, o.Style)
			}

			return
		}
		for i, r := range []rune(fmt.Sprintf(" %d/%d", h.GetCurrent(), h.GetMax())) {
			SetContentWithStyle(trackXPossition+i+1, o.Y, r, o.Style)
		}
	}
}
