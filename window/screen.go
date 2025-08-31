// Package window
package window

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvader-game-cli/core"
)

type OptsFunc func(*WindowOpts)

type WindowOpts struct {
	TickerDurationMil time.Duration
	EnableMouse       bool
}

var (
	screen      tcell.Screen
	initOnce    sync.Once
	cleanupOnce sync.Once
	ticker      *time.Ticker
	style       tcell.Style
	Delta       float64
)

func ChangeTickerDuration(duration time.Duration) OptsFunc {
	return func(opts *WindowOpts) {
		opts.TickerDurationMil = duration
	}
}

func EnableMouse(opts *WindowOpts) {
	opts.EnableMouse = true
}

func defautlOpts() WindowOpts {
	return WindowOpts{
		TickerDurationMil: 33,
	}
}

func InitScreen(opts ...OptsFunc) tcell.Screen {
	var err error
	initOnce.Do(func() {
		o := defautlOpts()
		for _, fn := range opts {
			fn(&o)
		}
		screen, err = tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		if err = screen.Init(); err != nil {
			log.Fatal(err)
		}
		// enable mouse
		if o.EnableMouse {
			screen.EnableMouse()
		}
		screen.SetTitle("not set")
		ticker = time.NewTicker(o.TickerDurationMil * time.Millisecond)
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreenYellow)
	})
	return screen
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

func InputEvent(exitCha chan struct{}, keys func(tcell.Event)) {
	if screen == nil {
		log.Fatal("[InputEvent] Screen must be initialized first. Call InitScreen()")
	}
	go func() {
		for {
			event := screen.PollEvent()
			switch ev := event.(type) {
			case *tcell.EventResize:
				screen.Clear()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyESC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
					cleanupOnce.Do(func() {
						screen.Fini()
					})
					close(exitCha)
					return
				}
			}
			keys(event)
		}
	}()
}

func Update(exitCha chan struct{}, updates func(delta float64)) {
	if screen == nil || ticker == nil {
		log.Fatal("Screen and/or ticker must be initialized first. Call InitScreen()")
	}
	go func() {
		last := time.Now()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				Delta = now.Sub(last).Seconds()
				last = now

				screen.Clear()

				lenStr := []rune(fmt.Sprintf("FPS: %.2f", (1 / Delta)))
				for i, r := range lenStr {
					screen.SetContent(i, 0, r, nil, style)
				}

				updates(Delta)

				screen.Show()
			case <-exitCha:
				cleanupOnce.Do(func() {
					screen.Fini()
				})
				return
			}
		}
	}()
}

func SetContent(x, y int, r rune) {
	if screen == nil {
		log.Fatal("[SetContent] Screen must be initialized first. Call InitScreen()")
	}
	screen.SetContent(x, y, r, nil, style)
}

func SetContentWithStyle(x, y int, r rune, style tcell.Style) {
	if screen == nil {
		log.Fatal("[SetContentWithStyle] Screen must be initialized first. Call InitScreen()")
	}
	screen.SetContent(x, y, r, nil, style)
}

func StyleIt(background, forground tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(background).Foreground(forground)
}

func GetCenterPoint() core.Point {
	if screen == nil {
		log.Fatal("[GetCenterPoint] Screen must be initialized first. Call InitScreen()")
	}
	w, h := GetSize()
	return core.Point{X: w / 2, Y: h / 2}
}
