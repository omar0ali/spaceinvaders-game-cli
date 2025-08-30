// Package main
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/entities"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

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
	spaceship := entities.InitSpaceShip()

	alienProducer := entities.AlienProducer{
		Aliens: []*entities.Alien{},
	}

	starProducer := entities.StarProducer{
		Stars: []*entities.Star{},
	}

	// order is important since some objects might overlap others
	gameContext.AddEntity(&starProducer)
	gameContext.AddEntity(&alienProducer)
	gameContext.AddEntity(&spaceship)

	// ----------------------------------------- window ------------------------------------
	window.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'r' {
					close(exit)
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
			for _, entity := range gameContext.GetEntities() {
				entity.Draw(&gameContext)
				entity.Update(&gameContext, delta)
			}
		},
	)

	// exit
	<-exit
}
