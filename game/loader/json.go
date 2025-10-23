// Package loader
package loader

import (
	"encoding/json"

	"github.com/omar0ali/spaceinvaders-game-cli/game/assets"
)

func LoadAsset[T any](filePath string) (T, error) {
	file, err := assets.Files.Open(filePath)
	var design T
	if err != nil {
		return design, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&design); err != nil {
		return design, err
	}

	return design, nil
}

func LoadListOfAssets[T any](filePath string) ([]T, error) {
	file, err := assets.Files.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []T
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
