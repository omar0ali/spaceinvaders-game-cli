package entities

import (
	"fmt"
	"strings"
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
			if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {

				upgrade := func(up func() bool) {
					if up() {
						u.LevelUpScreen = false
					}
					layout.SetLayout(nil)
				}

				SetStatus("Level Up")
				u.LevelUpScreen = true
				boxes := []*ui.Box{
					ui.NewUIBox(
						[]string{
							"     /\\      ",
							"    /***\\    ",
							"   | *+* |   ",
							"    \\***/    ",
							"     ||      ",
							"    POWER    ",
						},
						[]string{
							fmt.Sprintf("(%d) Increase Gun Power by %d", s.GetPower(), IncreaseGunPowerBy),
						},
						func() {
							upgrade(func() bool {
								SetStatus(fmt.Sprintf("Gun Power: +%d", IncreaseGunPowerBy))
								return s.IncreaseGunPower(IncreaseGunPowerBy)
							})
						},
					),

					ui.NewUIBox(
						[]string{
							"     ___     ",
							"    />>>\\    ",
							"   |>>+>>|   ",
							"    \\>>>/    ",
							"     ~~~     ",
							"    SPEED    ",
						},
						[]string{
							fmt.Sprintf("(%d/%d) Increase Gun Speed by %d", int(s.GetSpeed()), s.cfg.SpaceShipConfig.GunMaxSpeed, IncreaseGunSpeedBy),
						},
						func() {
							upgrade(func() bool {
								SetStatus(fmt.Sprintf("Gun Speed: +%d", IncreaseGunSpeedBy))
								return s.IncreaseGunSpeed(IncreaseGunSpeedBy, s.cfg.SpaceShipConfig.GunMaxSpeed)
							})
						},
					),

					ui.NewUIBox(
						[]string{
							"   _______   ",
							"  |[1] [2]|  ",
							"  |[3] [+]|  ",
							"  |_______|  ",
							"     CAP     ",
						},
						[]string{
							fmt.Sprintf("(%d/%d) Increase Gun Capacity by %d", s.GetCapacity(), s.cfg.SpaceShipConfig.GunMaxCap, IncreaseGunCapBy),
						},
						func() {
							upgrade(func() bool {
								SetStatus(fmt.Sprintf("Gun Capcity: +%d", IncreaseGunCapBy))
								return s.IncreaseGunCap(IncreaseGunCapBy, s.cfg.SpaceShipConfig.GunMaxCap)
							})
						},
					),

					ui.NewUIBox(
						[]string{
							"    _____    ",
							"   | -=- |   ",
							"   | -3- |   ",
							"   |_____|   ",
							"    \\___/    ",
							"  COOL DOWN  ",
						},
						[]string{
							fmt.Sprintf("(%d) Decrease Gun Cooldown by %d", s.GetCooldown(), DecreaseGunCooldownBy),
						},
						func() {
							upgrade(func() bool {
								if s.DecreaseCooldown(DecreaseGunCooldownBy) {
									SetStatus(fmt.Sprintf("Gun Cooldown: -%d", DecreaseGunCooldownBy))
									return true
								}
								SetStatus("Gun Cooldown: Maxed Out!")
								return false
							})
						},
					),

					ui.NewUIBox(
						[]string{
							"    _____    ",
							"   |\\___/|   ",
							"   || - ||   ",
							"   || 3 ||   ",
							"   ||___||   ",
							"   RLD CLD   ",
						},
						[]string{
							fmt.Sprintf("(%d) Decrease Gun Reload Cooldown by %d", s.GetReloadCooldown(), DecreaseGunCooldownBy),
						},
						func() {
							upgrade(func() bool {
								if s.DecreaseGunReloadCooldown(DecreaseGunCooldownBy) {
									SetStatus(fmt.Sprintf("Gun Reload Cooldown: -%d", DecreaseGunCooldownBy))
									return true
								}
								SetStatus("Gun Reload Cooldown: Maxed Out!")
								return false
							})
						},
					),

					ui.NewUIBox(
						[]string{
							"    _____    ",
							"   / +++ \\   ",
							"  | +100+ |  ",
							"   \\ ___ /   ",
							"    |___|    ",
							"     HP%     ",
						},
						[]string{
							fmt.Sprintf("(%d/%d) Restore Full Health", s.Health, s.SelectedSpaceship.EntityHealth),
						},
						func() {
							upgrade(func() bool {
								SetStatus("[H] Spaceship health has been restored!")
								return s.RestoreFullHealth()
							})
						},
					),
				}
				layout.SetLayout(
					ui.InitLayout(21, 10, boxes...),
				)
			}
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

				[LM] hold to shoot a beam to coming alien-ships.

				[E] Consume Health Kit.
				[P] To pause the game.

				[Ctrl+R] to restart the game.
				[Ctrl+Q] To quit the game.

				Press [S] to start the game
			`,
			"Space Invaders Game")
	}

	// show controls at the bottom of the screen
	_, h := base.GetSize()
	for i, r := range []rune("[LM] Shoot Beams ◆ [E] Consume Health Kit ◆ [R] Reload Gun ◆ [P] Pause Game ◆ [Ctrl+R] Restart Game ◆ [Ctrl+Q] Quit") {
		base.SetContentWithStyle(i, h-1, r, whiteColor)
	}

	// timer
	minutes = int(u.timeElapsed) / 60
	seconds = int(u.timeElapsed) % 60

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
	switch ev := events.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 's' || ev.Rune() == 'S' {
			if u.MenuScreen {
				SetStatus("Select a Spaceship")
				u.SpaceShipSelection = true
				if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
					var boxes []*ui.Box
					if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {
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
				}

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
	showStatus = true
	status = text
	go func() {
		time.Sleep(4 * time.Second)
		showStatus = false
	}()
}

func DrawRectStatus(text string) {
	w, _ := base.GetSize()
	color := base.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	lines := strings.Split(text, "\n")

	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	width := maxLen + 4
	height := len(lines) + 4
	ui.DrawRect(base.Point{X: (w * 2) - width - 6, Y: 15}, width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				base.SetContentWithStyle(x+col+2, y+row+2, r, color)
			}
		}
	})
}
