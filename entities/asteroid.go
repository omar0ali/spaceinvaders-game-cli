package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

const (
	MaxAsteroidsDeployed = 4
	MaxSpeed             = 7
)

type Asteroid struct {
	base.FallingObjectBase
	game.Design
}

type AsteroidProducer struct {
	Asteroids []*Asteroid
	Level     float64
}

func NewAsteroidProducer(gc *game.GameContext) *AsteroidProducer {
	a := &AsteroidProducer{
		Asteroids: []*Asteroid{},
		Level:     1.0,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			a.Level += 0.1
		})
	}

	return a
}

func (a *AsteroidProducer) Deploy() {
	w, _ := base.GetSize()
	const padding = 23
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding

	designs, err := game.LoadListOfAssets[game.Design]("asteroids.json")
	if err != nil {
		panic(err)
	}

	design := designs[rand.Intn(len(designs))]
	width := len(design.Shape[0])
	height := len(design.Shape)

	speed := rand.Float64()*float64(min(MaxAsteroidsDeployed, a.Level+1)) + 2

	a.Asteroids = append(a.Asteroids, &Asteroid{
		FallingObjectBase: base.FallingObjectBase{
			ObjectBase: base.ObjectBase{
				Health:    design.EntityHealth + int(a.Level),
				MaxHealth: design.EntityHealth + int(a.Level),
				ObjectEntity: base.ObjectEntity{
					Position: base.PointFloat{X: float64(xPos), Y: -5},
					Width:    width,
					Height:   height,
					Speed:    speed,
				},
			},
		},
		Design: design,
	})
}

func (a *AsteroidProducer) Update(gc *game.GameContext, delta float64) {
	if len(a.Asteroids) < min(int(a.Level), MaxAsteroidsDeployed) {
		a.Deploy()
	}

	var spaceship *SpaceShip
	var alienProducer *AlienProducer

	if ship, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		spaceship = ship
	}

	if aliens, ok := gc.FindEntity("alien").(*AlienProducer); ok {
		alienProducer = aliens
	}

	activeAsteroids := a.Asteroids[:0]

	for i, asteroid := range a.Asteroids {
		for j, asteroid2 := range a.Asteroids {
			if i >= j {
				continue
			}
			if Crash(&asteroid.ObjectBase, &asteroid2.ObjectEntity, gc) {
				asteroid.TakeDamage(40)
				asteroid2.TakeDamage(40)
			}
		}
		Move(&asteroid.ObjectBase, delta)

		// can collid with a meteroid
		if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
			for _, p := range ps.ParticleProducable {
				switch p.(type) {
				case *particles.MeteroidProducer:
					for _, m := range p.GetParticles() {
						if Crash(&asteroid.ObjectBase, &m.ObjectEntity, gc) {
							asteroid.TakeDamage(1)
							p.RemoveParticle(m)
						}
					}
				}
			}
		}

		// get get hit from the spaceship
		for _, beam := range spaceship.GetBeams() {
			if GettingHit(asteroid, beam, gc) {
				asteroid.TakeDamage(spaceship.GetPower())
				spaceship.RemoveBeam(beam)
			}
		}

		// get get hit from the aliens
		for _, alien := range alienProducer.Aliens {
			for _, beam := range alien.GetBeams() {
				if GettingHit(asteroid, beam, gc) {
					asteroid.TakeDamage(alien.GetPower())
					alien.RemoveBeam(beam)
				}
			}
		}

		_, h := base.GetSize()

		if asteroid.IsDead() {
			if ps, ok := gc.FindEntity("particles").(*particles.ParticleSystem); ok {
				ps.AddParticles(
					particles.InitExplosion(10,
						particles.WithDimensions(
							asteroid.Position.X,
							asteroid.Position.Y,
							asteroid.Width,
							asteroid.Height,
						), particles.WithStyle(base.StyleIt(tcell.ColorReset, tcell.ColorWhite)),
					),
				)

				ps.AddParticles(particles.InitMeteroids(
					particles.WithDimensions(
						asteroid.Position.X,
						asteroid.Position.Y,
						asteroid.Width,
						asteroid.Height,
					)))
			}

			spaceship.ScoreHit()
		}

		if !asteroid.IsDead() && !asteroid.IsOffScreen(h) {
			activeAsteroids = append(activeAsteroids, asteroid)
		}
	}

	a.Asteroids = activeAsteroids
}

func (a *AsteroidProducer) Draw(gc *game.GameContext) {
	for _, asteroid := range a.Asteroids {
		color := base.StyleIt(tcell.ColorReset, asteroid.GetColor())
		// asteroid.DisplayHealth(5, true, color, nil)
		for rowIndex, line := range asteroid.Shape {
			for colIndex, char := range line {
				if char != ' ' {
					x := int(asteroid.Position.GetX()) + colIndex
					y := int(asteroid.Position.GetY()) + rowIndex
					base.SetContentWithStyle(x, y, char, color)
				}
			}
		}

	}
}

func (a *AsteroidProducer) InputEvents(event tcell.Event, gc *game.GameContext) {
	// testing mode

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'i' { // dev mode to create a star
	// 		a.Deploy()
	// 	}
	// }
}

func (a *AsteroidProducer) GetType() string {
	return "asteroid"
}
