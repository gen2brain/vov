// VoV game
package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
)

// Credit structure
type Credit struct {
	Role *Sprite
	Name *Sprite
}

// Returns new credit
func NewCredit(e *engine.Engine, r *sdl.Texture, n *sdl.Texture) (c *Credit) {
	c = &Credit{}

	c.Role = NewSprite(e, r)
	c.Name = NewSprite(e, n)

	return
}

// Credits structure
type Credits struct {
	Engine   *engine.Engine
	Resource *engine.Resource

	Fog  *Fog
	Dust *Dust

	ScrollSpeed float64

	Credits []*Credit
}

// Returns new credits
func NewCredits(e *engine.Engine, r *engine.Resource) (c *Credits) {
	c = &Credits{}
	c.Engine = e
	c.Resource = r

	c.Fog = NewFog(e, r)
	c.Dust = NewDust(e)

	return
}

// Initializes state
func (c *Credits) OnInit() bool {
	c.Fog.Init()
	c.Dust.Init()

	c.ScrollSpeed = c.Engine.Cfg.ScrollSpeed

	c.Credits = make([]*Credit, 0)
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.ProgrammingText, c.Resource.ProgrammingCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.GraphicsText, c.Resource.GraphicsCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.MusicAndSoundsText, c.Resource.MusicAndSoundsCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.FontText, c.Resource.FontCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.BasedText, c.Resource.BasedCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.SDLText, c.Resource.SDLCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.GoText, c.Resource.GoCreditText))
	c.Credits = append(c.Credits, NewCredit(c.Engine, c.Resource.VoVText, c.Resource.VoVCreditText))

	for i := 0; i < len(c.Credits); i++ {
		c.Credits[i].Role.X = (c.Engine.Cfg.WinWidth - c.Credits[i].Role.Width) / 2
		c.Credits[i].Name.X = (c.Engine.Cfg.WinWidth - c.Credits[i].Name.Width) / 2

		c.Credits[i].Role.Y = c.Engine.Cfg.WinHeight + c.Credits[i].Role.Height + float64(i*120)
		c.Credits[i].Name.Y = c.Engine.Cfg.WinHeight + c.Credits[i].Role.Height + c.Credits[i].Name.Height + float64(i*120)
	}

	if !mix.PlayingMusic() {
		c.Resource.PlayMusic(c.Resource.MusicMenu, -1)
	}

	return true
}

// Quits state
func (c *Credits) OnQuit() bool {
	return true
}

// Returns state string
func (c *Credits) String() string {
	return "Credits"
}

// Handles input events
func (c *Credits) HandleEvents() {
	if engine.Paused {
		event := sdl.WaitEvent()
		if event != nil {
			c.HandleEvent(event)
		}
	} else {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			c.HandleEvent(event)
		}
	}
}

// Handles input event
func (c *Credits) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		// Handle quit event
		c.Engine.Quit()

	case *sdl.KeyDownEvent:
		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
			// Change state on back/escape
			c.Resource.PlaySound(c.Resource.SoundClick, -1, 0)
			c.Engine.State.Change(NewMenu(c.Engine, c.Resource))
		} else if (t.Keysym.Mod&sdl.KMOD_ALT != 0 && t.Keysym.Scancode == sdl.SCANCODE_RETURN) || t.Keysym.Scancode == sdl.SCANCODE_F11 {
			// Fullscreen
			c.Engine.Fullscreen()
		}

	case *sdl.ControllerDeviceEvent:
		// Initialize/Remove controller
		if t.Type == sdl.CONTROLLERDEVICEADDED {
			c.Engine.Controller = sdl.GameControllerOpen(int(t.Which))
			if c.Engine.Cfg.HapticEnabled {
				c.Engine.SetHaptic()
			}
		} else if t.Type == sdl.CONTROLLERDEVICEREMOVED {
			c.Engine.CloseController()
		}

	case *sdl.ControllerButtonEvent:
		// Controller buttons
		if t.Type == sdl.CONTROLLERBUTTONDOWN {
			if t.Button == sdl.CONTROLLER_BUTTON_B || t.Button == sdl.CONTROLLER_BUTTON_BACK {
				c.Resource.PlaySound(c.Resource.SoundClick, -1, 0)
				c.Engine.State.Change(NewMenu(c.Engine, c.Resource))
			}
		}

	case *sdl.TouchFingerEvent:
		// Pause scroll on touch
		if t.Type == sdl.FINGERDOWN {
			c.ScrollSpeed = 0
		} else if t.Type == sdl.FINGERUP {
			c.ScrollSpeed = c.Engine.Cfg.ScrollSpeed
		}

	case *sdl.MouseButtonEvent:
		// Pause scroll on click
		if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
			c.ScrollSpeed = 0
		} else if t.Type == sdl.MOUSEBUTTONUP && t.Button == sdl.BUTTON_LEFT {
			c.ScrollSpeed = c.Engine.Cfg.ScrollSpeed
		}

	default:
		break
	}
}

// Updates credits
func (c *Credits) Update() {
	// Don't update if timer is paused
	if engine.Paused {
		return
	}

	// Scrolling
	c.Engine.ScreenDX = c.Engine.Cfg.BarrierSpeed
	c.Engine.ScreenDY = 0.0

	// Update dust
	c.Dust.Update()

	// Update background
	c.Fog.Update()

	// Update text
	for i := 0; i < len(c.Credits); i++ {
		c.Credits[i].Role.Y -= c.ScrollSpeed
		c.Credits[i].Name.Y -= c.ScrollSpeed

		if c.Credits[i].Role.Y < -(c.Engine.Cfg.WinHeight * 2.0) {
			c.Credits[i].Role.Y = c.Engine.Cfg.WinHeight
		}

		if c.Credits[i].Name.Y < -(c.Engine.Cfg.WinHeight * 2.0) {
			c.Credits[i].Name.Y = c.Engine.Cfg.WinHeight
		}
	}
}

// Draws credits
func (c *Credits) Draw() {
	// Draw dust
	c.Dust.Draw()

	// Draw background
	c.Fog.Draw()

	// Draw text
	for i := 0; i < len(c.Credits); i++ {
		c.Credits[i].Role.Draw()
		c.Credits[i].Name.Draw()
	}
}
