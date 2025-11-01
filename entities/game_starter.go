package entities

import "github.com/omar0ali/spaceinvaders-game-cli/game"
import "github.com/omar0ali/spaceinvaders-game-cli/game/design"
import "github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
import "github.com/omar0ali/spaceinvaders-game-cli/entities/ui"

func StartGame(gc *game.GameContext, cfg game.GameConfig, exitCha chan struct{}) {
	// order is important since some objects might overlap others
	loadedUIDesigns := design.LoadDesigns()
	gc.AddEntity(NewStarsProducer(cfg))
	gc.AddEntity(NewSpaceShip(cfg, gc, loadedUIDesigns))
	gc.AddEntity(NewModifierProducer(gc, loadedUIDesigns))
	if cfg.Dev.Asteroids { // includeing asteroids is optional
		gc.AddEntity(NewAsteroidProducer(gc, loadedUIDesigns))
	}
	gc.AddEntity(NewAlienProducer(gc, loadedUIDesigns))
	gc.AddEntity(NewBossAlienProducer(gc, loadedUIDesigns))
	gc.AddEntity(particles.NewParticleSystem())
	gc.AddEntity(ui.NewUISystem())
	gc.AddEntity(NewUI(gc, cfg, exitCha))
}

func RestartGame(gc *game.GameContext, cfg game.GameConfig, exitCha chan struct{}) {
	gc.RemoveAllEntities()
	StartGame(gc, cfg, exitCha)
}
