// VoV game
package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/system/home"
	"github.com/gen2brain/vov/src/system/log"
)

var fileId string = ""
var jsonKey = []byte(``)

// Score structure
type Score struct {
	Name      string
	Time      int
	X         float64
	Y         float64
	Width     float64
	Height    float64
	Formatted string
}

// Scores structure
type Scores struct {
	Engine   *engine.Engine
	Resource *engine.Resource

	Fog  *Fog
	Dust *Dust

	// Current score
	Current int

	// String from input
	TextInput string

	// Is passed score highscore
	IsHighScore bool

	// Continue to menu state after timeout if true
	Continue bool

	// If scores are loaded
	Loaded bool

	// New hiscore text
	HiScoreText      *Sprite
	HiScoreEnterText *Sprite

	// Loading text
	LoadingText *Sprite

	// List of scores
	Scores []Score

	// Timers
	FadeTimer  float64
	StateTimer float64
}

// Returns new scores
func NewScores(e *engine.Engine, r *engine.Resource, score int, cont bool) (s *Scores) {
	s = &Scores{}

	s.Engine = e
	s.Resource = r
	s.Current = score
	s.Continue = cont

	s.Fog = NewFog(e, r)
	s.Dust = NewDust(e)

	return
}

// Initializes state
func (s *Scores) OnInit() bool {
	s.Fog.Init()
	s.Dust.Init()

	// Initialize scores
	s.Scores = make([]Score, s.Engine.Cfg.NScores)

	s.HiScoreText = NewSprite(s.Engine, s.Resource.HiScoreText)
	s.HiScoreEnterText = NewSprite(s.Engine, s.Resource.HiScoreEnterText)

	s.HiScoreText.X = (s.Engine.Cfg.WinWidth - s.HiScoreText.Width) / 2
	s.HiScoreText.Y = 100

	s.HiScoreEnterText.X = (s.Engine.Cfg.WinWidth - s.HiScoreEnterText.Width) / 2
	s.HiScoreEnterText.Y = s.HiScoreText.Y + s.HiScoreText.Height + s.HiScoreEnterText.Height

	s.LoadingText = NewSprite(s.Engine, s.Resource.LoadingText)
	s.LoadingText.X = s.Engine.Cfg.WinWidth/2 - (s.LoadingText.Width / 2)
	s.LoadingText.Y = s.Engine.Cfg.WinHeight/2 - (s.LoadingText.Height / 2)

	// Play music
	if !mix.PlayingMusic() {
		s.Resource.PlayMusic(s.Resource.MusicMenu, -1)
	}

	// Load scores
	if s.Exists() {
		s.Load()
	} else {
		s.Default()
	}

	s.Format()

	// Check score
	if s.Loaded {
		s.IsHighScore = s.HighScore()
		if s.IsHighScore {
			// Accept text input
			sdl.StartTextInput()
		}
	}

	if s.Continue {
		s.StateTimer = s.Engine.Cfg.ScoresLength
	}

	return true
}

// Quits state
func (s *Scores) OnQuit() bool {
	return true
}

// Returns state string
func (s *Scores) String() string {
	return "Scores"
}

// Handles input events
func (s *Scores) HandleEvents() {
	if engine.Paused {
		event := sdl.WaitEvent()
		if event != nil {
			s.HandleEvent(event)
		}
	} else {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			s.HandleEvent(event)
		}
	}
}

