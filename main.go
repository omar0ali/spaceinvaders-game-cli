// Package main
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/window"
)

func main() {
	exit := make(chan int)

	window.SetTitle("Space Invader Game")

	window.InputEvent(exit,
		func(event tcell.Event) {
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Rune() == 'r' {
					// testing exit with letter r
					exit <- 0
				}
			}
		},
	)

	window.Update(exit,
		func(delta float64) {
		},
	)

	// exit
	if val := <-exit; val == 0 {
		return
	}
}
