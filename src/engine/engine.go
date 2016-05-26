// VoV engine
package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	// Haptic device
	Haptic *sdl.Haptic
)

// SDL engine structure
type Engine struct {
	// Config
	Cfg *Config

	// Game state
	State StateMachine

	// SDL Window
	Window *sdl.Window
	// SDL Renderer
	Renderer *sdl.Renderer

	// Game controller
	Controller *sdl.GameController

	// Joystick
	Joystick *sdl.Joystick

	// Boolean set to true until exit
	Running bool

	// Screen X distance
	ScreenDX float64
	// Screen Y distance
	ScreenDY float64

	// Frames counter
	Frames int
	// Frames per second
	Fps float64
	// Maximum fps in ms
	FrameMs uint32

	// Start of frame (milliseconds)
	StartTicks uint32
	// End of frame (milliseconds)
	EndTicks uint32
	// Delta time between frames
	FrameDelta uint32
	// Length of frame adjusted for gamespeed
	TFrame float64
}

// Returns new engine
func NewEngine(c *Config) (e *Engine) {
	e = &Engine{}
	e.Cfg = c
	e.Running = true

	StartTimer()

	return
}

// Initializes engine
func (e *Engine) Init() (err error) {

	// Initialize SDL
	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		return
	}

	// Initialize mixer
	err = mix.Init(mix.INIT_OGG)
	if err != nil {
		return
	}

	// Open audio
	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE)
	if err != nil {
		return
	}

	// Initialize ttf
	err = ttf.Init()
	if err != nil {
		return
	}

	// Initialize controller
	err = sdl.InitSubSystem(sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER | sdl.INIT_HAPTIC)
	if err != nil {
		return
	}

	// Get window dimensions based on display aspect ratio
	width, height := e.GetDimensions()

	// Create window
	e.Window, err = sdl.CreateWindow("VoV", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return
	}

	// Updates window dimensions
	e.UpdateDimensions()

	// Create renderer
	e.Renderer, err = sdl.CreateRenderer(e.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return
	}

	// Set logical size
	err = e.Renderer.SetLogicalSize(int(e.Cfg.WinWidth), int(e.Cfg.WinHeight))
	if err != nil {
		return
	}

	// Event filter callback
	FilterEvent := func(event sdl.Event, userdata interface{}) bool {
		switch t := event.(type) {
		case *sdl.CommonEvent:
			if t.Type == sdl.APP_WILLENTERBACKGROUND {
				if Paused {
					WasPaused = true
				} else {
					Pause()
				}
			}

			if t.Type == sdl.APP_WILLENTERFOREGROUND {
				if WasPaused {
					WasPaused = false
				} else {
					Unpause()
				}
			}
		}

		return true
	}

	// Add event watch function
	sdl.AddEventWatchFunc(FilterEvent, nil)

	// Set android hints
	sdl.SetHint(sdl.HINT_ACCELEROMETER_AS_JOYSTICK, "1")
	sdl.SetHint(sdl.HINT_ANDROID_SEPARATE_MOUSE_AND_TOUCH, "1")

	// Maximum FPS in milliseconds
	e.FrameMs = uint32(1000 / e.Cfg.MaxFps)

	return
}

// Sets window icon
func (e *Engine) SetIcon(icon *sdl.Surface) {
	e.Window.SetIcon(icon)
	icon.Free()
}

// Gets window dimensions based on display aspect ratio
func (e *Engine) GetDimensions() (width, height int) {
	m := &sdl.DisplayMode{}
	sdl.GetDesktopDisplayMode(0, m)
	aspect := float64(m.W) / float64(m.H)

	switch aspect {
	case 1.33:
		// 4:3
		width = 1024
		height = 768
	case 1.50:
		// 3:2
		width = 1024
		height = 682
	case 1.60:
		// 16:10/8:5
		width = 1024
		height = 640
	case 1.66:
		// 5:3
		width = 1024
		height = 614
	case 1.70:
		// 17:10
		width = 1024
		height = 602
	case 1.77:
		// 16:9
		width = 1024
		height = 576
	default:
		width = 1024
		height = 640
	}

	return
}

