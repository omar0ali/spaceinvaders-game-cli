package entities

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type UI struct {
	MenuScreen     bool
	PauseScreen    bool
	GameOverScreen bool
}

func (u *UI) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	// start screen
	if u.MenuScreen {
		u.MessageBox(window.GetCenterPoint(),
			`
				[Introduction]
				The game starts at wave 1, and each subsequent wave will increase the number of 
				alien ships and their power. The number of waves is endless, and as the waves 
				progress, your score will increase. You can also collect loot boxes to gain extra
				health or enhance the power of your beams, allowing you to destroy the alien ships
				more quickly.

				[Controls]
				[LM] Click to shoot coming alienships.
				[P] To pause the game.
				[Q] To quit the game.

				Press [S] to start
			`,
			"Space Invaders Game")
		return
	}
	// pause ui
	if u.PauseScreen {
		if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {

			u.MessageBox(
				window.GetCenterPoint(),
				fmt.Sprintf(`
				----------- PAUSED -----------
				- HP: %d
				- Score: %d
				- Kills: %d
				- Level: %d

				[R] To restart the game.
				[Q] To quit the game.
				[P] To continue the game.
			`, spaceship.Health, spaceship.Score, spaceship.Kills, spaceship.Kills),
				"Paused",
			)
			return
		}
	}
	// game over ui
	if u.GameOverScreen {
		if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			u.MessageBox(
				window.GetCenterPoint(),
				fmt.Sprintf(`
				Thank you for playing :)

				- Score: %d
				- Kills: %d
				- Level: %d

				Would you like to play again?

				[R] To restart the game.
				[Q] To quit the game.

			`, spaceship.Score, spaceship.Kills, spaceship.Level),
				"Game Over",
			)
			return
		}
	}
	// show controls at the bottom of the screen
	_, h := window.GetSize()
	for i, r := range []rune("[LM] Shoot Beams ◆ [Q] Quit ◆ [P] Pause Game ◆ [R] Restart Game") {
		window.SetContentWithStyle(0+i, h-1, r, whiteColor)
	}
	// display spacehsip details
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship.UISpaceshipData(gc)
	}
}

func (u *UI) Update(gc *core.GameContext, delta float64) {
	if u.MenuScreen || u.PauseScreen || u.GameOverScreen {
		gc.Halt = true
	} else {
		gc.Halt = false
	}
}

func (u *UI) InputEvents(events tcell.Event, gc *core.GameContext) {
	switch ev := events.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 's' || ev.Rune() == 'S' {
			if u.MenuScreen {
				u.MenuScreen = false
			}
		}
		if ev.Rune() == 'p' || ev.Rune() == 'P' {
			if u.MenuScreen || u.GameOverScreen { // skip
				return
			}
			u.PauseScreen = !u.PauseScreen
		}
	}
}

func (u *UI) GetType() string {
	return "ui"
}

func (u *UI) MessageBox(origin core.Point, message string, title string) {
	padding := 2
	wrappedLines := u.wrapText(message)

	contentWidth := 0
	for _, line := range wrappedLines {
		if len(line) > contentWidth {
			contentWidth = len(line)
		}
	}

	boxWidth := contentWidth + (padding * 2)
	boxHeight := len(wrappedLines) + (padding * 2)

	// update x and y position based on the maxWidth.
	// start position should be moved slightly to the left or top from the origin for real centering.
	origin.X = origin.X - (boxWidth / 2)
	origin.Y = origin.Y - (boxHeight / 2)

	// title string inside top border
	titleStr := "-[ " + title + " ]"
	for x := range boxWidth {
		ch := tcell.RuneHLine
		if x < len(titleStr) {
			ch = rune(titleStr[x])
		}
		window.SetContent(int(origin.GetX())+x, origin.Y, ch)
	}

	// corners
	window.SetContent(int(origin.GetX()), int(origin.GetY()), tcell.RuneULCorner)
	window.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY()), tcell.RuneURCorner)
	window.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY())+boxHeight-1, tcell.RuneLRCorner)
	window.SetContent(int(origin.GetX()), int(origin.GetY())+boxHeight-1, tcell.RuneLLCorner)

	// sides and inner space
	for y := 1; y < boxHeight-1; y++ {
		// left side
		window.SetContent(int(origin.GetX()), int(origin.GetY())+y, tcell.RuneVLine)

		// right side
		window.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY())+y, tcell.RuneVLine)

		// inner space
		for x := 1; x < boxWidth-1; x++ {
			window.SetContent(int(origin.GetX())+x, int(origin.GetY())+y, ' ')
		}
	}

	for i, line := range wrappedLines {
		for j, r := range line {
			window.SetContent(int(origin.GetX())+padding+j, int(origin.GetY())+padding+i, r)
		}
	}

	// bottom line
	for x := 0; x < boxWidth-2; x++ {
		window.SetContent(int(origin.GetX())+x+1, int(origin.GetY())+boxHeight-1, tcell.RuneHLine)
	}
}

func (u UI) wrapText(message string) []string {
	var lines []string
	messageLines := strings.Split(message, "\n")

	maxWidth := 0 // find the maximum line
	for _, messageLine := range messageLines {
		maxWidth = max(len(messageLine), maxWidth)
	}

	for _, paragraph := range messageLines {
		words := strings.Fields(paragraph)
		var line string

		for _, word := range words {
			if len(line)+len(word)+1 > maxWidth {
				lines = append(lines, strings.TrimSpace(line))
				line = word + " "
			} else {
				line += word + " "
			}
		}

		lines = append(lines, strings.TrimSpace(line))
	}

	return lines
}
