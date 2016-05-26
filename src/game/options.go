// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/system/rumble"
)

// Options structure
type Options struct {
	Engine   *engine.Engine
	Resource *engine.Resource

	Fog  *Fog
	Dust *Dust

	Buttons      []*Button
	ButtonActive int

	YesText *Sprite
	NoText  *Sprite

	FadeTimer float64
}

// Returns new options
func NewOptions(e *engine.Engine, r *engine.Resource) (m *Options) {
	m = &Options{}
	m.Engine = e
	m.Resource = r

	m.Fog = NewFog(e, r)
	m.Dust = NewDust(e)

	return
}

// Initializes menu state
func (m *Options) OnInit() bool {
	m.Fog.Init()
	m.Dust.Init()

	// Create buttons
	m.Buttons = make([]*Button, 0)
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.MusicText, m.Resource.MusicTextHi, nil, m.Engine.Cfg.MusicEnabled))
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.SoundsText, m.Resource.SoundsTextHi, nil, m.Engine.Cfg.SoundsEnabled))

	if m.Engine.Joystick != nil {
		m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.AccelerometerText, m.Resource.AccelerometerTextHi, nil, m.Engine.Cfg.AccelerometerEnabled))
	}

	if rumble.RumbleAvailable() {
		m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.HapticText, m.Resource.HapticTextHi, nil, m.Engine.Cfg.HapticEnabled))
	}

	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.ShowFpsText, m.Resource.ShowFpsTextHi, nil, m.Engine.Cfg.ShowFps))

	m.ButtonActive = -1

	m.YesText = NewSprite(m.Engine, m.Resource.YesText)
	m.NoText = NewSprite(m.Engine, m.Resource.NoText)

	// Play menu music
	if !mix.PlayingMusic() {
		m.Resource.PlayMusic(m.Resource.MusicMenu, -1)
	}

	return true
}

// Quits game state
func (m *Options) OnQuit() bool {
	m.Engine.Cfg.Save()
	return true
}

// Returns state string
func (m *Options) String() string {
	return "Options"
}

// Handles input events
func (m *Options) HandleEvents() {
	if engine.Paused {
		event := sdl.WaitEvent()
		if event != nil {
			m.HandleEvent(event)
		}
	} else {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			m.HandleEvent(event)
		}
	}
}

// Handles input event
func (m *Options) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		// Handle quit event
		m.Engine.Quit()

	case *sdl.KeyDownEvent:

		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
			// Change state on back/escape
			m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
			m.Engine.State.Change(NewMenu(m.Engine, m.Resource))
		} else if (t.Keysym.Mod&sdl.KMOD_ALT != 0 && t.Keysym.Scancode == sdl.SCANCODE_RETURN) || t.Keysym.Scancode == sdl.SCANCODE_F11 {
			// Fullscreen
			m.Engine.Fullscreen()
		} else if t.Keysym.Scancode == sdl.SCANCODE_RETURN {
			// Change state on enter
			if m.ButtonActive != -1 {
				m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
				m.Buttons[m.ButtonActive].Selected = !m.Buttons[m.ButtonActive].Selected
				m.UpdateConfig()
			}
		} else if t.Keysym.Scancode == sdl.SCANCODE_UP || t.Keysym.Scancode == sdl.SCANCODE_DOWN {
			// Change active button on UP/DOWN
			if t.Keysym.Scancode == sdl.SCANCODE_UP {
				if m.ButtonActive == -1 {
					m.ButtonActive = len(m.Buttons) - 1
				} else {
					m.ButtonActive--
					if m.ButtonActive == -1 {
						m.ButtonActive = len(m.Buttons) - 1
					}
				}
			} else if t.Keysym.Scancode == sdl.SCANCODE_DOWN {
				if m.ButtonActive == -1 {
					m.ButtonActive = 0
				} else {
					m.ButtonActive++
					if m.ButtonActive == len(m.Buttons) {
						m.ButtonActive = 0
					}
				}
			}

			// Set active button
			for i := 0; i < len(m.Buttons); i++ {
				if i == m.ButtonActive {
					m.Buttons[i].Active = true
				} else {
					m.Buttons[i].Active = false
				}
			}
		}

	case *sdl.MouseMotionEvent:
		// Highlight button on hover
		point := sdl.Point{t.X, t.Y}
		for i := 0; i < len(m.Buttons); i++ {
			if point.InRect(m.Buttons[i].Image.Rect()) {
				for n := 0; n < len(m.Buttons); n++ {
					m.Buttons[n].Active = false
				}
				m.Buttons[i].Hovered = true
			} else {
				m.Buttons[i].Hovered = false
			}
		}

	case *sdl.MouseButtonEvent:
		// Set clicked button on mouse left button
		if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
			point := sdl.Point{t.X, t.Y}
			for i := 0; i < len(m.Buttons); i++ {
				if point.InRect(m.Buttons[i].Image.Rect()) {
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
					m.Buttons[i].Selected = !m.Buttons[i].Selected
					m.UpdateConfig()
				}
			}
		}

	case *sdl.TouchFingerEvent:
		// Set clicked button on touch
		point := sdl.Point{}

		// normalize touch coordinates
		point.X = int32(float64(t.X) * m.Engine.Cfg.WinWidth)
		point.Y = int32(float64(t.Y) * m.Engine.Cfg.WinHeight)

		if t.Type == sdl.FINGERDOWN {
			for i := 0; i < len(m.Buttons); i++ {
				if point.InRect(m.Buttons[i].Image.Rect()) {
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
					m.Buttons[i].Clicked = true
					m.Buttons[i].Selected = !m.Buttons[i].Selected
					m.UpdateConfig()
				} else {
					m.Buttons[i].Clicked = false
				}
			}
		}

	case *sdl.ControllerDeviceEvent:
		// Initialize/Remove controller
		if t.Type == sdl.CONTROLLERDEVICEADDED {
			m.Engine.Controller = sdl.GameControllerOpen(int(t.Which))
			if m.Engine.Cfg.HapticEnabled {
				m.Engine.SetHaptic()
			}
		} else if t.Type == sdl.CONTROLLERDEVICEREMOVED {
			m.Engine.CloseController()
		}

	case *sdl.ControllerButtonEvent:
		// Controller buttons
		if t.Type == sdl.CONTROLLERBUTTONDOWN {
			if t.Button == sdl.CONTROLLER_BUTTON_BACK {
				m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
				m.Engine.State.Change(NewMenu(m.Engine, m.Resource))
			} else if t.Button == sdl.CONTROLLER_BUTTON_A {
				if m.ButtonActive != -1 {
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
					m.Buttons[m.ButtonActive].Selected = !m.Buttons[m.ButtonActive].Selected
					m.UpdateConfig()
				}
			} else if t.Button == sdl.CONTROLLER_BUTTON_B || t.Button == sdl.CONTROLLER_BUTTON_BACK {
				m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
				m.Engine.State.Change(NewMenu(m.Engine, m.Resource))
			}
		}

	case *sdl.ControllerAxisEvent:
		// Controller axises
		switch t.Axis {
		case sdl.CONTROLLER_AXIS_LEFTY:
			// Change active button on UP/DOWN
			if t.Value < 0 {
				if m.ButtonActive == -1 {
					m.ButtonActive = len(m.Buttons) - 1
				} else {
					m.ButtonActive--
					if m.ButtonActive == -1 {
						m.ButtonActive = len(m.Buttons) - 1
					}
				}
			} else if t.Value > 0 {
				if m.ButtonActive == -1 {
					m.ButtonActive = 0
				} else {
					m.ButtonActive++
					if m.ButtonActive == len(m.Buttons) {
						m.ButtonActive = 0
					}
				}
			}

			for i := 0; i < len(m.Buttons); i++ {
				if i == m.ButtonActive {
					m.Buttons[i].Active = true
				} else {
					m.Buttons[i].Active = false
				}
			}

			break

		}

	default:
		break
	}
}

