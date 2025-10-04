package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/base"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
	"github.com/omar0ali/spaceinvaders-game-cli/window"
)

type Star struct {
	base.FallingObjectBase
}

type StarProducer struct {
	Stars []*Star
	Limit int
	Cfg   game.GameConfig
}

func NewStarsProducer(cfg game.GameConfig) *StarProducer {
	return &StarProducer{
		Stars: []*Star{},
		Limit: max(cfg.StarsConfig.Limit, 15),
		Cfg:   cfg,
	}
}

func (s *StarProducer) Deployment() {
	w, _ := window.GetSize()
	xPos := rand.Intn(w)

	randSpeed := rand.Float64()*float64(max(s.Cfg.StarsConfig.Speed, 15)) + 10

	s.Stars = append(s.Stars, &Star{
		FallingObjectBase: base.FallingObjectBase{
			ObjectBase: base.ObjectBase{
				Health:   1,
				Position: game.PointFloat{X: float64(xPos), Y: -5},
				Width:    1,
				Height:   1,
				Speed:    randSpeed,
			},
		},
	})
}

func (s *StarProducer) Update(gc *game.GameContext, delta float64) {
	if len(s.Stars) < s.Limit {
		s.Deployment()
	}

	// Update the coordinates of the stars.
	for _, star := range s.Stars {
		star.Move(delta)
	}

	// -------- this will ensure to clean up stars --------

	var activeStars []*Star

	// on each star avaiable check its position
	for _, star := range s.Stars {
		// check the star height position
		// clear
		_, h := window.GetSize()
		if !star.IsOffScreen(h) {
			activeStars = append(activeStars, star)
		}
	}

	s.Stars = activeStars
}

func (s *StarProducer) Draw(gc *game.GameContext) {
	whiteColor := window.StyleIt(tcell.ColorReset, game.HexToColor("445559"))
	for _, star := range s.Stars {
		switch {
		case star.Speed < 15:
			window.SetContentWithStyle(int(star.Position.GetX()), int(star.Position.GetY()), '☼', whiteColor)
		case star.Speed >= 15 && star.Speed < 25:
			window.SetContentWithStyle(int(star.Position.GetX()), int(star.Position.GetY()), '+', whiteColor)
		case star.Speed >= 25 && star.Speed < 45:
			window.SetContentWithStyle(int(star.Position.GetX()), int(star.Position.GetY()), '*', whiteColor)
		case star.Speed >= 45 && star.Speed < 50:
			window.SetContentWithStyle(int(star.Position.GetX()), int(star.Position.GetY()), tcell.RuneDegree, whiteColor)
		default:
			window.SetContentWithStyle(int(star.Position.GetX()), int(star.Position.GetY()), '.', whiteColor)
		}
	}
}

func (s *StarProducer) InputEvents(event tcell.Event, gc *game.GameContext) {
	// testing mode

	// switch ev := event.(type) {
	// case *tcell.EventKey:
	// 	if ev.Rune() == 'n' { // dev mode to create a star
	// 		s.Deployment()
	// 	}
	// }
}

func (s *StarProducer) GetType() string {
	return "star"
}
