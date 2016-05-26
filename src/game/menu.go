// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
)

// Menu structure
type Menu struct {
	Engine   *engine.Engine
	Resource *engine.Resource

	Fog  *Fog
	Dust *Dust

	Buttons      []*Button
	ButtonActive int

	TitleText *Sprite

	FadeTimer float64
}

// Menu button
type Button struct {
	Image     *Sprite
	Highlight *Sprite

	Active   bool
	Hovered  bool
	Clicked  bool
	Selected bool

	State engine.State
}

// Draws button
func (b *Button) Draw() {
	if b.Active || b.Hovered || b.Clicked {
		b.Highlight.Draw()
	} else {
		b.Image.Draw()
	}
}

// Returns new button
func NewButton(e *engine.Engine, i *sdl.Texture, h *sdl.Texture, s engine.State, selected bool) (b *Button) {
	b = &Button{}
	b.State = s

	b.Image = NewSprite(e, i)
	b.Highlight = NewSprite(e, h)
	b.Selected = selected

	return
}

// Returns new menu
func NewMenu(e *engine.Engine, r *engine.Resource) (m *Menu) {
	m = &Menu{}
	m.Engine = e
	m.Resource = r

	m.Fog = NewFog(e, r)
	m.Dust = NewDust(e)

	return
}

// Initializes menu state
func (m *Menu) OnInit() bool {
	m.Fog.Init()
	m.Dust.Init()

	// Create buttons
	m.Buttons = make([]*Button, 0)
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.StartText, m.Resource.StartTextHi, NewGame(m.Engine, m.Resource), false))
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.ScoresText, m.Resource.ScoresTextHi, NewScores(m.Engine, m.Resource, 0, false), false))
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.OptionsText, m.Resource.OptionsTextHi, NewOptions(m.Engine, m.Resource), false))
	m.Buttons = append(m.Buttons, NewButton(m.Engine, m.Resource.CreditsText, m.Resource.CreditsTextHi, NewCredits(m.Engine, m.Resource), false))

	m.ButtonActive = -1

	// Create sprite from rendered text
	m.TitleText = NewSprite(m.Engine, m.Resource.TitleText)

	// Play menu music
	if !mix.PlayingMusic() {
		m.Resource.PlayMusic(m.Resource.MusicMenu, -1)
	}

	return true
}

// Quits game state
func (m *Menu) OnQuit() bool {
	return true
}

// Returns state string
func (m *Menu) String() string {
	return "Menu"
}

// Handles input events
func (m *Menu) HandleEvents() {
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
func (m *Menu) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		// Handle quit event
		m.Engine.Quit()

	case *sdl.KeyDownEvent:

		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
			// Handle quit event
			m.Engine.Quit()
		} else if (t.Keysym.Mod&sdl.KMOD_ALT != 0 && t.Keysym.Scancode == sdl.SCANCODE_RETURN) || t.Keysym.Scancode == sdl.SCANCODE_F11 {
			// Fullscreen
			m.Engine.Fullscreen()
		} else if t.Keysym.Scancode == sdl.SCANCODE_RETURN {
			// Change state on enter
			m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
			if m.ButtonActive == -1 {
				m.Engine.State.Change(NewGame(m.Engine, m.Resource))
			} else {
				m.Engine.State.Change(m.Buttons[m.ButtonActive].State)
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
					m.Buttons[i].Clicked = true
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
				} else {
					m.Buttons[i].Clicked = false
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
					m.Buttons[i].Clicked = true
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
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
				m.Engine.Quit()
			} else if t.Button == sdl.CONTROLLER_BUTTON_A {
				if m.ButtonActive != -1 {
					m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
					m.Engine.State.Change(m.Buttons[m.ButtonActive].State)
				}
			} else if t.Button == sdl.CONTROLLER_BUTTON_B {
				for i := 0; i < len(m.Buttons); i++ {
					m.Buttons[i].Active = false
				}
			} else if t.Button == sdl.CONTROLLER_BUTTON_START {
				m.Resource.PlaySound(m.Resource.SoundClick, -1, 0)
				m.Engine.State.Change(NewGame(m.Engine, m.Resource))
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

// Updates menu
func (m *Menu) Update() {
	// Don't update if timer is paused
	if engine.Paused {
		return
	}

	// Scrolling
	m.Engine.ScreenDX = m.Engine.Cfg.BarrierSpeed
	m.Engine.ScreenDY = 0.0

	// Update fadetimer
	m.FadeTimer += m.Engine.TFrame / 2.0

	// Update title
	m.TitleText.X = (m.Engine.Cfg.WinWidth-m.TitleText.Width)/2 + math.Cos(m.FadeTimer/6.5)*10
	m.TitleText.Y = (m.Engine.Cfg.WinHeight/2 - m.TitleText.Height - (float64(len(m.Buttons)) * m.Buttons[0].Image.Height)) + math.Sin(m.FadeTimer/5.0)*10

	// Update buttons
	for i := 0; i < len(m.Buttons); i++ {
		m.Buttons[i].Image.X = (m.Engine.Cfg.WinWidth-m.Buttons[i].Image.Width)/2 + math.Cos(m.FadeTimer/6.5)*10
		m.Buttons[i].Image.Y = (m.Engine.Cfg.WinHeight/2 - m.TitleText.Height - m.Buttons[i].Image.Height + float64(i*70) + 30) + math.Sin(m.FadeTimer/5.0)*10
		m.Buttons[i].Highlight.X = (m.Engine.Cfg.WinWidth-m.Buttons[i].Highlight.Width)/2 + math.Cos(m.FadeTimer/6.5)*10
		m.Buttons[i].Highlight.Y = (m.Engine.Cfg.WinHeight/2 - m.TitleText.Height - m.Buttons[i].Highlight.Height + float64(i*70) + 30) + math.Sin(m.FadeTimer/5.0)*10
	}

	// Update dust
	m.Dust.Update()

	// Update background
	m.Fog.Update()
}

// Draws menu
func (m *Menu) Draw() {
	// Draw dust
	m.Dust.Draw()

	// Draw background
	m.Fog.Draw()

	// Draw title
	m.TitleText.Draw()

	// Draw buttons
	for i := 0; i < len(m.Buttons); i++ {
		m.Buttons[i].Draw()
	}

	// Change state if button is clicked
	for i := 0; i < len(m.Buttons); i++ {
		if m.Buttons[i].Clicked == true {
			// Update screen
			m.Engine.Renderer.Present()

			// Change state
			m.Engine.State.Change(m.Buttons[i].State)
		}
	}
}
