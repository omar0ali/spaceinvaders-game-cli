package entities

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/ui"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

const (
	IncreaseGunCapBy      = 1
	IncreaseGunPowerBy    = 1
	IncreaseGunSpeedBy    = 1
	DecreaseGunCooldownBy = 3
)

var (
	nextMinute   int
	minutes      int
	seconds      int
	listOfStatus []string
	mu           sync.Mutex
)

type UI struct {
	MenuScreen         bool
	PauseScreen        bool
	GameOverScreen     bool
	LevelUpScreen      bool
	SpaceShipSelection bool
	timeElapsed        float64
}

func NewUI(gc *game.GameContext, exitCha chan struct{}) *UI {
	nextMinute = 0

	u := &UI{
		MenuScreen:         true,
		PauseScreen:        false,
		GameOverScreen:     false,
		LevelUpScreen:      false,
		SpaceShipSelection: false,
	}

	if u.MenuScreen {
		if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {

			boxes := []*ui.Box{
				ui.NewUIBox(
					[]string{
						"Start New Game",
					},
					[]string{},
					func() {
						// here we should start the game
						SetStatus("Select a Spaceship")
						u.SpaceShipSelection = true
						if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
							var boxes []*ui.Box
							for i, shipDesign := range s.ListOfSpaceships {
								descriptions := []string{
									fmt.Sprintf("- [%s]", shipDesign.Name),
									fmt.Sprintf("* HP:         %d", shipDesign.EntityHealth),
									fmt.Sprintf("* Gun PWD:    %d", shipDesign.GunPower),
									fmt.Sprintf("* Gun CAP:    %d", shipDesign.GunCap),
									fmt.Sprintf("* Gun SPD:    %d", shipDesign.GunSpeed),
									fmt.Sprintf("* Gun CD:     %d ms", shipDesign.GunCooldown),
									fmt.Sprintf("* Gun RLD CD: %d ms", shipDesign.GunReloadCooldown),
								}

								boxes = append(boxes, ui.NewUIBox(
									shipDesign.Shape,
									descriptions,
									func() {
										name := s.SpaceshipSelection(i)
										SetStatus(fmt.Sprintf("%s Selected", name))
										u.SpaceShipSelection = false
										layout.SetLayout(nil)
									},
								))
							}
							layout.SetLayout(
								ui.InitLayout(21, 10, boxes...),
							)
						}
						u.MenuScreen = false
					},
				),
				ui.NewUIBox([]string{
					"Quit Game",
				}, []string{}, func() {
					base.ExitGame(exitCha)
				}),
			}
			layout.SetLayout(
				ui.InitMainMenu(20, 5, boxes...),
			)
		}
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			LevelUpPopUp(gc, u, s)
		})
	}
	return u
}

func (u *UI) Draw(gc *game.GameContext) {
	whiteColor := base.StyleIt(tcell.ColorWhite)

	n := len(listOfStatus)
	for i, notification := range listOfStatus {
		yIndex := n - 1 - i // invert y order
		DrawRectStatus(notification, yIndex)
	}

	// show controls at the bottom of the screen
	w, h := base.GetSize()
	controlsUI := []rune("[LM] Shoot Beams ◆ [E] Consume Health Kit ◆ [R] Reload Gun ◆ [P] Pause Game ◆ [Ctrl+R] Restart Game ◆ [Ctrl+Q] Quit")
	for i, r := range controlsUI {
		base.SetContentWithStyle(w/2-(len(controlsUI)/2)+i, h-1, r, whiteColor)
	}

	// timer
	minutes = int(u.timeElapsed) / 60
	seconds = int(u.timeElapsed) % 60

	if !u.MenuScreen {

		timeStr := []rune(fmt.Sprintf("  * Time: %02d:%02d", minutes, seconds))
		// display objects details
		for i, r := range timeStr {
			base.SetContentWithStyle(i, 1, r, whiteColor)
		}

		// display spacehsip details - Also drop a health kit every minute
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			s.UISpaceshipData(gc)
		}

		// display aliens details
		if aliens, ok := gc.FindEntity("alien").(*AlienProducer); ok {
			aliens.UIAlienShipData(gc)
		}

	}
	// pause ui
	if u.PauseScreen {
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {

			u.MessageBox(
				base.GetCenterPoint(),
				fmt.Sprintf(`
				----------- PAUSED -----------
				- HP: %d
				- Score: %d
				- Kills: %d
				- Level: %d

				[Ctrl+R] To restart the game.
				[Ctrl+Q] To quit the game.

				[P] To continue the game.
			`, s.Health, s.Score.Score, s.Kills, s.Level),
				"Paused",
			)
			return
		}
	}
	// game over ui
	if u.GameOverScreen {
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			u.MessageBox(
				base.GetCenterPoint(),
				fmt.Sprintf(`
				Thank you for playing :)

				- Score: %d
				- Kills: %d
				- Level: %d

				Would you like to play again?

				[Ctrl+R] To restart the game.
				[Ctrl+Q] To quit the game.

			`, s.Score.Score, s.Kills, s.Level),
				"Game Over",
			)
		}
		return
	}
}