// Handles input event
func (s *Scores) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		// Handle quit event
		s.Engine.Quit()

	case *sdl.KeyDownEvent:
		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
			// Change state on back/escape
			s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)
			s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
		} else if (t.Keysym.Mod&sdl.KMOD_ALT != 0 && t.Keysym.Scancode == sdl.SCANCODE_RETURN) || t.Keysym.Scancode == sdl.SCANCODE_F11 {
			// Fullscreen
			s.Engine.Fullscreen()
		} else if t.Keysym.Scancode == sdl.SCANCODE_BACKSPACE {
			// Handle backspace
			if sdl.IsTextInputActive() && s.TextInput != "" {
				s.TextInput = s.TextInput[:len(s.TextInput)-1]
			}
		} else if t.Keysym.Scancode == sdl.SCANCODE_RETURN {
			// Stop accepting input on enter and save
			if s.IsHighScore && sdl.IsTextInputActive() && s.TextInput != "" {
				sdl.StopTextInput()

				// Insert highscore
				s.Insert()

				// Save highscore
				s.Save()

				s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)

				s.Current = 0
				s.Continue = false
				s.IsHighScore = false

				// Switch to menu state after duration
				s.Continue = true
				s.StateTimer = s.Engine.Cfg.ScoresLength
			} else {
				s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)
				s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
			}
		}

	case *sdl.MouseButtonEvent:
		if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
			// Change state on mouse button
			s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)
			s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
		}

	case *sdl.TouchFingerEvent:
		if t.Type == sdl.FINGERDOWN {
			// Change state on touch
			s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)
			s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
		}

	case *sdl.ControllerDeviceEvent:
		// Initialize/Remove controller
		if t.Type == sdl.CONTROLLERDEVICEADDED {
			s.Engine.Controller = sdl.GameControllerOpen(int(t.Which))
			if s.Engine.Cfg.HapticEnabled {
				s.Engine.SetHaptic()
			}
		} else if t.Type == sdl.CONTROLLERDEVICEREMOVED {
			s.Engine.CloseController()
		}

	case *sdl.ControllerButtonEvent:
		// Controller buttons
		if t.Type == sdl.CONTROLLERBUTTONDOWN {
			if t.Button == sdl.CONTROLLER_BUTTON_B || t.Button == sdl.CONTROLLER_BUTTON_BACK {
				s.Resource.PlaySound(s.Resource.SoundClick, -1, 0)
				s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
			}
		}

	case *sdl.TextInputEvent:
		// Enter name for highscore
		if t.Type == sdl.TEXTINPUT {
			if len(s.TextInput) < 16 {
				b := t.Text[:]
				s.TextInput += string(b[:clen(b)])
			}
		}

	default:
		break
	}
}

// Adds default scores
func (s *Scores) Default() {
	for i := 0; i < s.Engine.Cfg.NScores; i++ {
		s.Scores[i] = Score{"-", 150000 - (i * 15000), 0, 0, 0, 0, ""}
	}

	s.Loaded = true
}

// Returns score rank
func (s *Scores) Rank() int {
	for i := 0; i < s.Engine.Cfg.NScores; i++ {
		if s.Current > s.Scores[i].Time {
			return i
		}
	}
	return -1
}

// Checks if score is highscore
func (s *Scores) HighScore() bool {
	if s.Current <= s.Scores[s.Engine.Cfg.NScores-1].Time {
		return false
	}

	return s.Rank() >= 0
}

// Formats scores and gets text dimensions
func (s *Scores) Format() {
	max := 0
	for i := 0; i < s.Engine.Cfg.NScores; i++ {
		s.Scores[i].Formatted = formatTime(s.Scores[i].Time, true)

		w, _, _ := s.Resource.FontMedium.SizeUTF8(s.Scores[i].Name)
		if w > max {
			max = w
		}
	}

	for i := 0; i < s.Engine.Cfg.NScores; i++ {
		w1, h1, _ := s.Resource.FontSmall.SizeUTF8("0.")
		w2, _, _ := s.Resource.FontMedium.SizeUTF8(s.Scores[0].Formatted)
		s.Scores[i].Width, s.Scores[i].Height = float64(w1+w2+max), float64(h1)
	}
}

