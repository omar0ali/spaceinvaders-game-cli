package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type UILayoutMenuBoxesProducer struct {
	Boxes         []*Box
	Width, Height int
	SelectedDesc  []string
}

func (u *UILayoutMenuBoxesProducer) GetTotalBoxes() int {
	return len(u.Boxes)
}

func (u *UILayoutMenuBoxesProducer) GetBoxes() []*Box {
	return u.Boxes
}

func InitMainMenu(boxWidth, boxHeight int, boxes ...*Box) *UILayoutMenuBoxesProducer {
	return &UILayoutMenuBoxesProducer{
		Boxes:  boxes,
		Width:  boxWidth,
		Height: boxHeight,
		SelectedDesc: []string{
			"* Space Invaders Game v1.8.0.alpha.2",
			"The game is an endless space shooter where players face increasingly difficult",
			"waves of alien ships that scale with their level.",
			"",
			"Each time the player levels up, they can choose an upgrade to improve their spaceship,",
			"such as boosting firepower to handle tougher aliens with stronger armor.",
			"",
			"The objective is to survive as long as possible, destroy alien ships, and push for",
			"a higher score while managing health through occasional drop-down health packs that",
			"restore the spaceship health.",
			"",
			"(*) Controls",
			"",
			"[LM] hold to shoot a beam to coming alien-ships.",
			"[E] Consume Health Kit.",
			"[R] or [RM] Reload Gun.",
			"[P] To Pause The Game.",
			"[Ctrl+R] To Restart The Game.",
		},
	}
}

func (u *UILayoutMenuBoxesProducer) Update(gc *game.GameContext, delta float64) {}

func (u *UILayoutMenuBoxesProducer) InputEvents(events tcell.Event, gc *game.GameContext) {
	switch ev := events.(type) {
	case *tcell.EventMouse:
		mx, my := ev.Position()
		for _, b := range u.Boxes {
			if mx >= b.Position.X && mx < b.Position.X+b.Width && my >= b.Position.Y && my < b.Position.Y+b.Height {
				if !b.Hovered {
					b.Hovered = true
					if len(b.Description) > 0 {
						u.SelectedDesc = b.Description
					}
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

func (u *UILayoutMenuBoxesProducer) Draw(gc *game.GameContext) {
	w, h := base.GetSize()

	const spaceBetween = 1
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

	// Description Box
	style := base.StyleIt(tcell.ColorWhite)
	width := 90
	height := len(u.SelectedDesc) + 2
	DrawBoxOverlap(base.Point{X: (w / 2) - (width / 2) + 20, Y: (h / 2) - height/2}, width, height, func(innerX, innerY int) {
		for j, line := range u.SelectedDesc {
			for i, r := range line {
				base.SetContentWithStyle(innerX+i+1, innerY+j+1, r, style)
			}
		}
	}, style)
}
