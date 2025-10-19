package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type UIProducable interface {
	Update(gc *game.GameContext, delta float64)
	Draw(gc *game.GameContext)
	InputEvents(events tcell.Event, gc *game.GameContext)
	GetTotalBoxes() int
	GetBoxes() []*Box
}

type UISystem struct {
	UIProducable UIProducable
	target       time.Time // this is used to give time before player can use the menu
	// this used to avoiud accidents
}

func NewUISystem() *UISystem {
	return &UISystem{}
}

type Box struct {
	Position    base.Point
	Height      int
	Width       int
	Shape       []string
	Description []string
	OnClick     func()
	Hovered     bool
}

func NewUIBox(shape, desc []string, onClick func()) *Box {
	return &Box{
		Description: desc,
		Shape:       shape,
		OnClick:     onClick,
	}
}

func (ui *UISystem) SetLayout(layout UIProducable) {
	now := time.Now()
	ui.target = now.Add(2 * time.Second)

	ui.UIProducable = layout
}

func (ui *UISystem) Draw(gc *game.GameContext) {
	// show loading symbol
	if !time.Now().After(ui.target) && ui.UIProducable != nil {
		w, h := base.GetSize()
		// frames := []rune{'.', 'o', 'O', 'o'}
		// frames := []rune{'▁', '▃', '▄', '▅', '▆', '▇', '█', '▇', '▆', '▅', '▄', '▃'}
		frames := []rune{'•', '◦'}

		i := int(time.Now().UnixNano()/400_000_000) % len(frames)
		j := int(time.Now().UnixNano()/200_000_000) % len(frames)

		loadingWord := "Loading"
		for i, r := range loadingWord {
			base.SetContent(i+(w/2)-len(loadingWord)/2, h-3, r)
		}

		base.SetContent((w/2)-1, h-2, frames[i])
		base.SetContent((w/2)-2, h-2, frames[i])
		base.SetContent(w/2, h-2, frames[j])
		base.SetContent((w/2)+1, h-2, frames[i])
		base.SetContent((w/2)+2, h-2, frames[i])
	}

	if ui.UIProducable != nil {
		ui.UIProducable.Draw(gc)
	}
}

func (ui *UISystem) Update(gc *game.GameContext, delta float64) {
	if ui.UIProducable != nil {
		ui.UIProducable.Update(gc, delta)
	}
}

func (ui *UISystem) InputEvents(events tcell.Event, gc *game.GameContext) {
	if time.Now().After(ui.target) {
		if ui.UIProducable != nil {
			ui.UIProducable.InputEvents(events, gc)
		}
	}
}

func (ui *UISystem) GetType() string {
	return "layout"
}

func DrawBoxHover(pos base.Point, width, height int, hover bool, fn func(initX, initY int)) {
	style := base.StyleIt(tcell.ColorWhite)
	if hover {
		style = base.StyleIt(tcell.ColorYellowGreen)
	}
	const padding = 2
	startX := pos.X
	startY := pos.Y
	for i := range height {
		for j := range width {
			switch {
			case j == 0 && i == 0:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneULCorner, style)
			case j == width-1 && i == 0:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneURCorner, style)
			case j == 0 && i == height-1:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneLLCorner, style)
			case j == width-1 && i == height-1:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneLRCorner, style)
			case i == 0 || i == height-1:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneHLine, style)
			case j == 0 || j == width-1:
				base.SetContentWithStyle(startX+j, startY+i, tcell.RuneVLine, style)

			default:
				base.SetContent(startX+j, startY+i, ' ')
			}
		}
	}
	fn(startX+padding, startY+padding)
}

func DrawBox(pos base.Point, width, height int, style tcell.Style) {
	for i := range height {
		for j := range width {
			switch {
			case j == 0 && i == 0:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneULCorner, style)
			case j == width-1 && i == 0:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneURCorner, style)
			case j == 0 && i == height-1:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneLLCorner, style)
			case j == width-1 && i == height-1:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneLRCorner, style)
			case i == 0 || i == height-1:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneHLine, style)
			case j == 0 || j == width-1:
				base.SetContentWithStyle(pos.X+j, pos.Y+i, tcell.RuneVLine, style)

			default:
				base.SetContent(pos.X+j, pos.Y+i, ' ')
			}
		}
	}
}

func DrawBoxOverlap(pos base.Point, width, height int, fn func(initX, initY int), style tcell.Style) {
	const padding = 2
	startX := pos.X
	startY := pos.Y
	DrawBox(pos, width, height, style)
	fn(startX, startY)
}

// TODO: Refactor

func DrawRect(pos base.Point, width, height int, fn func(initX, initY int)) {
	const padding = 2
	centerOfW := pos.X / 2
	centerOfH := pos.Y / 2
	startX := centerOfW - (width / 2)
	startY := centerOfH - (height / 2)
	for i := range height {
		for j := range width {
			switch {
			case j == 0 && i == 0:
				base.SetContent(startX+j, startY+i, tcell.RuneULCorner)
			case j == width-1 && i == 0:
				base.SetContent(startX+j, startY+i, tcell.RuneURCorner)
			case j == 0 && i == height-1:
				base.SetContent(startX+j, startY+i, tcell.RuneLLCorner)
			case j == width-1 && i == height-1:
				base.SetContent(startX+j, startY+i, tcell.RuneLRCorner)
			case i == 0 || i == height-1:
				base.SetContent(startX+j, startY+i, tcell.RuneHLine)
			case j == 0 || j == width-1:
				base.SetContent(startX+j, startY+i, tcell.RuneVLine)

			default:
				base.SetContent(startX+j, startY+i, ' ')
			}
		}
	}
	fn(startX, startY)
}
