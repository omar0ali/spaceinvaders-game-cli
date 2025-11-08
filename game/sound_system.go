package game

import (
	"bytes"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/omar0ali/spaceinvaders-game-cli/game/assets"
)

// to limit how many sounds playing at the same time.
var (
	soundsPlaying int32
	maxSounds     int32 = 10
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type Sound struct {
	Data   []byte
	Format beep.Format
}

type SoundSystem struct {
	Sounds map[string]Sound
	cfg    GameConfig
}

func InitSoundSystem(cfg GameConfig) *SoundSystem {
	if !cfg.Dev.Sounds {
		return &SoundSystem{}
	}
	var sounds = map[string]Sound{}
	entries, _ := assets.SoundFS.ReadDir("sounds")
	for _, e := range entries {
		name := e.Name()
		data, _ := assets.SoundFS.ReadFile("sounds/" + name)
		r := nopCloser{bytes.NewBuffer(data)}

		Log(Info, "LOAD SOUND: %s", name)
		streamer, format, _ := mp3.Decode(r)
		streamer.Close()
		sounds[name] = Sound{Data: data, Format: format}
	}

	// prepare speaker only once.
	sampleRate := beep.SampleRate(44100) // usign mp3
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	return &SoundSystem{
		Sounds: sounds,
		cfg:    cfg,
	}
}

func (s *SoundSystem) PlayRandom(names []string, vol float64) {
	if len(names) == 0 {
		return
	}
	idx := rand.Intn(len(names))
	s.PlaySound(names[idx], vol)
}

func (s *SoundSystem) PlaySound(name string, vol float64) {
	if !s.cfg.Dev.Sounds {
		return
	}
	sound, ok := s.Sounds[name]
	if !ok {
		Log(Error, "Failed to locate the file. %s ", name)
	}

	if atomic.LoadInt32(&soundsPlaying) >= maxSounds {
		Log(Debug, "Skipping sound: too many playing")
		return
	}

	Log(Debug, "Sounds Playing: %s", name)

	go func() {
		atomic.AddInt32(&soundsPlaying, 1)
		r := nopCloser{bytes.NewReader(sound.Data)}
		streamer, _, err := mp3.Decode(r)
		if err != nil {
			Log(Error, "Failed to decode %s", err)
			return
		}

		v := &effects.Volume{
			Streamer: streamer,
			Base:     2,
			Volume:   vol,
			Silent:   false,
		}
		speaker.Play(beep.Seq(v, beep.Callback(func() {
			streamer.Close()
			atomic.AddInt32(&soundsPlaying, -1)
		})))
	}()
}