func (u *UI) Update(gc *game.GameContext, delta float64) {
	if u.MenuScreen || u.PauseScreen || u.GameOverScreen || u.LevelUpScreen || u.SpaceShipSelection {
		gc.Halt = true
	} else {
		gc.Halt = false
		u.timeElapsed += delta
	}
}

func (u *UI) InputEvents(events tcell.Event, gc *game.GameContext) {
	switch ev := events.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'p' || ev.Rune() == 'P' {
			if u.MenuScreen || u.GameOverScreen || u.SpaceShipSelection { // skip
				return
			}
			u.PauseScreen = !u.PauseScreen
		}
	}
}

func (u *UI) GetType() string {
	return "ui"
}

func (u *UI) MessageBox(origin base.Point, message string, title string) {
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
		base.SetContent(int(origin.GetX())+x, origin.Y, ch)
	}

	// corners
	base.SetContent(int(origin.GetX()), int(origin.GetY()), tcell.RuneULCorner)
	base.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY()), tcell.RuneURCorner)
	base.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY())+boxHeight-1, tcell.RuneLRCorner)
	base.SetContent(int(origin.GetX()), int(origin.GetY())+boxHeight-1, tcell.RuneLLCorner)

	// sides and inner space
	for y := 1; y < boxHeight-1; y++ {
		// left side
		base.SetContent(int(origin.GetX()), int(origin.GetY())+y, tcell.RuneVLine)

		// right side
		base.SetContent(int(origin.GetX())+boxWidth-1, int(origin.GetY())+y, tcell.RuneVLine)

		// inner space
		for x := 1; x < boxWidth-1; x++ {
			base.SetContent(int(origin.GetX())+x, int(origin.GetY())+y, ' ')
		}
	}

	color := base.StyleIt(tcell.ColorWhite)
	for i, line := range wrappedLines {
		for j, r := range line {
			base.SetContentWithStyle(int(origin.GetX())+padding+j, int(origin.GetY())+padding+i, r, color)
		}
	}

	// bottom line
	for x := 0; x < boxWidth-2; x++ {
		base.SetContent(int(origin.GetX())+x+1, int(origin.GetY())+boxHeight-1, tcell.RuneHLine)
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

func DrawRectCenter(width, height int, fn func(x, y int)) {
	w, h := base.GetSize()
	ui.DrawRect(base.Point{X: w, Y: h}, width, height, func(x, y int) {
		fn(x, y)
	})
}

func DrawBoxedText(text string) {
	bottomPadding := 2
	lines := strings.Split(text, "\n")

	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	width := maxLen + 4
	height := len(lines) + 2 + bottomPadding

	DrawRectCenter(width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				base.SetContent(x+col, y+row, r)
			}
		}
	})
}

func SetStatus(text string) {
	mu.Lock()
	listOfStatus = append(listOfStatus, text) // safe add
	mu.Unlock()

	go func(msg string) {
		time.Sleep(3 * time.Second)

		mu.Lock()
		for i, v := range listOfStatus {
			// safe remove
			if v == msg {
				listOfStatus = append(listOfStatus[:i], listOfStatus[i+1:]...)
				break
			}
		}
		mu.Unlock()
	}(text)
}

func DrawRectStatus(text string, y int) {
	// Get terminal width
	w, _ := base.GetSize()
	color := base.StyleIt(tcell.ColorWhite)

	lines := strings.Split(text, "\n")
	jumpBy := 10

	// Find the longest line to determine rectangle width
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	width := maxLen + 4 // Padding around text
	height := len(lines) + 4

	// Calculate top-left corner of the rectangle
	start := base.Point{
		X: (w * 2) - width - 2,
		Y: 1 + jumpBy*(y+1),
	}

	// Draw the rectangle and render text inside
	ui.DrawRect(start, width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				base.SetContentWithStyle(x+col+2, y+row+2, r, color)
			}
		}
	})
}