// Updates config
func (m *Options) UpdateConfig() {
	for i := 0; i < len(m.Buttons); i++ {
		switch m.Buttons[i].Image.Texture {
		case m.Resource.SoundsText:
			m.Engine.Cfg.SoundsEnabled = m.Buttons[i].Selected

		case m.Resource.AccelerometerText:
			m.Engine.Cfg.AccelerometerEnabled = m.Buttons[i].Selected

		case m.Resource.HapticText:
			m.Engine.Cfg.HapticEnabled = m.Buttons[i].Selected

		case m.Resource.ShowFpsText:
			m.Engine.Cfg.ShowFps = m.Buttons[i].Selected

		case m.Resource.MusicText:
			m.Engine.Cfg.MusicEnabled = m.Buttons[i].Selected

			if m.Buttons[i].Selected && !mix.PlayingMusic() {
				m.Resource.PlayMusic(m.Resource.MusicMenu, -1)
			} else if !m.Buttons[i].Selected && mix.PlayingMusic() {
				mix.HaltMusic()
			}
		}
	}
}

// Updates menu
func (m *Options) Update() {
	// Don't update if timer is paused
	if engine.Paused {
		return
	}

	// Scrolling
	m.Engine.ScreenDX = m.Engine.Cfg.BarrierSpeed

	// Update fadetimer
	m.FadeTimer += m.Engine.TFrame / 2.0

	// Update buttons
	for i := 0; i < len(m.Buttons); i++ {
		m.Buttons[i].Image.X = (m.Engine.Cfg.WinWidth-m.Buttons[i].Image.Width-m.NoText.Width)/2 + math.Cos(m.FadeTimer/6.5)*10
		m.Buttons[i].Image.Y = (m.Engine.Cfg.WinHeight/2 - (float64(len(m.Buttons)) * m.Buttons[0].Image.Height) + float64(i*70)) + math.Sin(m.FadeTimer/5.0)*10
		m.Buttons[i].Highlight.X = (m.Engine.Cfg.WinWidth-m.Buttons[i].Highlight.Width-m.NoText.Width)/2 + math.Cos(m.FadeTimer/6.5)*10
		m.Buttons[i].Highlight.Y = (m.Engine.Cfg.WinHeight/2 - (float64(len(m.Buttons)) * m.Buttons[0].Image.Height) + float64(i*70)) + math.Sin(m.FadeTimer/5.0)*10
	}

	// Update dust
	m.Dust.Update()

	// Update background
	m.Fog.Update()
}

// Draws menu
func (m *Options) Draw() {
	// Draw dust
	m.Dust.Draw()

	// Draw background
	m.Fog.Draw()

	// Draw buttons
	for i := 0; i < len(m.Buttons); i++ {
		m.Buttons[i].Draw()

		if m.Buttons[i].Selected {
			m.YesText.X = m.Buttons[i].Image.X + m.Buttons[i].Image.Width + m.YesText.Width/2
			m.YesText.Y = m.Buttons[i].Image.Y
			m.YesText.Draw()
		} else {
			m.NoText.X = m.Buttons[i].Image.X + m.Buttons[i].Image.Width + m.YesText.Width/2
			m.NoText.Y = m.Buttons[i].Image.Y
			m.NoText.Draw()
		}
	}

	// Show highlight on touch
	for i := 0; i < len(m.Buttons); i++ {
		if m.Buttons[i].Clicked == true {
			// Update screen
			m.Engine.Renderer.Present()

			m.Buttons[i].Clicked = false
		}
	}
}
