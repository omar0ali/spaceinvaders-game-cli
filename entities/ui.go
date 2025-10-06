package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

const (
	IncreaseGunCapBy      = 1
	IncreaseGunPowerBy    = 1
	IncreaseGunSpeedBy    = 1
	DecreaseGunCooldownBy = 3
)

var (
	nextMinute int
	status     string
	showStatus bool
	minutes    int
	seconds    int
)

type UI struct {
	MenuScreen         bool
	PauseScreen        bool
	GameOverScreen     bool
	LevelUpScreen      bool
	SpaceShipSelection bool
	timeElapsed        float64
}

func NewUI(gc *game.GameContext) *UI {
	nextMinute = 0
	u := &UI{true, false, false, false, false, 0}
	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			SetStatus("Level Up")
			u.LevelUpScreen = true
		})
	}
	return u
}

func (u *UI) Draw(gc *game.GameContext) {
	whiteColor := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)

	// start screen
	if u.MenuScreen {
		u.MessageBox(base.GetCenterPoint(),
			`
				The game is an endless space shooter where players face increasingly difficult 
				waves of alien ships that scale with their level.

				Each time the player levels up, they can choose an upgrade to improve their spaceship,
				such as boosting firepower to handle tougher aliens with stronger armor.

				The objective is to survive as long as possible, destroy alien ships, and push for 
				a higher score while managing health through occasional drop-down health packs that
				restore the spaceship health.

				 -----------------
				(*) Controls
				 -----------------

				[LM] Or [Space] to shoot a beam to coming alien-ships.

				[E] Consume Health Kit.
				[P] To pause the game.

				[Ctrl+R] to restart the game.
				[Ctrl+Q] To quit the game.

				Press [S] to start the game
			`,
			"Space Invaders Game")
	}

	if u.SpaceShipSelection {
		w, _ := base.GetSize()
		rectWidth := 45
		rectHeight := 10
		startPosY := 25
		startPosX := (w / 2) - rectWidth/2
		colGap := 54
		rowGap := 10
		columnsPerRow := 3

		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			for i, shape := range s.ListOfSpaceships {
				rowIndex := i / columnsPerRow
				colIndex := i % columnsPerRow

				x := startPosX + colIndex*(rectWidth+colGap)
				y := startPosY + rowIndex*(rectHeight+rowGap)

				DrawRect(base.Point{X: x, Y: y}, rectWidth, rectHeight, func(initX int, initY int) {
					// draw index (spaceship selection key)
					for j, r := range fmt.Sprintf("[%d]", i+1) {
						base.SetContent(initX+j, initY, r)
					}

					// draw the spaceship shape inside the rectangle
					gap := 4
					for rowIndex, line := range shape.Shape {
						color := base.StyleIt(tcell.ColorReset, shape.GetColor())
						for colIndex, char := range line {
							if char != ' ' {
								base.SetContentWithStyle(initX+colIndex+gap, initY+rowIndex, char, color)
							}
						}
						// draw details of the spaceship
						str := []string{
							fmt.Sprintf("[%s]", shape.Name),
							fmt.Sprintf("* HP:         %d", shape.EntityHealth),
							fmt.Sprintf("* Gun PWD:    %d", shape.GunPower),
							fmt.Sprintf("* Gun CAP:    %d", shape.GunCap),
							fmt.Sprintf("* Gun SPD:    %d", shape.GunSpeed),
							fmt.Sprintf("* Gun CD:     %d ms", shape.GunCooldown),
							fmt.Sprintf("* Gun RLD CD: %d ms", shape.GunReloadCooldown),
						}

						for j, line := range str {
							for i, r := range line {
								base.SetContentWithStyle(initX+colIndex+(gap*4)+i, initY+j, r, color)
							}
						}
					}
				})
			}
		}
	}

	if u.LevelUpScreen {
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			u.MessageBox(base.GetCenterPoint(),
				fmt.Sprintf(`
				(*) Level: %d

				(*) Choose a stat to upgrade:                  

				[A] (%d) Increase Gun Power by %d
				[S] (%d/%d) Increase Gun Speed by %d
				[D] (%d/%d) Increase Gun Capacity by %d
				[F] (%d) Decrease Gun Cooldown by %d
				[G] (%d) Decrease Gun Reload Cooldown by %d
				[H] (%d/%d) Restore Full Health
				`, s.Level,
					s.GetPower(),
					IncreaseGunPowerBy,
					int(s.Gun.GetSpeed()),
					s.cfg.SpaceShipConfig.GunMaxSpeed,
					IncreaseGunSpeedBy,
					s.GetCapacity(),
					s.cfg.SpaceShipConfig.GunMaxCap,
					IncreaseGunCapBy,
					s.GetCooldown(),
					DecreaseGunCooldownBy,
					s.GetReloadCooldown(),
					DecreaseGunCooldownBy,
					s.Health,
					s.SelectedSpaceship.EntityHealth),
				"Level Up")
		}
	}

	// show controls at the bottom of the screen
	_, h := base.GetSize()
	for i, r := range []rune("[LM] or [Space] Shoot Beams ◆ [E] Consume Health Kit ◆ [R] Reload Gun ◆ [P] Pause Game ◆ [Ctrl+R] Restart Game ◆ [Ctrl+Q] Quit") {
		base.SetContentWithStyle(i, h-1, r, whiteColor)
	}

	// timer
	minutes = int(u.timeElapsed) / 60
	seconds = int(u.timeElapsed) % 60

	w, _ := base.GetSize()
	timeStr := []rune(fmt.Sprintf("Time: %02d:%02d", minutes, seconds))
	// display objects details
	for i, r := range timeStr {
		base.SetContent((w-len(timeStr))+i, 0, r)
	}

	// display spacehsip details - Also drop a health kit every minute
	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.UISpaceshipData(gc)
	}

	// display aliens details
	if aliens, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		aliens.UIAlienShipData(gc)
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

	if showStatus {
		DrawRectStatus(status)
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
	upgrade := func(up func() bool) {
		if up() {
			u.LevelUpScreen = false
		}
	}
	switch ev := events.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 's' || ev.Rune() == 'S' {
			if u.MenuScreen {
				SetStatus("Select a Spaceship (1 - 5)")
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
						name := s.SpaceshipSelection(n - 1)
						SetStatus(fmt.Sprintf("[%d] %s Selected", n, name))
						u.SpaceShipSelection = false
					}
				}
			}
			if u.LevelUpScreen {
				if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
					if ev.Rune() == 'A' || ev.Rune() == 'a' {
						upgrade(func() bool {
							SetStatus(fmt.Sprintf("[A] Gun Power: +%d", IncreaseGunPowerBy))
							return s.IncreaseGunPower(IncreaseGunPowerBy)
						})
					}
					if ev.Rune() == 'S' || ev.Rune() == 's' {
						upgrade(func() bool {
							SetStatus(fmt.Sprintf("[S] Gun Speed: +%d", IncreaseGunSpeedBy))
							return s.IncreaseGunSpeed(IncreaseGunSpeedBy, s.cfg.SpaceShipConfig.GunMaxSpeed)
						})
					}
					if ev.Rune() == 'D' || ev.Rune() == 'd' {
						upgrade(func() bool {
							SetStatus(fmt.Sprintf("[D] Gun Capcity: +%d", IncreaseGunCapBy))
							return s.IncreaseGunCap(IncreaseGunCapBy, s.cfg.SpaceShipConfig.GunMaxCap)
						})
					}
					if ev.Rune() == 'F' || ev.Rune() == 'f' {
						upgrade(func() bool {
							if s.DecreaseCooldown(DecreaseGunCooldownBy) {
								SetStatus(fmt.Sprintf("[C] Gun Cooldown: -%d", DecreaseGunCooldownBy))
								return true
							}
							SetStatus("[C] Gun Cooldown: Maxed Out!")
							return false
						})
					}

					if ev.Rune() == 'G' || ev.Rune() == 'g' {
						upgrade(func() bool {
							if s.DecreaseGunReloadCooldown(DecreaseGunCooldownBy) {
								SetStatus(fmt.Sprintf("[G] Gun Reload Cooldown: -%d", DecreaseGunCooldownBy))
								return true
							}
							SetStatus("[G] Gun Reload Cooldown: Maxed Out!")
							return false
						})
					}

					if ev.Rune() == 'H' || ev.Rune() == 'h' {
						upgrade(func() bool {
							SetStatus("[H] Spaceship health has been restored!")
							return s.RestoreFullHealth()
						})
					}

				}
			}
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

	color := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)
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
	fn(startX+padding, startY+padding)
}

func DrawRectCenter(width, height int, fn func(x, y int)) {
	w, h := base.GetSize()
	DrawRect(base.Point{X: w, Y: h}, width, height, func(x, y int) {
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
	showStatus = true
	status = text
	go func() {
		time.Sleep(3 * time.Second)
		showStatus = false
	}()
}

func DrawRectStatus(text string) {
	w, _ := base.GetSize()
	color := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)
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
	DrawRect(base.Point{X: (w * 2) - width - 6, Y: 15}, width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				base.SetContentWithStyle(x+col, y+row, r, color)
			}
		}
	})
}
