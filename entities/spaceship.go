// Package entities
package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
)

type SpaceShip struct{}

func (s *SpaceShip) Update(gs *core.GameContext) {}
func (s *SpaceShip) Draw(gs *core.GameContext)   {}
func (s *SpaceShip) InputEvents(event tcell.Event, gc *core.GameContext) {
}

func (s *SpaceShip) GetType() string {
	return "spaceship"
}
