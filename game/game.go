// Package game
package game

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type LogType string

const (
	Info  LogType = "INFO"
	Warn  LogType = "WARN"
	Error LogType = "ERROR"
	Debug LogType = "DEBUG"
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

func Log(logType LogType, format string, v ...any) {
	cfg := LoadConfig()
	if !cfg.Dev.Debug {
		return
	}
	log.Printf(string("["+logType+"] ")+format, v...)
}

func SetupLogs() *os.File {
	f, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile) // optional: timestamp + file info
	log.Println("Starting game...")
	return f
}
