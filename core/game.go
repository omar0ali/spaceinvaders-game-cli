// Package core
package core

import (
	"github.com/gdamore/tcell/v2"
)

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
		Halt     bool
	}
)

func (gc *GameContext) AddEntity(entity ...Entity) {
	gc.entities = append(gc.entities, entity...)
}

func (gc *GameContext) RemoveEntity(entity Entity) {
	for i, v := range gc.entities {
		if v == entity {
			gc.entities = append(gc.entities[:i], gc.entities[i+1:]...)
			return
		}
	}
}

func (gc *GameContext) RemoveAllEntities() {
	gc.entities = []Entity{}
}

func (gc *GameContext) GetEntities() []Entity {
	return gc.entities
}

func (gc *GameContext) FindEntity(typeName string) Entity {
	for _, entity := range gc.entities {
		if entity.GetType() == typeName {
			return entity
		}
	}
	return nil
}
