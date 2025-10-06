// Package main
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/entities"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

func DeployEntities(gc *game.GameContext, cfg game.GameConfig) {
	// order is important since some objects might overlap others
	gc.AddEntity(entities.NewSpaceShip(cfg, gc))
	gc.AddEntity(entities.NewAlienProducer(gc))
	gc.AddEntity(entities.NewBossAlienProducer(cfg, gc))
	gc.AddEntity(&entities.Producer{})
	gc.AddEntity(entities.NewStarsProducer(cfg))
	gc.AddEntity(entities.NewUI(gc))
}

func main() {
	exit := make(chan struct{})

	// ------------------------------- Setup ------------------------------------
	screen := window.InitScreen(window.EnableMouse)
	screen.SetTitle("Space Invader Game")
	cfg := game.LoadConfig()

	// ------------------------------------- Objects ----------------------------------
	gameContext := game.GameContext{
		Screen: screen,
	}
	// ---------------------------------- entities --------------------------------------

	DeployEntities(&gameContext, cfg)

	// ----------------------------------------- window ------------------------------------
	window.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlR {
					gameContext.RemoveAllEntities()
					DeployEntities(&gameContext, cfg)
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
				if star, ok := gameContext.FindEntity("star").(*entities.StarProducer); ok {
					star.Update(&gameContext, delta)
					star.Draw(&gameContext)
				}
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
