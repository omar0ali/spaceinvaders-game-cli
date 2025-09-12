package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/core"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Star struct {
	FallingObjectBase
}

type StarProducer struct {
	Stars []*Star
	Limit int
	Cfg   core.GameConfig
}

func NewStarsProducer(cfg core.GameConfig) *StarProducer {
	return &StarProducer{
		Stars: []*Star{},
		Limit: max(cfg.StarsConfig.Limit, 15),
		Cfg:   cfg,
	}
}

func (s *StarProducer) Deployment() {
	w, _ := window.GetSize()
	xPos := rand.Intn(w)

	randSpeed := rand.Intn(max(s.Cfg.StarsConfig.Speed, 30)) + 10 // ensure the speed its always high < 30 larger than 30
	s.Stars = append(s.Stars, &Star{
		FallingObjectBase: FallingObjectBase{
			Health:      1,
			Speed:       randSpeed,
			OriginPoint: core.PointFloat{X: float64(xPos), Y: -5},
			Width:       1,
			Height:      1,
		},
	})
}

func (s *StarProducer) Update(gc *core.GameContext, delta float64) {
	// create stars at least 15
	if len(s.Stars) < s.Limit {
		s.Deployment()
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
			s.Deployment()
		}
	}
}

func (s *StarProducer) GetType() string {
	return "star"
}
