// Package ui
package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type UILayoutBoxesProducer struct {
	Boxes         []*Box
	selectedDesc  []string
	Width, Height int
}

func (u *UILayoutBoxesProducer) GetTotalBoxes() int {
	return len(u.Boxes)
}

func (u *UILayoutBoxesProducer) GetBoxes() []*Box {
	return u.Boxes
}

func InitLayout(boxWidth, boxHeight int, boxes ...*Box) *UILayoutBoxesProducer {
	return &UILayoutBoxesProducer{
		Boxes:        boxes,
		Width:        boxWidth,
		Height:       boxHeight,
		selectedDesc: []string{"- No detials to show."},
	}
}

func (u *UILayoutBoxesProducer) Update(gc *game.GameContext, delta float64) {}

func (u *UILayoutBoxesProducer) InputEvents(events tcell.Event, gc *game.GameContext) {
	switch ev := events.(type) {
	case *tcell.EventMouse:
		mx, my := ev.Position()
		for _, b := range u.Boxes {
			if mx >= b.Position.X && mx < b.Position.X+b.Width && my >= b.Position.Y && my < b.Position.Y+b.Height {
				if !b.Hovered {
					b.Hovered = true
					u.selectedDesc = b.Description
				}
				if ev.Buttons() == tcell.Button1 {
					b.OnClick()
				}
			} else {
				b.Hovered = false
			}
		}
	}
}

func (u *UILayoutBoxesProducer) Draw(gc *game.GameContext) {
	w, h := base.GetSize()

	const spaceBetween = 1
	totalWidth := (u.Width * len(u.Boxes)) + (spaceBetween * (len(u.Boxes) - 1))
	startX := (w / 2) - (totalWidth / 2)
	startY := (h / 2) - u.Height
	for i, b := range u.Boxes {
		b.Position.X = startX + (i * (u.Width + spaceBetween))
		b.Position.Y = startY
		b.Width = u.Width
		b.Height = u.Height

		// style based on hover
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		if b.Hovered {
			style = tcell.StyleDefault.Foreground(tcell.ColorYellow)
		}

		// draw the box using DrawRect
		DrawBoxHover(base.Point{X: b.Position.X, Y: b.Position.Y}, b.Width, b.Height, b.Hovered, func(innerX, innerY int) {
			startX := (b.Width / 2) - len(b.Shape[0])/2
			startY := (b.Height / 2) - len(b.Shape) + 1
			for rowIndex, line := range b.Shape {
				for colIndex, char := range line {
					if char != ' ' {
						base.SetContentWithStyle(innerX+colIndex+startX-2, innerY+rowIndex+startY, char, style)
					}
				}
			}
		})
	}

	style := base.StyleIt(tcell.ColorWhite)
	DrawRect(base.Point{X: w, Y: h + 15}, 50, len(u.selectedDesc)+2, func(innerX, innerY int) {
		for j, line := range u.selectedDesc {
			for i, r := range line {
				base.SetContentWithStyle(innerX+i+1, innerY+j+1, r, style)
			}
		}
	})
}
