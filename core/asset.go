package core

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

type Design struct {
	Name         string   `json:"name"`
	Shape        []string `json:"shape"`
	EntityHealth int      `json:"health"`
	Color        string   `json:"color"`
}

func (d *Design) GetColor() tcell.Color {
	return HexToColor(d.Color)
}

func LoadAsset(filePath string) (Design, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Design{}, err
	}
	defer file.Close()

	var rawDesign Design
	if err := json.NewDecoder(file).Decode(&rawDesign); err != nil {
		return Design{}, err
	}

	return rawDesign, nil
}

func LoadListOfAssets(filePath string) ([]Design, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rawDesigns []Design
	if err := json.NewDecoder(file).Decode(&rawDesigns); err != nil {
		return nil, err
	}

	return rawDesigns, nil
}

func HexToColor(hex string) tcell.Color {
	if len(hex) != 6 {
		return tcell.ColorWhite
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
