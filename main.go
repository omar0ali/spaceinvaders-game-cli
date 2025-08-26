// Package main
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/entities"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

func main() {
	exit := make(chan int)

	// window by default is set to 30 FPS
	// window.InitScreen(window.ChangeTickerDuration(16)) // this can update the framerate to 60
	screen := window.InitScreen()
	w, h := window.GetSize()
	window.SetTitle("Space Invader Game")
	// ------------------------------------- Objects ----------------------------------
	gameContext := core.GameContext{
		Screen: screen,
	}
	// ---------------------------------- entities --------------------------------------
	spaceship := entities.InitSpaceShip(core.Point{
		X: w / 2,
		Y: h - 4,
	})

	gameContext.AddEntity(spaceship)

	window.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'r' {
					// testing exit with letter r
					exit <- 0
				}
				for _, entity := range gameContext.GetEntities() {
					entity.InputEvents(event, &gameContext)
				}
			}
		},
	)

	window.Update(exit,
		func(delta float64) {
			// update game
			for _, entity := range gameContext.GetEntities() {
				entity.Draw(&gameContext)
				entity.Update(&gameContext)
			}
		},
	)

	// exit
	if val := <-exit; val == 0 {
		return
	}
}
