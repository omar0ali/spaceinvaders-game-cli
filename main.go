// Package main
package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/ui"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

func main() {
	cfg := game.LoadConfig()

	// setup logs
	if cfg.Dev.Debug {
		logFile := game.SetupLogs()
		defer logFile.Close()
	}

	exit := make(chan struct{})

	// ------------------------------- Setup ------------------------------------
	screen := base.InitScreen(base.EnableMouse)
	screen.SetTitle("Space Invader Game")

	sounds := game.InitSoundSystem(cfg)

	// ------------------------------------- Objects ----------------------------------
	gameContext := game.GameContext{
		Screen: screen,
		Sounds: sounds,
	}
	// ---------------------------------- entities --------------------------------------

	log.Println("Game running...")

	entities.StartGame(&gameContext, cfg, exit)

	// ----------------------------------------- window ------------------------------------
	base.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlR {
					entities.RestartGame(&gameContext, cfg, exit)
				}
			}
			for _, entity := range gameContext.GetEntities() {
				entity.InputEvents(event, &gameContext)
			}
		},
	)

	base.Update(exit,
		func(delta float64) {
			// update game
			if cfg.Dev.FPSCounter {
				// fps
				for i, r := range []rune(fmt.Sprintf("FPS: %.2f", (1 / delta))) {
					base.SetContent(i, 0, r)
				}
			}
			// only let ui to be displayed
			if gameContext.Halt {
				if star, ok := gameContext.FindEntity("star").(*entities.StarProducer); ok {
					star.Update(&gameContext, delta)
					star.Draw(&gameContext)
				}
				if ui, ok := gameContext.FindEntity("ui").(*entities.UI); ok {
					ui.Update(&gameContext, delta)
					ui.Draw(&gameContext)
				}
				if layout, ok := gameContext.FindEntity("layout").(*ui.UISystem); ok {
					layout.Update(&gameContext, delta)
					layout.Draw(&gameContext)
				}
				// draw only (to display all objects when game is paused)
				for _, entity := range gameContext.GetEntities() {
					entity.Draw(&gameContext)
				}
			} else { // update everything
				for _, entity := range gameContext.GetEntities() {
					entity.Draw(&gameContext)
					entity.Update(&gameContext, delta)
				}
			}
		},
	)

	// exit
	<-exit
}
