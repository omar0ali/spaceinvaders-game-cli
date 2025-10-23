package design

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func HexToColor(hex string) tcell.Color {
	if len(hex) != 6 {
		return tcell.ColorWhite
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
