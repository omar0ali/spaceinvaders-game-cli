package game

import (
	"bytes"
	"io"
	"math/rand"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/omar0ali/spaceinvaders-game-cli/game/assets"
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
		return
	}
	Log(Debug, "Sounds Playing: %s", name)

	go func() {
		r := nopCloser{bytes.NewReader(sound.Data)}
		streamer, format, err := mp3.Decode(r)
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
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(beep.Seq(v, beep.Callback(func() {
			streamer.Close()
		})))
	}()
}
