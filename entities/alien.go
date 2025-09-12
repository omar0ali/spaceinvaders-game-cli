package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Alien struct {
	FallingObjectBase
}

type AlienProducer struct {
	Aliens []*Alien
	limit  int // NOTE: always starts with 0 | its linked to the spaceship level up
	health int
	Cfg    core.GameConfig
}

func NewAlienProducer(cfg core.GameConfig) *AlienProducer {
	return &AlienProducer{
		limit:  cfg.AliensConfig.Limit, // start limit
		health: cfg.AliensConfig.Health,
		Cfg:    cfg,
	}
}

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// limit amount of ships Falling (generate alien ships)
	if len(a.Aliens) < a.limit-1 {
		a.DeployAliens()
	}

	// -------- this will ensure to clean up dead aliens and beams --------
	spaceship := a.MovementAndCollision(delta, gc)

	// -------- progression ---------
	spaceship.LevelUp(func() { // on every spaceship level up, deployed alien health increases
		a.limit += 1
		a.health += 1
	})
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	brownColor := window.StyleIt(tcell.ColorReset, tcell.ColorBrown)
	color := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)
	for _, alien := range a.Aliens {
		// drawing the points
		// header
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX()), int(alien.OriginPoint.GetY())+alien.Height+1, 'v', color,
		)

		// draw the left lines
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())-1, int(alien.OriginPoint.GetY())+alien.Height, '\\', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())-2, int(alien.OriginPoint.GetY())+alien.Height-1, '\\', brownColor,
		)
		// draw the right lines
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())+1, int(alien.OriginPoint.GetY())+alien.Height, '/', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())+2, int(alien.OriginPoint.GetY())+alien.Height-1, '/', brownColor,
		)
		// draw the bottom lines
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX()), int(alien.OriginPoint.GetY())+alien.Height-2, '^', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())+1, int(alien.OriginPoint.GetY())+alien.Height-2, '-', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())-1, int(alien.OriginPoint.GetY())+alien.Height-2, '-', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())+2, int(alien.OriginPoint.GetY())+alien.Height-2, '-', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())-2, int(alien.OriginPoint.GetY())+alien.Height-2, '-', brownColor,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())+3, int(alien.OriginPoint.GetY())+alien.Height-2, ']', color,
		)
		window.SetContentWithStyle(
			int(alien.OriginPoint.GetX())-3, int(alien.OriginPoint.GetY())+alien.Height-2, '[', color,
		)

	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'm' { // dev mode
			a.DeployAliens()
		}
	}
}

func (a *AlienProducer) DeployAliens() {
	w, _ := window.GetSize()
	const padding = 18
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding // starting from 18
	randSpeed := rand.Intn(a.Cfg.AliensConfig.Speed) + 2
	// spawn alien
	a.Aliens = append(a.Aliens, &Alien{
		FallingObjectBase: FallingObjectBase{
			Health:      a.health,
			Speed:       randSpeed,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       3,
			Height:      5,
		},
	})
}

func (a *AlienProducer) UIAlienShipData(gc *core.GameContext) {
	w, _ := window.GetSize()
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	aliensStr := []rune(fmt.Sprintf("Aliens Limit: %d * ", a.limit-1))
	for i, r := range aliensStr {
		window.SetContentWithStyle(w+i-len(aliensStr), 2, r, whiteColor)
	}
	alienMSPD := []rune(fmt.Sprintf("Max SPD: %d * ", a.Cfg.AliensConfig.Speed))
	for i, r := range alienMSPD {
		window.SetContentWithStyle(w+i-len(alienMSPD), 3, r, whiteColor)
	}
	aliensHP := []rune(fmt.Sprintf("Max HP: %d * ", a.health))
	for i, r := range aliensHP {
		window.SetContentWithStyle(w+i-len(aliensHP), 4, r, whiteColor)
	}
}

func (a *AlienProducer) MovementAndCollision(delta float64, gc *core.GameContext) *SpaceShip {
	var activeAliens []*Alien
	var gun *Gun
	var spaceship *SpaceShip

	// look for the spaceship since it has the gun and the number of beams
	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
		gun = &spaceship.Gun
	}

	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.Aliens {
		// Update the coordinates of the aliens.
		alien.move(delta)
		for _, beam := range gun.Beams {
			if alien.isHit(&beam.position, gun.Power) {
				spaceship.ScoreHit()
				gun.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien
		_, h := window.GetSize()
		if alien.isOffScreen(h) {
			spaceship.Health -= 1
		}
		if alien.isDead() {
			spaceship.ScoreKill()
		}
		if !alien.isDead() && !alien.isOffScreen(h) { // still flying
			activeAliens = append(activeAliens, alien)
		}
	}
	a.Aliens = activeAliens
	return spaceship
}

func (a *AlienProducer) GetType() string {
	return "alien"
}
