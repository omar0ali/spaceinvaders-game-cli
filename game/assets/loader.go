// Package assets
package assets

import "embed"

//go:embed *.json
var Files embed.FS

//go:embed sounds/*.mp3
var SoundFS embed.FS
