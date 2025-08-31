// Package main
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/entities"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

func DeployEntities(gc *core.GameContext) {
	// order is important since some objects might overlap others
	gc.AddEntity(
		entities.InitSpaceShip(5, 5),
		&entities.AlienProducer{
			Aliens: []*entities.Alien{},
		},
		&entities.StarProducer{
			Stars: []*entities.Star{},
		},
		&entities.UI{
			MenuScreen:  true,
			PauseScreen: false,
		},
	)
}

func main() {
	exit := make(chan struct{})

	// window by default is set to 30 FPS
	// window.InitScreen(window.ChangeTickerDuration(16), window.EnableMouse) // this can update the framerate to 60
	// ------------------------------- Setup ------------------------------------
	screen := window.InitScreen(window.EnableMouse)
	screen.SetTitle("Space Invader Game")
	// ------------------------------------- Objects ----------------------------------
	gameContext := core.GameContext{
		Screen: screen,
	}
	// ---------------------------------- entities --------------------------------------

	DeployEntities(&gameContext)

	// ----------------------------------------- window ------------------------------------
	window.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'r' || ev.Rune() == 'R' {
					gameContext.RemoveAllEntities()
					DeployEntities(&gameContext)
				}
			}
			for _, entity := range gameContext.GetEntities() {
				entity.InputEvents(event, &gameContext)
			}
		},
	)

	window.Update(exit,
		func(delta float64) {
			// update game

			// only let ui to be displayed
			if gameContext.Halt {
				if ui, ok := gameContext.FindEntity("ui").(*entities.UI); ok {
					ui.Update(&gameContext, delta)
					ui.Draw(&gameContext)
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
