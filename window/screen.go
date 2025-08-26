// Package window
package window

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type OptsFunc func(*WindowOpts)

type WindowOpts struct {
	TickerDurationMil time.Duration
}

var (
	screen      tcell.Screen
	initOnce    sync.Once
	cleanupOnce sync.Once
	ticker      *time.Ticker
	style       tcell.Style
)

func ChangeTickerDuration(duration time.Duration) OptsFunc {
	return func(opts *WindowOpts) {
		opts.TickerDurationMil = duration
	}
}

func defautlOpts() WindowOpts {
	return WindowOpts{
		TickerDurationMil: 33,
	}
}

func InitScreen(opts ...OptsFunc) tcell.Screen {
	var s tcell.Screen
	var err error
	initOnce.Do(func() {
		o := defautlOpts()
		for _, fn := range opts {
			fn(&o)
		}
		s, err = tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		if err = s.Init(); err != nil {
			log.Fatal(err)
		}
		s.SetTitle("not set")
		ticker = time.NewTicker(o.TickerDurationMil * time.Millisecond)
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreenYellow)
		screen = s
	})
	return s
}

func GetStyle() tcell.Style {
	return style
}

func GetScreen() tcell.Screen {
	if screen != nil {
		log.Fatal("[SCREEN] Screen must be initialized first. Call InitScreen()")
	}
	return screen
}

func GetTicker() *time.Ticker {
	if ticker == nil {
		log.Fatal("[TICKER] Screen must be initialized first. Call InitScreen()")
	}
	return ticker
}

func GetSize() (int, int) {
	if screen == nil {
		log.Fatal("[SCREEN] Screen must be initialized first. Call InitScreen()")
	}
	return screen.Size()
}

func SetTitle(title string) {
	if screen == nil {
		log.Fatal("[TITLE] Screen must be initialized first. Call InitScreen()")
	}
	screen.SetTitle(title)
}

func InputEvent(exitCha chan int, keys func(tcell.Event)) {
	if screen == nil {
		log.Fatal("Screen must be initialized first. Call InitScreen()")
	}
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
	if screen == nil || ticker == nil {
		log.Fatal("Screen and/or ticker must be initialized first. Call InitScreen()")
		return
	}
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
