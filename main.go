// Package main
package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/ui"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

func DeployEntities(gc *game.GameContext, cfg game.GameConfig, exitCha chan struct{}) {
	// order is important since some objects might overlap others
	gc.AddEntity(entities.NewStarsProducer(cfg))
	gc.AddEntity(entities.NewSpaceShip(cfg, gc))
	gc.AddEntity(entities.NewModifierProducer(gc))
	if cfg.Dev.Asteroids { // includeing asteroids is optional
		gc.AddEntity(entities.NewAsteroidProducer(gc))
	}
	gc.AddEntity(entities.NewAlienProducer(gc))
	gc.AddEntity(entities.NewBossAlienProducer(gc))
	gc.AddEntity(particles.NewParticleSystem())
	gc.AddEntity(ui.NewUISystem())
	gc.AddEntity(entities.NewUI(gc, exitCha))
}

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

	// ------------------------------------- Objects ----------------------------------
	gameContext := game.GameContext{
		Screen: screen,
	}
	// ---------------------------------- entities --------------------------------------

	log.Println("Game running...")
	DeployEntities(&gameContext, cfg, exit)

	// ----------------------------------------- window ------------------------------------
	base.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlR {
					gameContext.RemoveAllEntities()
					DeployEntities(&gameContext, cfg, exit)
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
