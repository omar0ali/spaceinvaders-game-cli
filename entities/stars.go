package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type Star struct {
	FallingObjectBase
}

type StarProducer struct {
	Stars []*Star
}

func (s *StarProducer) Update(gc *core.GameContext, delta float64) {
	// create stars at least 15
	if len(s.Stars) < 15 {
		w, _ := window.GetSize()
		xPos := rand.Intn(w)
		randSpeed := rand.Intn(50) + 10
		// create star
		s.Stars = append(s.Stars, &Star{
			FallingObjectBase: *NewObject(1, randSpeed, core.PointFloat{X: float64(xPos), Y: -5}),
		})
	}

	// Update the coordinates of the stars.
	for _, star := range s.Stars {
		star.move(delta)
	}

	// -------- this will ensure to clean up stars --------

	var activeStars []*Star

	// on each star avaiable check its position
	for _, star := range s.Stars {
		// check the star height position
		// clear
		_, h := window.GetSize()
		if !star.isOffScreen(h) {
			activeStars = append(activeStars, star)
		}
	}

	s.Stars = activeStars
}

func (s *StarProducer) Draw(gc *core.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, tcell.ColorWhite)
	for _, star := range s.Stars {
		switch {
		case star.Speed < 15:
			window.SetContentWithStyle(int(star.OriginPoint.GetX()), int(star.OriginPoint.GetY()), '◆', whiteColor)
		case star.Speed >= 15 && star.Speed < 25:
			window.SetContentWithStyle(int(star.OriginPoint.GetX()), int(star.OriginPoint.GetY()), '+', whiteColor)
		case star.Speed >= 25 && star.Speed < 45:
			window.SetContentWithStyle(int(star.OriginPoint.GetX()), int(star.OriginPoint.GetY()), '*', whiteColor)
		case star.Speed >= 45 && star.Speed < 50:
			window.SetContentWithStyle(int(star.OriginPoint.GetX()), int(star.OriginPoint.GetY()), tcell.RuneDegree, whiteColor)
		default:
			window.SetContentWithStyle(int(star.OriginPoint.GetX()), int(star.OriginPoint.GetY()), '.', whiteColor)
		}
	}
}

func (s *StarProducer) InputEvents(event tcell.Event, gc *core.GameContext) {
	switch ev := event.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'n' { // dev mode to create a star
			w, _ := window.GetSize()
			xPos := rand.Intn(w)
			randSpeed := rand.Intn(50) + 4
			// create star
			s.Stars = append(s.Stars, &Star{
				FallingObjectBase: *NewObject(1, randSpeed, core.PointFloat{X: float64(xPos), Y: -5}),
			})
		}
	}
}

func (s *StarProducer) GetType() string {
	return "star"
}
