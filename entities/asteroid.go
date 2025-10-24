package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/entities/particles"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/game/design"
	"github.com/omar0ali/spaceinvaders-game-cli/game/loader"
)

type Asteroid struct {
	base.FallingObjectBase
	design.Design
}

type AsteroidProducer struct {
	Asteroids        []*Asteroid
	Design           design.AsteroidDesign
	Level            float64
	SelectedAsteroid *Asteroid
}

func NewAsteroidProducer(gc *game.GameContext) *AsteroidProducer {
	asteroidDesign, err := loader.LoadAsset[design.AsteroidDesign]("asteroids.json")
	if err != nil {
		panic(err)
	}

	a := &AsteroidProducer{
		Asteroids: []*Asteroid{},
		Level:     1.0,
		Design:    asteroidDesign,
	}

	if s, ok := gc.FindEntity("spaceship").(*SpaceShip); ok {
		s.AddOnLevelUp(func(newLevel int) {
			a.Level += 0.1
			game.Log(game.Warn, "Asteroid Level UP: %1.f", a.Level)
			// clear screen from asteroids when the player levels up
			a.Asteroids = nil
		})
	}

	return a
}

func (a *AsteroidProducer) Deploy() {
	w, _ := base.GetSize()

	pickAsteroid := a.Design.Asteroids[rand.Intn(len(a.Design.Asteroids))]

	width := len(pickAsteroid.Shape[0])
	height := len(pickAsteroid.Shape)

	speed := rand.Float64()*float64(min(a.Design.MaxSpeed, int(a.Level)+1)) + 2

	const padding = 30
	distance := (w - (padding * 2))
	xPos := rand.Intn(distance) + padding

	a.Asteroids = append(a.Asteroids, &Asteroid{
		FallingObjectBase: base.FallingObjectBase{
			ObjectBase: base.ObjectBase{
				Health:    pickAsteroid.EntityHealth + int(a.Level),
				MaxHealth: pickAsteroid.EntityHealth + int(a.Level),
				ObjectEntity: base.ObjectEntity{
					Position: base.PointFloat{X: float64(xPos), Y: -5},
					Width:    width,
					Height:   height,
					Speed:    speed,
				},
			},
		},
		Design: pickAsteroid,
	})
}

func (a *AsteroidProducer) Update(gc *game.GameContext, delta float64) {
	if len(a.Asteroids) < min(int(a.Level), a.Design.MaxLimit) {
		game.Log(game.Info, "Asteroids Deployed %d Level %.1f", len(a.Asteroids), a.Level)
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
			// asteroids can crash at each other and destroy
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
				a.SelectedAsteroid = asteroid
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
					particles.InitExplosion(15,
						particles.WithDimensions(
							asteroid.Position.X,
							asteroid.Position.Y,
							asteroid.Width,
							asteroid.Height,
						), particles.WithStyle(base.StyleIt(tcell.ColorWhite)),
					),
				)

				ps.AddParticles(particles.InitMeteroids(1,
					particles.WithDimensions(
						asteroid.Position.X,
						asteroid.Position.Y,
						asteroid.Width,
						asteroid.Height,
					)))
			}

			a.SelectedAsteroid = nil
			spaceship.ScoreHit()
		}

		if asteroid.IsOffScreen(h) {
			a.SelectedAsteroid = nil
		}

		if !asteroid.IsDead() && !asteroid.IsOffScreen(h) {
			activeAsteroids = append(activeAsteroids, asteroid)
		}
	}

	a.Asteroids = activeAsteroids
}

func (a *AsteroidProducer) Draw(gc *game.GameContext) {
	// display the last asteroid that was shot

	if a.SelectedAsteroid != nil {
		base.DisplayHealthLeft(
			&a.SelectedAsteroid.ObjectBase,
			8,
			a.SelectedAsteroid.Name,
			15,
			true,
			base.StyleIt(a.SelectedAsteroid.GetColor()),
			nil,
		)
	}

	for _, asteroid := range a.Asteroids {
		color := base.StyleIt(asteroid.GetColor())
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
