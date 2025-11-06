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
	exitCha            chan struct{}
	cfg                game.GameConfig
}

func NewUI(gc *game.GameContext, cfg game.GameConfig, exitCha chan struct{}) *UI {
	nextMinute = 0

	u := &UI{
		MenuScreen:         true,
		PauseScreen:        false,
		GameOverScreen:     false,
		LevelUpScreen:      false,
		SpaceShipSelection: false,
		exitCha:            exitCha,
		cfg:                cfg,
	}

	if u.MenuScreen {
		if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {
			boxes := []*ui.Box{
				ui.NewUIBox(
					[]string{
						"Start New Game",
					}, ui.StartGameDesc,
					func() {
						// here we should start the game
						SetStatus("Select a Spaceship", gc)
						u.SpaceShipSelection = true
						if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
							var boxes []*ui.Box
							for i, shipDesign := range s.LoadedDesigns.ListOfSpaceships {
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
										SetStatus(fmt.Sprintf("%s Selected", name), gc)
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
				ui.NewUIBox(
					[]string{
						"Compendium",
					},
					[]string{
						"Scan the battlefield: ships, asteroids and abilities.",
					}, func() {
						// init items for the menu
						abilitiesItems := make([]*ui.Box, 0)
						spaceshipsItems := make([]*ui.Box, 0)
						asteroidsItems := make([]*ui.Box, 0)
						alienShipsItems := make([]*ui.Box, 0)
						bossShipsItems := make([]*ui.Box, 0)
						modifiersItems := make([]*ui.Box, 0)

						// Load designs for each items
						if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
							for _, i := range ship.LoadedDesigns.ListOfAbilities {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("* Description:    %s", i.Description),
									fmt.Sprintf("* Status:    %s", i.Status),
								}

								abilitiesItems = append(
									abilitiesItems,
									ui.NewUIBox(i.Shape, descriptions, nil), // using hover
								)
							}
							for _, i := range ship.LoadedDesigns.ListOfSpaceships {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("* HP:         %d", i.EntityHealth),
									fmt.Sprintf("* Gun POW:    %d", i.GunPower),
									fmt.Sprintf("* Gun CAP:    %d", i.GunCap),
									fmt.Sprintf("* Gun SPD:    %d", i.GunSpeed),
									fmt.Sprintf("* Gun CD:     %d ms", i.GunCooldown),
									fmt.Sprintf("* Gun RLD CD: %d ms", i.GunReloadCooldown),
								}
								spaceshipsItems = append(
									spaceshipsItems,
									ui.NewUIBox(i.Shape, descriptions, nil), // using hover
								)
							}
							for _, i := range ship.LoadedDesigns.ListOfAsteroids.Asteroids {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("Color:  %s", i.Color),
									fmt.Sprintf("Health: %d", i.EntityHealth),
								}

								asteroidsItems = append(asteroidsItems,
									ui.NewUIBox(i.Shape, descriptions, nil))
							}
							for _, i := range ship.LoadedDesigns.ListOfAlienships {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("* HP:         %d", i.EntityHealth),
									fmt.Sprintf("* Gun POW:    %d", i.GunPower),
									fmt.Sprintf("* Gun CAP:    %d", i.GunCap),
									fmt.Sprintf("* Gun SPD:    %d", i.GunSpeed),
									fmt.Sprintf("* Gun CD:     %d ms", i.GunCooldown),
									fmt.Sprintf("* Gun RLD CD: %d ms", i.GunReloadCooldown),
								}
								alienShipsItems = append(alienShipsItems,
									ui.NewUIBox(i.Shape, descriptions, nil))
							}
							for _, i := range ship.LoadedDesigns.ListOfBossShips {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("* HP:         %d", i.EntityHealth),
									fmt.Sprintf("* Gun POW:    %d", i.GunPower),
									fmt.Sprintf("* Gun CAP:    %d", i.GunCap),
									fmt.Sprintf("* Gun SPD:    %d", i.GunSpeed),
									fmt.Sprintf("* Gun CD:     %d ms", i.GunCooldown),
									fmt.Sprintf("* Gun RLD CD: %d ms", i.GunReloadCooldown),
								}
								bossShipsItems = append(bossShipsItems,
									ui.NewUIBox(i.Shape, descriptions, nil))
							}
							for _, i := range ship.LoadedDesigns.ModifierDesign {
								descriptions := []string{
									fmt.Sprintf("- [%s]", i.Name),
									fmt.Sprintf("* Health:     %d", i.EntityHealth),
									fmt.Sprintf("* Modify Gun POW:     %d", i.ModifyGunPower),
									fmt.Sprintf("* Modify Gun CAP:     %d", i.ModifyGunCap),
									fmt.Sprintf("* Modify Gun SPD:     %d", i.ModifyGunSpeed),
									fmt.Sprintf("* Modify Gun CD:      %d", i.ModifyGunCoolDown),
									fmt.Sprintf("* Modify Gun CD RLD:  %d", i.ModifyGunReloadCoolDown),
									fmt.Sprintf("* Max:     %d", i.MaxValue),
								}

								modifiersItems = append(modifiersItems,
									ui.NewUIBox(i.Shape, descriptions, nil))
							}

						}
						layoutCodexMenu := ui.InitCodexMenu(20, 5)
						boxes := make([]*ui.Box, 0)
						boxes = append(boxes,
							ui.NewUIBox(
								[]string{
									"Abilities",
								},
								[]string{
									"Displaying the Abilities",
								}, func() {
									layoutCodexMenu.SetList(abilitiesItems)
								}),
							ui.NewUIBox(
								[]string{
									"Spaceships",
								},
								[]string{
									"Displaying the Spaceships",
								}, func() {
									layoutCodexMenu.SetList(spaceshipsItems)
								}),
							ui.NewUIBox(
								[]string{
									"Asteroids",
								},
								[]string{
									"Displaying the Asteroids",
								}, func() {
									layoutCodexMenu.SetList(asteroidsItems)
								}),
							ui.NewUIBox(
								[]string{
									"Alienships",
								},
								[]string{
									"Displaying the Alienships",
								}, func() {
									layoutCodexMenu.SetList(alienShipsItems)
								}),
							ui.NewUIBox(
								[]string{
									"Boss Spaceships",
								},
								[]string{
									"Displaying the Boss Spaceships",
								}, func() {
									layoutCodexMenu.SetList(bossShipsItems)
								}),
							ui.NewUIBox(
								[]string{
									"Modifiers",
								},
								[]string{
									"Displaying the Modifiers",
								}, func() {
									layoutCodexMenu.SetList(modifiersItems)
								}),
							ui.NewUIBox(
								[]string{
									"< Back",
								},
								[]string{
									"Back to main menu.",
								}, func() {
									RestartGame(gc, u.cfg, u.exitCha)
								}))
						layoutCodexMenu.SetMenuItems(boxes)
						layout.SetLayout(layoutCodexMenu)
					},
				),
				ui.NewUIBox([]string{
					"Quit Game",
				}, []string{"Quit the game."}, func() {
					base.ExitGame(exitCha)
				}),
			}
			layout.SetLayout(
				ui.InitMainMenu(20, 5, boxes...),
			)
		}
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

	// timer
	minutes = int(u.timeElapsed) / 60
	seconds = int(u.timeElapsed) % 60

	if !u.MenuScreen {
		w, h := base.GetSize()
		// draw a line

		for i := range w {
			base.SetContent(i, h-2, tcell.RuneHLine)
		}

		// show controls at the bottom of the screen
		controlsUI := []rune("[LM] Shoot Beams ◆ [E] Consume Health Kit ◆ [R] Reload Gun ◆ [P] Pause Game ◆ [Ctrl+R] Restart Game ◆ [Ctrl+Q] Quit")
		for i, r := range controlsUI {
			base.SetContentWithStyle(w/2-(len(controlsUI)/2)+i, h-1, r, whiteColor)
		}

		whiteColor := base.StyleIt(tcell.ColorWhite)
		greenColor := base.StyleIt(tcell.ColorGreenYellow)

		// top left box
		ui.DrawBoxOverlap(base.Point{X: 0, Y: 0}, 35, 5, func(x int, y int) {
			// display time details
			timeStr := []rune(fmt.Sprintf("Time: %02d:%02d", minutes, seconds))
			for i, r := range timeStr {
				base.SetContentWithStyle(i+x+2, y+1, r, whiteColor)
			}

			if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
				// display score
				txtScore := "Score: "
				barSize := 22
				for i, r := range txtScore {
					base.SetContentWithStyle(i+x+2, y+2, r, whiteColor)
				}

				base.DisplayBar(
					&s.Score,
					base.WithPosition(x+len(txtScore)+2, y+2),
					base.WithBarSize(barSize),
					base.WithStatus(false),
					base.WithStyle(whiteColor),
				)

				for i, r := range []rune(fmt.Sprintf("Kills: %d", s.Kills)) {
					base.SetContentWithStyle(i+x+2, y+3, r, whiteColor)
				}
			}
		}, greenColor)

		// display spacehsip details - Also drop a health kit every minute
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			s.UISpaceshipData(gc)
		}

		// display aliens details
		if aliens, ok := gc.FindEntity("alien").(*AlienProducer); ok {
			aliens.UIAlienShipData(gc)
		}

	}

	// game over ui
	if u.GameOverScreen {
		if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
			u.MessageBox(base.GetCenterPoint(),
				fmt.Sprintf(`
			Taken damage from:
			%v

			Killed By:
			%s Level: %d

			Thank you for playing :)
			---------------------------------------
			Would you like to play again?
			[Ctrl+R] To Restart.
			[Ctrl+Q] To Quit.
			`, strings.Join(s.GetRegisteredHits(), "\n"), s.KilledBy.Name, s.KilledBy.Power),
				"Game Over",
			)
		}
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
		if ev.Rune() == 'p' || ev.Rune() == 'P' || ev.Key() == tcell.KeyESC {
			if u.MenuScreen || u.GameOverScreen || u.SpaceShipSelection || u.LevelUpScreen { // skip
				return
			}
			gc.Sounds.PlaySound("8-bit-game-sfx-sound-select.mp3", -1)
			u.PauseScreen = !u.PauseScreen
			if layout, ok := gc.FindEntity("layout").(*ui.UISystem); ok {
				var spaceship *SpaceShip
				if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
					spaceship = ship
				}
				if u.PauseScreen {
					boxes := []*ui.Box{
						ui.NewUIBox(
							[]string{
								"Continue",
							},
							[]string{
								"Continue the game.",
							}, func() {
								u.PauseScreen = false
								layout.SetLayout(nil)
							},
						),
						ui.NewUIBox(
							[]string{
								"My Spaceship",
							}, []string{
								fmt.Sprintf("[%s] - [Level: %d]", spaceship.SelectedSpaceship.Name, spaceship.Level),
								"---------------------------------",
								fmt.Sprintf("Gun Capacity:         %d +(%d) -> %d",
									spaceship.SelectedSpaceship.GunCap,
									spaceship.GetCapacity()-spaceship.SelectedSpaceship.GunCap,
									spaceship.GetCapacity(),
								),
								fmt.Sprintf("Gun Speed:            %d +(%d) -> %d",
									spaceship.SelectedSpaceship.GunSpeed,
									spaceship.GetSpeed()-spaceship.SelectedSpaceship.GunSpeed,
									spaceship.GetSpeed(),
								),
								fmt.Sprintf("Gun Power:            %d +(%d) -> %d",
									spaceship.SelectedSpaceship.GunPower,
									spaceship.GetPower()-spaceship.SelectedSpaceship.GunPower,
									spaceship.GetPower(),
								),
								fmt.Sprintf("Gun Cooldown:         %d +(%d) -> %d",
									spaceship.SelectedSpaceship.GunCooldown,
									int(spaceship.GetCooldown())-spaceship.SelectedSpaceship.GunCooldown,
									spaceship.GetCooldown(),
								),
								fmt.Sprintf("Gun Reload Cooldown:  %d +(%d) -> %d",
									spaceship.SelectedSpaceship.GunReloadCooldown,
									int(spaceship.GetReloadCooldown())-spaceship.SelectedSpaceship.GunReloadCooldown,
									spaceship.GetReloadCooldown(),
								),
								fmt.Sprintf("Spaceship Health:     %d +(%d) -> %d",
									spaceship.SelectedSpaceship.EntityHealth,
									spaceship.MaxHealth-spaceship.SelectedSpaceship.EntityHealth,
									spaceship.MaxHealth,
								),
							}, func() {

							},
						),
						ui.NewUIBox(
							[]string{
								"Restart",
							},
							[]string{
								"Return to Main Menu.",
							},
							func() {
								RestartGame(gc, u.cfg, u.exitCha)
							},
						),
						ui.NewUIBox([]string{
							"Quit Game",
						}, []string{"Exit the game."}, func() {
							base.ExitGame(u.exitCha)
						}),
					}
					menuUi := ui.InitMainMenu(20, 5, boxes...)
					menuUi.SelectedDesc = []string{"Paused Game"}
					layout.SetLayout(menuUi)
				} else {
					layout.SetLayout(nil)
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

func SetStatus(text string, gc *game.GameContext) {
	mu.Lock()
	gc.Sounds.PlaySound("8-bit-game-sfx-notification.mp3", 0)
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
	jumpBy := 6

	// Find the longest line to determine rectangle width
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	width := maxLen + 4 // Padding around text
	height := len(lines) + 2

	// Calculate top-left corner of the rectangle
	start := base.Point{
		X: (w * 2) - width,
		Y: 2 + jumpBy*y,
	}

	// Draw the rectangle and render text inside
	ui.DrawRect(start, width, height, func(x, y int) {
		for row, line := range lines {
			for col, r := range line {
				base.SetContentWithStyle(x+col+2, y+row+1, r, color)
			}
		}
	})
}
