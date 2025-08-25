// Package window
package window

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	screen      tcell.Screen
	initOnce    sync.Once
	initLock    sync.Mutex
	cleanupOnce sync.Once
	ticker      *time.Ticker
	style       tcell.Style
	delta       float64
)

func initScreen() {
	initOnce.Do(func() {
		var s tcell.Screen
		var err error
		s, err = tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		if err = s.Init(); err != nil {
			log.Fatal(err)
		}
		s.SetTitle("not set")
		ticker = time.NewTicker(33 * time.Millisecond)
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreenYellow)
		screen = s
	})
}

func GetScreen() tcell.Screen {
	initLock.Lock()
	defer initLock.Unlock()
	if screen != nil {
		return screen
	}
	initScreen()
	return screen
}

func GetTicker() *time.Ticker {
	initLock.Lock()
	defer initLock.Unlock()
	if ticker == nil {
		// ensures both ticker and screen are initialized
		initScreen()
	}
	return ticker
}

func SetTitle(title string) {
	initScreen()
	screen.SetTitle(title)
}

func InputEvent(exitCha chan int, keys func(tcell.Event)) {
	initScreen()
	go func() {
		for {
			event := screen.PollEvent()
			switch ev := event.(type) {
			case *tcell.EventResize:
				screen.Clear()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyESC || ev.Rune() == 'q' {
					exitCha <- 0
					cleanupOnce.Do(func() {
						screen.Fini()
					})
					return
				}
			}
			keys(event)
		}
	}()
}

func Update(exitCha chan int, updates func(delta float64)) {
	initScreen()
	var delta float64
	go func() {
		last := time.Now()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				delta = now.Sub(last).Seconds()
				last = now

				screen.Clear()

				lenStr := []rune(fmt.Sprintf("FPS: %.2f", (1 / delta)))
				for i, r := range lenStr {
					screen.SetContent(i, 0, r, nil, style)
				}

				updates(delta)

				screen.Show()
			case val := <-exitCha:
				if val == 0 {
					cleanupOnce.Do(func() {
						screen.Fini()
					})
					return
				}
			}
		}
	}()
}