// Inserts new score
func (s *Scores) Insert() {
	rank := s.Rank()
	if rank == -1 {
		return
	}

	if s.Scores[rank].Name != "-" {
		// Move all lower scores down
		for i := s.Engine.Cfg.NScores - 1; i > rank; i-- {
			t := s.Scores[i-1]
			s.Scores[i] = t
		}
	}

	s.Scores[rank].Time = s.Current
	s.Scores[rank].Name = s.TextInput

	s.Format()
}

// Loads scores from file
func (s *Scores) Load() {
	file := filepath.Join(home.Dir(), ".vov", "scores")
	js, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("ReadFile: %s\n", err)
	}

	err = json.Unmarshal(js, &s.Scores)
	if err != nil {
		log.Error("Unmarshal: %s\n", err)
	}

	s.Loaded = true
}

// Saves scores to file
func (s *Scores) Save() {
	dir := filepath.Join(home.Dir(), ".vov")
	if _, err := os.Stat(dir); err != nil {
		os.Mkdir(dir, 0755)
	}

	js, err := json.MarshalIndent(s.Scores, "", "    ")
	if err != nil {
		log.Error("MarshalIndent: %s\n", err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "scores"), js, 0644)
	if err != nil {
		log.Error("WriteFile: %s\n", err)
	}
}

// Checks if scores file exists
func (s *Scores) Exists() bool {
	file := filepath.Join(home.Dir(), ".vov", "scores")
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// Updates state
func (s *Scores) UpdateState() {
	s.StateTimer -= float64(s.Engine.FrameDelta)
	if s.StateTimer > 0 {
		return
	}

	if s.Continue && !s.IsHighScore {
		s.Engine.State.Change(NewMenu(s.Engine, s.Resource))
	}
}

// Updates scores
func (s *Scores) Update() {
	// Don't update if timer is paused
	if engine.Paused {
		return
	}

	// Scrolling
	s.Engine.ScreenDX = s.Engine.Cfg.BarrierSpeed
	s.Engine.ScreenDY = 0.0

	// Update fadetimer
	s.FadeTimer += s.Engine.TFrame / 2.0

	// Update text
	for i := 0; i < s.Engine.Cfg.NScores; i++ {
		s.Scores[i].X = (s.Engine.Cfg.WinWidth-s.Scores[i].Width)/2 + math.Cos(s.FadeTimer/6.5)*10
		s.Scores[i].Y = (s.Engine.Cfg.WinHeight/2 - (float64(s.Engine.Cfg.NScores) * float64(s.Scores[i].Height)) + float64(i*40)) + math.Sin(s.FadeTimer/5.0)*10
	}

	// Update dust
	s.Dust.Update()

	// Update background
	s.Fog.Update()
}

// Draws scores
func (s *Scores) Draw() {
	// Draw dust
	s.Dust.Draw()

	// Draw background
	s.Fog.Draw()

	if !s.IsHighScore {
		// Draw scores
		if s.Loaded {
			for i := 0; i < s.Engine.Cfg.NScores; i++ {
				x := int32(s.Scores[i].X)
				y := int32(s.Scores[i].Y)

				s.Resource.DrawText(fmt.Sprintf("%d.", i+1), x, y+1, engine.FONT_SMALL)
				s.Resource.DrawText(s.Scores[i].Formatted, x+50, y, engine.FONT_MEDIUM)
				s.Resource.DrawText(s.Scores[i].Name, x+150, y, engine.FONT_MEDIUM)
			}
		} else {
			// Draw loading screen
			s.LoadingText.Draw()
		}
	} else {
		// Draw input screen
		s.HiScoreText.Draw()
		s.HiScoreEnterText.Draw()

		w, h, _ := s.Resource.FontMain.SizeUTF8(s.TextInput)
		x := (s.Engine.Cfg.WinWidth - float64(w)) / 2
		y := s.HiScoreEnterText.Y + s.HiScoreEnterText.Height + float64(h)
		s.Resource.DrawText(s.TextInput, int32(x), int32(y), engine.FONT_LARGE)
	}

	// Update state
	s.UpdateState()
}
