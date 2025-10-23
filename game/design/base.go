// Package design
package design

import "github.com/gdamore/tcell/v2"

type Design struct {
	Name         string   `json:"name"`
	Shape        []string `json:"shape"`
	EntityHealth int      `json:"health"`
	Color        string   `json:"color"`
	Speed        int      `json:"speed"`
}

type Designable interface {
	GetColor() tcell.Color
	GetHealth() int
	GetName() string
	GetShape() []string
	GetMaxSpeed() int
}

func (d *Design) GetColor() tcell.Color { return HexToColor(d.Color) }
func (d *Design) GetHealth() int        { return d.EntityHealth }
func (d *Design) GetName() string       { return d.Name }
func (d *Design) GetShape() []string    { return d.Shape }
func (d *Design) GetMaxSpeed() int      { return d.Speed }
