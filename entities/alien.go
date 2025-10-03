package entities

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Alien struct {
	FallingObjectBase
	Gun
	core.Design
	done chan struct{}
}

type AlienProducer struct {
	Aliens []*Alien
	limit  float64 // always starts with 0 | its linked to the spaceship level up
	health float64
	Cfg    core.GameConfig
}

func NewAlienProducer(cfg core.GameConfig, gc *core.GameContext) *AlienProducer {
	a := &AlienProducer{
		limit:  float64(cfg.AliensConfig.Limit),
		health: float64(cfg.AliensConfig.Health),
		Cfg:    cfg,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			a.limit += 0.2
			a.health += 0.8
		})
	}

	return a
}

func (a *AlienProducer) Update(gc *core.GameContext, delta float64) {
	// limit amount of ships Falling (generate alien ships)
	if len(a.Aliens) < int(a.limit) {
		a.DeployAliens()
	}

	// go through each alien's gun and shoot
	for _, alien := range a.Aliens {
		alien.Update(gc, delta)
	}

	// -------- this will ensure to clean up dead aliens and beams --------
	a.MovementAndCollision(delta, gc)
}

func (a *AlienProducer) Draw(gc *core.GameContext) {
	colorHealth := window.StyleIt(tcell.ColorReset, tcell.ColorIndianRed)
	for _, alien := range a.Aliens {
		color := window.StyleIt(tcell.ColorReset, alien.GetColor())
		alien.Draw(gc, alien.GetColor())

		alien.DisplayHealth(6, true, colorHealth)

		// draw shape
		for rowIndex, line := range alien.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(alien.OriginPoint.GetX()) + colIndex
					y := int(alien.OriginPoint.GetY()) + rowIndex
					window.SetContentWithStyle(x, y, char, color)
				}
			}
		}

	}
}

func (a *AlienProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	// testing code

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'm' { // dev mode
	// 		a.DeployAliens()
	// 	}
	// }
}

func (a *AlienProducer) DeployAliens() {
	w, _ := window.GetSize()
	const padding = 20
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding // starting from 18
	randSpeed := rand.Intn(a.Cfg.AliensConfig.Speed) + 2
	// spawn alien
	designs, err := core.LoadListOfAssets[core.Design]("alienship.json")
	if err != nil {
		panic(err)
	}

	// pick random design: based on the current health level. The higher the stronger the ships.
	randDesign := designs[rand.Intn(int(a.health))]
	width := len(randDesign.Shape[0])
	height := len(randDesign.Shape)

	alien := &Alien{
		FallingObjectBase: FallingObjectBase{
			ObjectBase: ObjectBase{
				Health:      int(a.health) + randDesign.EntityHealth,
				MaxHealth:   int(a.health) + randDesign.EntityHealth,
				Width:       width,
				Height:      height,
				OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			},
			Speed: randSpeed,
		},
		Gun: Gun{
			Beams: []*Beam{},
			Cap:   1,
			Power: a.Cfg.AliensConfig.GunPower,
			Speed: a.Cfg.AliensConfig.GunSpeed,
		},
		Design: randDesign,
	}

	go DoEvery(2*time.Second,
		func() {
			alien.initBeam(core.Point{
				X: int(alien.OriginPoint.X) + (alien.Width / 2),
				Y: int(alien.OriginPoint.Y) + (alien.Height) + 1,
			}, Down)
		},
		alien.done,
	)

	a.Aliens = append(a.Aliens, alien)
}

func (a *AlienProducer) UIAlienShipData(gc *core.GameContext) {
	w, _ := window.GetSize()
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	aliensStr := []rune(fmt.Sprintf("Aliens Limit: %d * ", int(a.limit)))
	for i, r := range aliensStr {
		window.SetContentWithStyle(w+i-len(aliensStr), 2, r, whiteColor)
	}
	alienMSPD := []rune(fmt.Sprintf("Max SPD: %d * ", a.Cfg.AliensConfig.Speed))
	for i, r := range alienMSPD {
		window.SetContentWithStyle(w+i-len(alienMSPD), 3, r, whiteColor)
	}
}

func (a *AlienProducer) MovementAndCollision(delta float64, gc *core.GameContext) *SpaceShip {
	var activeAliens []*Alien
	var spaceship *SpaceShip

	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	// on each alien avaiable check its position and check if the beam is at the same position
	for _, alien := range a.Aliens {
		// Update the coordinates of the aliens.
		alien.move(delta)
		for _, beam := range spaceship.Beams {
			if alien.isHit(&beam.position, spaceship.Power) {
				spaceship.ScoreHit()
				spaceship.RemoveBeam(beam) // removing a beam when hitting the ship
			}
		}

		// check the alien ship height position
		// check the health of each alien

		_, h := window.GetSize()
		if alien.isOffScreen(h) {
			spaceship.Health -= 1
		}
		if alien.isDead() {
			ScoreKill()
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
