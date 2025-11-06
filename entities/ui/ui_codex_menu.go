package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type UICodexMenuBoxesProducer struct {
	UIProducerBase
	CurrentDisplayList []*Box
}

func (u *UICodexMenuBoxesProducer) GetTotalBoxes() int {
	return len(u.Boxes)
}

func (u *UICodexMenuBoxesProducer) GetBoxes() []*Box {
	return u.Boxes
}

func InitCodexMenu(boxWidth, boxHeight int) *UICodexMenuBoxesProducer {
	return &UICodexMenuBoxesProducer{
		UIProducerBase: UIProducerBase{
			Width:  boxWidth,
			Height: boxHeight,
			SelectedDesc: []string{
				"Codex Menu",
			},
		},
	}
}

func (u *UICodexMenuBoxesProducer) SetMenuItems(boxes []*Box) {
	u.Boxes = boxes
}

func (u *UICodexMenuBoxesProducer) SetList(boxes []*Box) {
	u.CurrentDisplayList = boxes
}

func (u *UICodexMenuBoxesProducer) Update(gc *game.GameContext, delta float64) {}

func (u *UICodexMenuBoxesProducer) InputEvents(events tcell.Event, gc *game.GameContext) {
	switch ev := events.(type) {
	case *tcell.EventMouse:
		mx, my := ev.Position()
		for _, b := range u.Boxes {
			if mx >= b.Position.X && mx < b.Position.X+b.Width && my >= b.Position.Y && my < b.Position.Y+b.Height {
				if !b.Hovered {
					gc.Sounds.PlaySound("8-bit-hover-button.mp3", 0)
					b.Hovered = true
					if len(b.Description) > 0 {
						u.SelectedDesc = b.Description
					}
				}
				if ev.Buttons() == tcell.Button1 {
					gc.Sounds.PlaySound("8-bit-game-sfx-sound-select.mp3", 0)
					if b.OnClick != nil {
						b.OnClick()
					}
				}
			} else {
				b.Hovered = false
			}
		}
		for _, b := range u.CurrentDisplayList {
			if mx >= b.Position.X && mx < b.Position.X+b.Width && my >= b.Position.Y && my < b.Position.Y+b.Height {
				if !b.Hovered {
					gc.Sounds.PlaySound("8-bit-hover-button.mp3", 0)
					b.Hovered = true
					if len(b.Description) > 0 {
						u.SelectedDesc = b.Description
					}
				}
				if ev.Buttons() == tcell.Button1 {
					gc.Sounds.PlaySound("8-bit-game-sfx-sound-select.mp3", 0)
					if b.OnClick != nil {
						b.OnClick()
					}
				}
			} else {
				b.Hovered = false
			}
		}
	}
}

func (u *UICodexMenuBoxesProducer) Draw(gc *game.GameContext) {
	w, h := base.GetSize()

	const spaceBetween = 0
	totalHeight := (u.Height * len(u.Boxes)) + (spaceBetween * (len(u.Boxes) - 1))
	startX := 20
	startY := (h / 2) - (totalHeight / 2)
	for i, b := range u.Boxes {
		b.Position.X = startX
		b.Position.Y = startY + (i * (u.Height + spaceBetween))
		b.Width = u.Width
		b.Height = u.Height

		// style based on hover
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		if b.Hovered {
			style = tcell.StyleDefault.Foreground(tcell.ColorYellow)
		}

		// draw the box using DrawRect
		DrawBoxHover(base.Point{X: b.Position.X, Y: b.Position.Y}, b.Width, b.Height, b.Hovered, func(innerX, innerY int) {
			startX := ((b.Width / 2) - 2) - len([]rune(b.Shape[0]))/2
			startY := (b.Height / 2) - 2
			for rowIndex, line := range b.Shape {
				for colIndex, char := range line {
					if char != ' ' {
						base.SetContentWithStyle(innerX+colIndex+startX, innerY+rowIndex+startY, char, style)
					}
				}
			}
		})
	}

	// Items
	startPosition := 0
	gridAt := 1
	for _, b := range u.CurrentDisplayList {
		if startX+25+((startPosition+1)*(u.Width+spaceBetween)) > w {
			gridAt++
			startPosition = 0
		}

		b.Position.X = startX + 25 + (startPosition * (u.Width + spaceBetween))
		b.Position.Y = (gridAt * (u.Height * 2)) - 10 // fixed: base Y per grid row
		b.Width = u.Width
		b.Height = u.Height * 2

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
		startPosition++
	}

	// Description Box
	style := base.StyleIt(tcell.ColorWhite)
	width := w / 2
	height := len(u.SelectedDesc) + 2
	DrawBoxOverlap(base.Point{X: (w / 2) - (width / 2) + 20, Y: h - height}, width, height, func(innerX, innerY int) {
		for j, line := range u.SelectedDesc {
			for i, r := range line {
				base.SetContentWithStyle(innerX+i+1, innerY+j+1, r, style)
			}
		}
	}, style)
}
