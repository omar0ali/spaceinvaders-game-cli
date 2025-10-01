package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

const (
	IncreaseGunCapBy   = 1
	IncreaseGunPowerBy = 1
	IncreaseGunSpeedBy = 2
)

type UI struct {
	MenuScreen         bool
	PauseScreen        bool
	GameOverScreen     bool
	LevelUpScreen      bool
	SpaceShipSelection bool
	timeElapsed        float64
	nextMinute         int
	status             string
	showStatus         bool
}

func NewUI(gc *core.GameContext) *UI {
	u := &UI{true, false, false, false, false, 0, 0, "Start New Game", false}
	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			u.SetStatus("Level Up +1")
			u.LevelUpScreen = true
		})
	}
	return u
}

func (u *UI) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	// start screen
	if u.MenuScreen {
		u.MessageBox(window.GetCenterPoint(),
			`
				The game is an endless space shooter where players face increasingly difficult 
				waves of alien ships that scale with their level.

				Each time the player levels up, they can choose an upgrade to improve their spaceship,
				such as boosting firepower to handle tougher aliens with stronger armor.

				The objective is to survive as long as possible, destroy alien ships, and push for 
				a higher score while managing health through occasional drop-down health packs that
				restore the spaceship health.

				[Controls]
				[LM] Or [Space] to shoot coming alien-ships.
				[F] Consume Health Kit.
				[P] To pause the game.
				[Q] To quit the game.

				Press [S] to start the game
			`,
			"Space Invaders Game")
	}

	if u.SpaceShipSelection {
		w, _ := window.GetSize()
		rectWidth := 45
		rectHeight := 10
		startPosY := 25
		startPosX := (w / 2) - rectWidth/2
		colGap := 54
		rowGap := 10
		columnsPerRow := 3

		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			for i, shape := range s.ListOfDesigns {
				rowIndex := i / columnsPerRow
				colIndex := i % columnsPerRow

				x := startPosX + colIndex*(rectWidth+colGap)
				y := startPosY + rowIndex*(rectHeight+rowGap)

				DrawRect(core.Point{X: x, Y: y}, rectWidth, rectHeight, func(initX int, initY int) {
					// draw index (spaceship selection key)
					for j, r := range fmt.Sprintf("[%d]", i+1) {
						window.SetContent(initX+j, initY, r)
					}

					// draw the spaceship shape inside the rectangle
					gap := 4
					for rowIndex, line := range shape.Shape {
						color := window.StyleIt(tcell.ColorReset, shape.GetColor())
						for colIndex, char := range line {
							if char != ' ' {
								window.SetContentWithStyle(initX+colIndex+gap, initY+rowIndex, char, color)
							}
						}
						// draw details of the spaceship
						str := []string{
							fmt.Sprintf("[%s]", shape.Name),
							fmt.Sprintf("* Gun Power: %d", shape.GunPower),
							fmt.Sprintf("* Gun Capacity: %d", shape.GunCap),
							fmt.Sprintf("* Gun Speed: %d", shape.GunSpeed),
						}
						for j, line := range str {
							for i, r := range line {
								window.SetContentWithStyle(initX+colIndex+(gap*5)+i, initY+j, r, color)
							}
						}
					}
				})
			}
		}
	}

	if u.LevelUpScreen {
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			u.MessageBox(window.GetCenterPoint(),
				fmt.Sprintf(`
				[Player] Current Level: %d

				(*) Choose a stat to upgrade:

				[A] (%d) Increase Gun Power by %d
				[S] (%d/%d) Increase Gun Speed by %d
				[D] (%d/%d) Increase Gun Capacity by %d
				[C] (%d/%d) Restore Full Health
				`, s.Level,
					s.Power,
					IncreaseGunPowerBy,
					s.Speed,
					s.cfg.SpaceShipConfig.GunMaxSpeed,
					IncreaseGunSpeedBy,
					s.Cap,
					s.cfg.SpaceShipConfig.GunMaxCap,
					IncreaseGunCapBy,
					s.health,
					s.SpaceshipDesign.EntityHealth),
				"Level Up")
		}
	}

	// show controls at the bottom of the screen
	_, h := window.GetSize()
	for i, r := range []rune("[LM] Shoot Beams ◆ [F] Consume Health Kit ◆ [P] Pause Game ◆ [R] Restart Game ◆ [Q] Quit") {
		window.SetContentWithStyle(i, h-1, r, whiteColor)
	}

	// timer
	minutes := int(u.timeElapsed) / 60
	seconds := int(u.timeElapsed) % 60

	w, _ := window.GetSize()
	timeStr := []rune(fmt.Sprintf("Time: %02d:%02d", minutes, seconds))
	// display objects details
	for i, r := range timeStr {
		window.SetContent((w-len(timeStr))+i, 0, r)
	}

	// display spacehsip details - Also drop a health kit every minute
	if spaceship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		if u.nextMinute < minutes {
			spaceship.HealthProducer.DeployHealthKit()
			u.nextMinute++
		}

		spaceship.UISpaceshipData(gc)
	}

	// display aliens details
	if aliens, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		aliens.UIAlienShipData(gc)
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
			`, spaceship.health, spaceship.Score, spaceship.Kills, spaceship.Level),
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
	if u.showStatus {
		DrawRectStatus(u.status)
	}
}

func (u *UI) Update(gc *core.GameContext, delta float64) {
	if u.MenuScreen || u.PauseScreen || u.GameOverScreen || u.LevelUpScreen || u.SpaceShipSelection {
		gc.Halt = true
	} else {
		gc.Halt = false
		u.timeElapsed += delta
	}
}

func (u *UI) InputEvents(events tcell.Event, gc *core.GameContext) {
	upgrade := func(up func() bool) {
		if up() {
			u.LevelUpScreen = false
		}
	}
	switch ev := events.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 's' || ev.Rune() == 'S' {
			if u.MenuScreen {
				u.SetStatus("Select a Spaceship")
				u.SpaceShipSelection = true
				u.MenuScreen = false
			}
		}
		if ev.Rune() == 'p' || ev.Rune() == 'P' {
			if u.MenuScreen || u.GameOverScreen { // skip
				return
			}
			u.PauseScreen = !u.PauseScreen
		}
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			if u.SpaceShipSelection {
				if ev.Key() == tcell.KeyRune {
					n := int(ev.Rune() - '0')
					switch n {
					case 1, 2, 3, 4, 5:
						s.SpaceshipSelection(n - 1)
						u.SpaceShipSelection = false
						u.SetStatus("Get Ready!")
					}
				}
			}
			if u.LevelUpScreen {
				if ev.Rune() == 'A' || ev.Rune() == 'a' {
					upgrade(func() bool {
						u.SetStatus(fmt.Sprintf("Gun Power: +%d", IncreaseGunPowerBy))
						return s.IncreaseGunPower(IncreaseGunPowerBy)
					})
				}
				if ev.Rune() == 'S' || ev.Rune() == 's' {
					upgrade(func() bool {
						u.SetStatus(fmt.Sprintf("Gun Speed: +%d", IncreaseGunSpeedBy))
						return s.IncreaseGunSpeed(IncreaseGunSpeedBy)
					})
				}
				if ev.Rune() == 'D' || ev.Rune() == 'd' {
					upgrade(func() bool {
						u.SetStatus(fmt.Sprintf("Gun Cap: +%d", IncreaseGunCapBy))
						return s.IncreaseGunCap(IncreaseGunCapBy)
					})
				}
				if ev.Rune() == 'C' || ev.Rune() == 'c' {
					upgrade(func() bool {
						u.SetStatus("Health Restored")
						return s.RestoreFullHealth()
					})
				}
			}
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

func DrawRect(pos core.Point, width, height int, fn func(initX, initY int)) {
	const padding = 2
	centerOfW := pos.X / 2
	centerOfH := pos.Y / 2
	startX := centerOfW - (width / 2)
	startY := centerOfH - (height / 2)
	for i := range height {
		for j := range width {
			switch {
			case j == 0 && i == 0:
				window.SetContent(startX+j, startY+i, tcell.RuneULCorner)
			case j == width-1 && i == 0:
				window.SetContent(startX+j, startY+i, tcell.RuneURCorner)
			case j == 0 && i == height-1:
				window.SetContent(startX+j, startY+i, tcell.RuneLLCorner)
			case j == width-1 && i == height-1:
				window.SetContent(startX+j, startY+i, tcell.RuneLRCorner)

			case i == 0 || i == height-1:
				window.SetContent(startX+j, startY+i, tcell.RuneHLine)
			case j == 0 || j == width-1:
				window.SetContent(startX+j, startY+i, tcell.RuneVLine)

			default:
				window.SetContent(startX+j, startY+i, ' ')
			}
		}
	}
	fn(startX+padding, startY+padding)
}

func DrawRectCenter(width, height int, fn func(x, y int)) {
	w, h := window.GetSize()
	DrawRect(core.Point{X: w, Y: h}, width, height, func(x, y int) {
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
				window.SetContent(x+col, y+row, r)
			}
		}
	})
}

func (u *UI) SetStatus(text string) {
	u.showStatus = true
	u.status = text
	go func() {
		time.Sleep(1 * time.Second)
		u.showStatus = false
	}()
}

func DrawRectStatus(text string) {
	w, _ := window.GetSize()
	color := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
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
	DrawRect(core.Point{X: (w * 2) - width - 1, Y: 15}, width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				window.SetContentWithStyle(x+col, y+row, r, color)
			}
		}
	})
}
