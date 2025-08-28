// Package core
package core

import "github.com/gdamore/tcell/v2"

type (
	Entity interface {
		Draw(gc *GameContext)
		Update(gc *GameContext, delta float64)
		InputEvents(event tcell.Event, gc *GameContext)
		GetType() string
	}
	GameContext struct {
		entities []Entity
		Screen   tcell.Screen
	}
)

func (gs *GameContext) AddEntity(entity ...Entity) {
	gs.entities = append(gs.entities, entity...)
}

func (gs *GameContext) RemoveEntity(entity Entity) {
	for i, v := range gs.entities {
		if v == entity {
			gs.entities = append(gs.entities[:i], gs.entities[i+1:]...)
			return
		}
	}
}

func (gs *GameContext) GetEntities() []Entity {
	return gs.entities
}

func (gs *GameContext) FindEntity(typeName string) Entity {
	for _, entity := range gs.entities {
		if entity.GetType() == typeName {
			return entity
		}
	}
	return nil
}
