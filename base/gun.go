package base

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/omar0ali/spaceinvaders-game-cli/game"
)

type Direction = int

const (
	Up Direction = iota
	Down
)

type beam struct {
	position  Point
	Symbol    rune
	Direction Direction
}

func (b *beam) GetPosition() *Point {
	return &b.position
}

type Gun struct {
	beams  []*beam
	cap    int
	loaded int
	power  int
	speed  int

	reloading       bool
	mu              sync.Mutex
	lastShot        time.Time
	cooldown        time.Duration
	reloadCooldown  time.Duration
	reloadStartTime time.Time
}

func NewGun(cap, power, speed int, cooldown, reloadCooldown int) Gun {
	return Gun{
		beams:          []*beam{},
		cap:            cap,
		loaded:         cap,
		power:          power,
		speed:          speed,
		cooldown:       time.Duration(cooldown) * time.Millisecond,
		reloadCooldown: time.Duration(reloadCooldown) * time.Millisecond,
	}
}

func (g *Gun) GetPower() int {
	return g.power
}

func (g *Gun) GetSpeed() int {
	return g.speed
}

func (g *Gun) GetCooldown() time.Duration {
	return g.cooldown / time.Millisecond
}

func (g *Gun) GetReloadCooldown() time.Duration {
	return g.reloadCooldown / time.Millisecond
}

func (g *Gun) IncreaseGunSpeed(i, limit int) bool {
	if g.speed < limit {
		g.speed += i
		g.speed = min(g.speed, limit)
		return true
	}
	return false
}

func (g *Gun) IncreaseGunPower(i int) bool {
	g.power += i
	return true
}

func (g *Gun) DecreaseCooldown(i int) bool {
	if g.cooldown < 30 { // min cooldown set to 30
		return false
	}
	g.cooldown += time.Duration(i) * time.Millisecond
	return true
}

func (g *Gun) DecreaseGunReloadCooldown(i int) bool {
	if g.reloadCooldown < 30 { // min cooldown set to 30
		return false
	}
	g.reloadCooldown += time.Duration(i) * time.Millisecond
	return true
}

func (g *Gun) IncreaseGunCap(i, limit int) bool {
	if g.cap < limit {
		g.cap += i
		g.cap = min(g.cap, limit)
		return true
	}
	return false
}

func (g *Gun) GetCapacity() int {
	return g.cap
}

func (g *Gun) GetLoaded() int {
	return g.loaded
}

func (g *Gun) GetBeams() []*beam {
	return g.beams
}

func (g *Gun) IsReloading() bool {
	return g.reloading
}

func (g *Gun) ReloadGun(sounds *game.SoundSystem) {
	if !g.reloading {
		g.reloading = true
		sounds.PlaySound("sfx-tank-reload.mp3", 0)
		done := make(chan struct{})
		go DoOnce(g.reloadCooldown, func() {
			g.mu.Lock()
			defer g.mu.Unlock()
			g.loaded = g.cap
			g.reloading = false
		}, done)
	}
}

func (g *Gun) InitBeam(pos Point, dir Direction, sounds *game.SoundSystem) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.IsReloading() {
		return
	}

	if time.Since(g.lastShot) < g.cooldown {
		return
	}

	if g.loaded <= 0 {
		g.ReloadGun(sounds)
		return
	}

	symbol := '↑'
	if dir == Down {
		symbol = '↓'
	}

	beam := beam{
		position: Point{
			X: pos.X,
			Y: pos.Y,
		},
		Symbol:    symbol,
		Direction: dir,
	}

	// sounds.PlaySound("8-bit-explosion-1.mp3", -1)
	sounds.PlaySound("8-bit-laser.mp3", -1)

	g.beams = append(g.beams, &beam)
	g.lastShot = time.Now()
	g.loaded -= 1
}

func (g *Gun) RemoveBeam(beam *beam) {
	for i, b := range g.beams {
		if beam == b {
			g.beams = append(g.beams[:i], g.beams[i+1:]...)
			break
		}
	}
}

func (g *Gun) Update(gc *game.GameContext, delta float64) {
	// update the coordinates of the beam
	_, h := GetSize()
	var activeBeams []*beam
	for _, beam := range g.beams {
		distance := int(float64(g.speed) * delta)
		switch beam.Direction {
		case Up:
			beam.position.Y -= distance
		case Down:
			beam.position.Y += distance
		}
		if beam.position.Y >= 0 && beam.position.Y <= h {
			activeBeams = append(activeBeams, beam)
		}
	}

	g.beams = activeBeams
}

func (g *Gun) Draw(gc *game.GameContext, color tcell.Color) {
	// draw the beam new position
	style := StyleIt(color)

	for _, beam := range g.beams {
		SetContentWithStyle(beam.position.X, beam.position.Y, beam.Symbol, style)
	}
}

func (g *Gun) InputEvents(event tcell.Event, gc *game.GameContext) {}

func (g *Gun) GetType() string {
	return "gun"
}