// Updates window dimensions
func (e *Engine) UpdateDimensions() {
	w, h := e.Window.GetSize()
	e.Cfg.WinWidth, e.Cfg.WinHeight = float64(w), float64(h)

	e.Cfg.XScrollTo = e.Cfg.WinWidth / 3
	e.Cfg.YScrollTo = e.Cfg.WinHeight / 2

	e.Cfg.MaxDistAhead = e.Cfg.WinWidth
}

// Calculates start frame
func (e *Engine) StartFrame() {
	// Get start ticks
	e.StartTicks = GetTicks()

	// All movements are based on TFrame (1/20th of a second)
	e.TFrame = e.Cfg.GameSpeed * float64(e.FrameDelta) / 50
}

// Calculates end frame
func (e *Engine) EndFrame() {
	if !Paused {
		// Increment frame counter
		e.Frames++
	}

	// Get end ticks and calculate delta
	e.EndTicks = GetTicks()
	e.FrameDelta = e.EndTicks - e.StartTicks

	// Cap the frame rate
	if e.FrameMs > e.FrameDelta {
		// Sleep the remaining frame time
		sdl.Delay(e.FrameMs - e.FrameDelta)
	}

	// Calculate FPS
	e.Fps = float64(e.Frames) / (float64(e.EndTicks) / 1000)
}

// Clears screen
func (e *Engine) Clear() {
	e.Renderer.Clear()
	e.Renderer.SetDrawColor(0, 0, 0, 255)
	e.Renderer.FillRect(nil)
}

// Toggles fullscreen
func (e *Engine) Fullscreen() {
	flag := uint32(sdl.WINDOW_FULLSCREEN_DESKTOP)

	if e.Window.GetFlags()&flag != 0 {
		e.Window.SetFullscreen(0)
	} else {
		e.Window.SetFullscreen(flag)
	}
}

// Initializes game controller
func (e *Engine) SetController(mappings []string) {
	// Add controller mappings
	for _, mapping := range mappings {
		sdl.GameControllerAddMapping(mapping)
	}

	// Open game controller
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			e.Controller = sdl.GameControllerOpen(i)
			break
		}
	}
}

// Closes game controller
func (e *Engine) CloseController() {
	if e.Controller != nil {
		e.Controller.Close()
		e.Controller = nil
	}

	if Haptic != nil {
		Haptic.Close()
		Haptic = nil
	}
}

// Sets android accelerometer
func (e *Engine) SetAccelerometer() {
	for i := 0; i < sdl.NumJoysticks(); i++ {
		e.Joystick = sdl.JoystickOpen(sdl.JoystickID(i))

		if e.Joystick.Name() == "Android Accelerometer" || e.Joystick.NumAxes() == 3 {
			break
		} else {
			e.Joystick.Close()
		}
	}
}

// Sets joystick haptic
func (e *Engine) SetHaptic() {
	if e.Controller != nil && sdl.JoystickIsHaptic(e.Controller.GetJoystick()) == 1 {
		Haptic = sdl.HapticOpenFromJoystick(e.Controller.GetJoystick())
	} else if e.Joystick != nil && sdl.JoystickIsHaptic(e.Joystick) == 1 {
		Haptic = sdl.HapticOpenFromJoystick(e.Joystick)
	}

	if Haptic != nil && Haptic.RumbleSupported() != 0 {
		Haptic.RumbleInit()
	}
}

// Destroys SDL and releases the memory
func (e *Engine) Destroy() {
	e.Renderer.Destroy()
	e.Window.Destroy()

	ttf.Quit()

	mix.HaltMusic()
	mix.CloseAudio()
	mix.Quit()

	e.CloseController()
	if e.Joystick != nil {
		e.Joystick.Close()
	}

	sdl.Quit()
}

// Quits game loop
func (e *Engine) Quit() {
	e.Running = false
	StopTimer()
}
