package entities

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

type FallingObjectBase struct {
	Health        int
	Speed         int
	OriginPoint   core.PointFloat
	TrianglePoint core.Triangle
}

func (f *FallingObjectBase) isOffScreen(h int) bool {
	return int(f.OriginPoint.Y) > h-2
}

func (f *FallingObjectBase) isDead() bool {
	return f.Health <= 0
}

func (f *FallingObjectBase) isHit(point core.PointInterface, power int) bool {
	grayColor := window.StyleIt(tcell.ColorReset, tcell.ColorDarkGray)
	redColor := window.StyleIt(tcell.ColorReset, tcell.ColorRed)
	yellowColor := window.StyleIt(tcell.ColorReset, tcell.ColorYellow)

	if f.TrianglePoint.A.GetY() > point.GetY() &&
		f.TrianglePoint.C.GetY()-2 < point.GetY() &&
		(f.TrianglePoint.C.GetX()-1 < point.GetX() && f.TrianglePoint.B.GetX()+1 > point.GetX()) {

		f.Health -= power // update health of the falling object

		window.SetContentWithStyle(
			int(point.GetX()-1), int(point.GetY()+1), tcell.RuneBoard, grayColor)
		window.SetContentWithStyle(
			int(point.GetX()-1), int(point.GetY()), tcell.RuneCkBoard, yellowColor)
		window.SetContentWithStyle(
			int(point.GetX()+1), int(point.GetY()), tcell.RuneBoard, grayColor)
		window.SetContentWithStyle(
			int(point.GetX()), int(point.GetY()+1), tcell.RuneCkBoard, redColor)
		window.SetContentWithStyle(
			int(point.GetX()), int(point.GetY()-1), tcell.RuneBoard, yellowColor)
		window.SetContentWithStyle(
			int(point.GetX()+1), int(point.GetY()+1), tcell.RuneCkBoard, grayColor)
		return true
	}
	return false
}

func (f *FallingObjectBase) move(delta float64) {
	distance := float64(f.Speed) * delta
	f.OriginPoint.Y += distance
	f.TrianglePoint.A.AppendY(distance)
	f.TrianglePoint.B.AppendY(distance)
	f.TrianglePoint.C.AppendY(distance)
}

func NewObject(health, speed int, origin core.PointFloat) *FallingObjectBase {
	return &FallingObjectBase{
		Health:      health,
		Speed:       speed,
		OriginPoint: origin,
		TrianglePoint: core.Triangle{
			A: &core.PointFloat{
				X: origin.X,
				Y: origin.Y + 2,
			},
			B: &core.PointFloat{
				X: origin.X + 3,
				Y: origin.Y,
			},
			C: &core.PointFloat{
				X: origin.X - 3,
				Y: origin.Y,
			},
		},
	}
}

type FallingObjects interface {
	move(distance float64)
	isHit(point core.PointInterface)
	NewObject(health, speed int, origin core.PointFloat)
}
