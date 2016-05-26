// VoV game
package game

import (
	"fmt"
	"math"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/system/rumble"
)

// Game states
const (
	GamePlay = iota
	GameOver
	GamePause
	GameQuit
	DeadPause
)

// Directions
const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

// Input directions
type Direction struct {
	X1     float32
	Y1     float32
	X2     float32
	Y2     float32
	Motion bool
	State  []bool
}

// Sets all direction states to state
func (d *Direction) SetStates(state bool) {
	for i := 0; i < len(d.State); i++ {
		d.State[i] = state
	}
}

// Checks if any of direction states is enabled
func (d *Direction) StateEnabled() bool {
	for i := 0; i < len(d.State); i++ {
		if d.State[i] {
			return true
		}
	}

	return false
}

// Game structure
type Game struct {
	Engine   *engine.Engine
	Cfg      *engine.Config
	Resource *engine.Resource

	Life *Sprite

	FpsText  *Sprite
	TimeText *Sprite

	ShieldsText     *Sprite
	AttackText      *Sprite
	InvincibleText  *Sprite
	EngineBlastText *Sprite
	SlowdownText    *Sprite

	PausedText   *Sprite
	GameOverText *Sprite

	Fog    *Fog
	Dust   *Dust
	Dots   *Dots
	Ship   *Ship
	Rocks  *Rocks
	Powups *Powups

	State        int
	LastState    int
	StateTimeout float64

	Direction *Direction

	Score int
}

// Returns new game
func NewGame(e *engine.Engine, r *engine.Resource) (g *Game) {
	g = &Game{}
	g.Engine = e
	g.Cfg = e.Cfg
	g.Resource = r

	g.Direction = &Direction{}

	g.Fog = NewFog(e, r)
	g.Dust = NewDust(e)
	g.Dots = NewDots(g)
	g.Ship = NewShip(g)
	g.Powups = NewPowups(g)
	g.Rocks = NewRocks(e, r)

	return
}

// Initializes game state
func (g *Game) OnInit() bool {
	g.Direction.State = make([]bool, 4)

	// Initialize ship
	g.Ship.Init()

	g.Ship.DX = g.Engine.ScreenDX
	g.Ship.DY = g.Engine.ScreenDY

	g.Engine.ScreenDX = g.Cfg.BarrierSpeed
	g.Engine.ScreenDY = 0.0

	// Create sprites
	g.Life = NewSprite(g.Engine, g.Resource.Life)

	g.FpsText = NewSprite(g.Engine, g.Resource.FpsText)
	g.TimeText = NewSprite(g.Engine, g.Resource.TimeText)

	g.ShieldsText = NewSprite(g.Engine, g.Resource.ShieldsText)
	g.AttackText = NewSprite(g.Engine, g.Resource.AttackText)
	g.InvincibleText = NewSprite(g.Engine, g.Resource.InvincibleText)
	g.EngineBlastText = NewSprite(g.Engine, g.Resource.EngineBlastText)
	g.SlowdownText = NewSprite(g.Engine, g.Resource.SlowdownText)

	g.PausedText = NewSprite(g.Engine, g.Resource.PausedText)
	g.GameOverText = NewSprite(g.Engine, g.Resource.GameOverText)

	// Sprites positions
	g.FpsText.X = g.Cfg.WinWidth - g.FpsText.Width - 300
	g.FpsText.Y = 10

	g.TimeText.X = g.Cfg.WinWidth - g.TimeText.Width - 150
	g.TimeText.Y = 10

	g.ShieldsText.X = 200
	g.ShieldsText.Y = 10
	g.AttackText.X = 200
	g.AttackText.Y = 10
	g.InvincibleText.X = 200
	g.InvincibleText.Y = 10
	g.EngineBlastText.X = 200
	g.EngineBlastText.Y = 10
	g.SlowdownText.X = 200
	g.SlowdownText.Y = 10

	g.PausedText.X = g.Cfg.WinWidth/2 - (g.PausedText.Width / 2)
	g.PausedText.Y = g.Cfg.WinHeight/2 - (g.PausedText.Height / 2)

	g.GameOverText.X = g.Cfg.WinWidth/2 - (g.GameOverText.Width / 2)
	g.GameOverText.Y = g.Cfg.WinHeight/2 - (g.GameOverText.Height / 2)

	// Initialize objects
	g.Fog.Init()
	g.Dust.Init()
	g.Dots.Init()
	g.Rocks.Init()
	g.Powups.Init()

	// Play game music
	g.Resource.PlayMusic(g.Resource.MusicGame, -1)

	return true
}

// Quits game state
func (g *Game) OnQuit() bool {
	mix.HaltMusic()
	return true
}

// Returns state string
func (g *Game) String() string {
	return "Game"
}

// Handles input events
func (g *Game) HandleEvents() {
	//if engine.Paused {
	//event := sdl.WaitEvent()
	//if event != nil {
	//g.HandleEvent(event)
	//}
	//} else {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		g.HandleEvent(event)
	}
	//}
}

// Handles input event
func (g *Game) HandleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		// Handle quit event
		g.Engine.Quit()

	case *sdl.KeyDownEvent:
		// Handle pause
		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK ||
			t.Keysym.Scancode == sdl.SCANCODE_P || t.Keysym.Scancode == sdl.SCANCODE_PAUSE ||
			t.Keysym.Scancode == sdl.SCANCODE_SPACE {

			if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
				if engine.Paused {
					engine.Unpause()
					g.Engine.State.Change(NewMenu(g.Engine, g.Resource))
					break
				}
			}

			g.TogglePause()
		}

		// Handle fullscreen
		if (t.Keysym.Mod&sdl.KMOD_ALT != 0 && t.Keysym.Scancode == sdl.SCANCODE_RETURN) || t.Keysym.Scancode == sdl.SCANCODE_F11 {
			g.Engine.Fullscreen()
			break
		}

		if g.State != GamePlay {
			break
		}

		// Handle keys directions
		switch t.Keysym.Scancode {
		case sdl.SCANCODE_LEFT:
			g.Ship.PlaySound()
			g.Direction.State[LEFT] = true
		case sdl.SCANCODE_RIGHT:
			g.Ship.PlaySound()
			g.Direction.State[RIGHT] = true
		case sdl.SCANCODE_UP:
			g.Ship.PlaySound()
			g.Direction.State[UP] = true
		case sdl.SCANCODE_DOWN:
			g.Ship.PlaySound()
			g.Direction.State[DOWN] = true
		}
	case *sdl.KeyUpEvent:
		g.Ship.FadeSound()

		// Handle keys directions
		switch t.Keysym.Scancode {
		case sdl.SCANCODE_LEFT:
			g.Direction.State[LEFT] = false
		case sdl.SCANCODE_RIGHT:
			g.Direction.State[RIGHT] = false
		case sdl.SCANCODE_UP:
			g.Direction.State[UP] = false
		case sdl.SCANCODE_DOWN:
			g.Direction.State[DOWN] = false
		}

	case *sdl.TouchFingerEvent:
		if engine.Paused {
			g.TogglePause()
			break
		}

		if g.State != GamePlay {
			break
		}

		// Handle touch directions
		if t.Type == sdl.FINGERDOWN {
			g.Direction.X1 = t.X
			g.Direction.Y1 = t.Y
		} else if t.Type == sdl.FINGERMOTION {
			g.Direction.X2 = t.X
			g.Direction.Y2 = t.Y

			g.Direction.Motion = true

			xDiff := g.Direction.X2 - g.Direction.X1
			yDiff := g.Direction.Y2 - g.Direction.Y1

			if math.Abs(float64(xDiff)) > math.Abs(float64(yDiff)) && math.Abs(float64(xDiff)) > g.Engine.Cfg.TouchThreshold {
				g.Ship.PlaySound()

				if g.Direction.X1 < g.Direction.X2 {
					g.Direction.State[RIGHT] = true
				} else {
					g.Direction.State[LEFT] = true
				}
			} else if math.Abs(float64(yDiff)) > math.Abs(float64(xDiff)) && math.Abs(float64(yDiff)) > g.Engine.Cfg.TouchThreshold {
				g.Ship.PlaySound()

				if g.Direction.Y1 > g.Direction.Y2 {
					g.Direction.State[UP] = true
				} else {
					g.Direction.State[DOWN] = true
				}
			}
		} else if t.Type == sdl.FINGERUP {
			g.Ship.FadeSound()
			g.Direction.Motion = false
			g.Direction.SetStates(false)
		}

	case *sdl.MouseMotionEvent:
		// Handle mouse directions
		if g.Direction.Motion {
			g.Direction.X2 = float32(t.X)
			g.Direction.Y2 = float32(t.Y)

			xDiff := g.Direction.X2 - g.Direction.X1
			yDiff := g.Direction.Y2 - g.Direction.Y1

			if math.Abs(float64(xDiff)) > math.Abs(float64(yDiff)) {
				g.Ship.PlaySound()

				if g.Direction.X1 < g.Direction.X2 {
					g.Direction.State[RIGHT] = true
				} else {
					g.Direction.State[LEFT] = true
				}
			} else {
				g.Ship.PlaySound()

				if g.Direction.Y1 > g.Direction.Y2 {
					g.Direction.State[UP] = true
				} else {
					g.Direction.State[DOWN] = true
				}
			}
		}

	case *sdl.MouseButtonEvent:
		if engine.Paused {
			g.TogglePause()
			break
		}

		if g.State != GamePlay {
			break
		}

		// Handle mouse directions
		if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
			g.Direction.X1 = float32(t.X)
			g.Direction.Y1 = float32(t.Y)

			g.Direction.Motion = true
		} else if t.Type == sdl.MOUSEBUTTONUP && t.Button == sdl.BUTTON_LEFT {
			g.Ship.FadeSound()
			g.Direction.Motion = false
			g.Direction.SetStates(false)
		}

	case *sdl.ControllerDeviceEvent:
		// Initialize/Remove controller
		if t.Type == sdl.CONTROLLERDEVICEADDED {
			g.Engine.Controller = sdl.GameControllerOpen(int(t.Which))
			if g.Engine.Cfg.HapticEnabled {
				g.Engine.SetHaptic()
			}
		} else if t.Type == sdl.CONTROLLERDEVICEREMOVED {
			g.Engine.CloseController()
		}

	case *sdl.ControllerButtonEvent:
		// Handle pause
		if t.Type == sdl.CONTROLLERBUTTONDOWN {
			if t.Button == sdl.CONTROLLER_BUTTON_A || t.Button == sdl.CONTROLLER_BUTTON_START || t.Button == sdl.CONTROLLER_BUTTON_B {
				if engine.Paused {
					g.TogglePause()
					break
				}

				g.TogglePause()
			}

			if t.Button == sdl.CONTROLLER_BUTTON_BACK {
				if engine.Paused {
					engine.Unpause()
				}

				g.Engine.State.Change(NewMenu(g.Engine, g.Resource))
			}
		}

	case *sdl.ControllerAxisEvent:
		if g.State != GamePlay {
			break
		}

		// Handle joystick directions
		switch t.Axis {
		case sdl.CONTROLLER_AXIS_LEFTX:
			if t.Value < 0 {
				g.Ship.PlaySound()
				g.Direction.State[LEFT] = true
			} else if t.Value > 0 {
				g.Ship.PlaySound()
				g.Direction.State[RIGHT] = true
			} else if t.Value == 0 {
				g.Ship.FadeSound()
				g.Direction.State[LEFT] = false
				g.Direction.State[RIGHT] = false
			}

		case sdl.CONTROLLER_AXIS_LEFTY:
			if t.Value < 0 {
				g.Ship.PlaySound()
				g.Direction.State[UP] = true
			} else if t.Value > 0 {
				g.Ship.PlaySound()
				g.Direction.State[DOWN] = true
			} else if t.Value == 0 {
				g.Ship.FadeSound()
				g.Direction.State[UP] = false
				g.Direction.State[DOWN] = false
			}
		}

	case *sdl.JoyAxisEvent:
		if runtime.GOOS != "android" || !g.Engine.Cfg.AccelerometerEnabled {
			break
		}

		if g.State != GamePlay || g.Direction.Motion {
			break
		}

		// TODO
		// Handle accelerometer directions
		switch t.Axis {
		case 0:
			xDiff := g.Direction.X1 - float32(t.Value)
			g.Direction.X1 = float32(t.Value)

			if t.Value < 0 && math.Abs(float64(xDiff)) > g.Engine.Cfg.AccelThreshold {
				g.Ship.PlaySound()
				g.Direction.State[LEFT] = true
			} else if t.Value > 0 && math.Abs(float64(xDiff)) > g.Engine.Cfg.AccelThreshold {
				g.Ship.PlaySound()
				g.Direction.State[RIGHT] = true
			} else {
				g.Ship.FadeSound()
				g.Direction.State[LEFT] = false
				g.Direction.State[RIGHT] = false
			}

		case 1:
			yDiff := g.Direction.Y1 - float32(t.Value)
			g.Direction.Y1 = float32(t.Value)

			if t.Value < 0 && math.Abs(float64(yDiff)) > g.Engine.Cfg.AccelThreshold {
				g.Ship.PlaySound()
				g.Direction.State[UP] = true
			} else if t.Value > 0 && math.Abs(float64(yDiff)) > g.Engine.Cfg.AccelThreshold {
				g.Ship.PlaySound()
				g.Direction.State[DOWN] = true
			} else {
				g.Ship.FadeSound()
				g.Direction.State[UP] = false
				g.Direction.State[DOWN] = false
			}
		}

	default:
		break
	}
}

// Toggles paused state
func (g *Game) TogglePause() {
	if g.State == GameOver {
		return
	}

	if !engine.Paused {
		g.Ship.FadeSound()
		g.Resource.PlaySound(g.Resource.SoundClick, -1, 0)

		g.Direction.SetStates(false)
		engine.Pause()

		g.LastState = g.State
		g.State = GamePause
	} else {
		g.Resource.PlaySound(g.Resource.SoundClick, -1, 0)
		engine.Unpause()
		g.State = g.LastState
	}
}

// Updates game state
func (g *Game) UpdateState() {
	if g.StateTimeout > 0 {
		g.StateTimeout -= g.Engine.TFrame * 3
		return
	}

	switch g.State {
	case GamePlay:
		// Update ship direction
		if g.Direction.State[LEFT] {
			g.Ship.DX -= g.Cfg.ThrusterStrength * g.Engine.TFrame
			g.Ship.Jets |= 1 << 0
		}
		if g.Direction.State[DOWN] {
			g.Ship.DY += g.Cfg.ThrusterStrength * g.Engine.TFrame
			g.Ship.Jets |= 1 << 1
		}
		if g.Direction.State[RIGHT] {
			g.Ship.DX += g.Cfg.ThrusterStrength * g.Engine.TFrame
			g.Ship.Jets |= 1 << 2
		}
		if g.Direction.State[UP] {
			g.Ship.DY -= g.Cfg.ThrusterStrength * g.Engine.TFrame
			g.Ship.Jets |= 1 << 3
		}

		if g.Ship.Jets != 0 {
			g.Ship.DX = fconstrain2(g.Ship.DX, -50, 50)
			g.Ship.DY = fconstrain2(g.Ship.DY, -50, 50)
		}

	case GameOver:
		// Stop rumble
		if g.Engine.Cfg.HapticEnabled {
			rumble.RumbleStop()
		}

		// Change state
		g.State = GameQuit

	case DeadPause:
		// Stop rumble
		if g.Engine.Cfg.HapticEnabled {
			rumble.RumbleStop()
		}

		// Restore the ship
		g.Ship.Flags = DRAW | MOVE | COLLIDE
		g.State = GamePlay

		// Change ship state
		g.Ship.State = INVINCIBLE
		g.Ship.StateTimeout = g.Engine.Cfg.InvinciblePauseLength
	}
}

// Updates game
func (g *Game) Update() {
	// Don't update if timer is paused
	if engine.Paused {
		return
	}

	// Update state
	g.UpdateState()

	// Update objects
	g.Ship.Update()
	g.Rocks.Update()
	g.Powups.Update()
	g.Dots.Update()
	g.Dust.Update()
	g.Fog.Update()

	// Generate rocks
	g.Rocks.New()

	// Generate powups
	g.Powups.New()

	// Rocks collisions
	g.Rocks.Collisions()

	// Powups collisions
	g.Powups.Collisions()

	if g.State != GameOver {
		// Update score
		g.Score += int(g.Engine.FrameDelta)

		// Ship engine dots
		g.Dots.NewShipDots()

		// Ship collisions
		g.Ship.Collisions()
	}
}

// Draws fps
func (g *Game) DrawFPS() {
	g.FpsText.Draw()
	fps := fmt.Sprintf("%.1f", g.Engine.Fps)
	x := int32(g.FpsText.X + (g.FpsText.Width))
	y := int32(g.FpsText.Y)
	g.Resource.DrawText(fps, x, y, engine.FONT_SMALL)
}

// Draws score
func (g *Game) DrawScore() {
	g.TimeText.Draw()
	x := int32(g.TimeText.X + (g.TimeText.Width))
	y := int32(g.TimeText.Y)
	score := formatTime(g.Score, false)
	g.Resource.DrawText(score, x, y, engine.FONT_SMALL)
}

// Draws lives
func (g *Game) DrawLives() {
	for i := 0; i < g.Ship.Lives-1; i++ {
		g.Life.X = float64(i+1) * (g.Life.Width * 1.5)
		g.Life.Y = g.Life.Width / 2
		g.Life.Draw()
	}
}

// Draws ship state timeout
func (g *Game) DrawState() {
	var x, y int32

	switch g.Ship.State {
	case SHIELDS:
		g.ShieldsText.Draw()
		x = int32(g.ShieldsText.X + (g.ShieldsText.Width))
		y = int32(g.ShieldsText.Y)

	case ATTACK:
		g.AttackText.Draw()
		x = int32(g.AttackText.X + (g.AttackText.Width))
		y = int32(g.AttackText.Y)

	case INVINCIBLE:
		g.InvincibleText.Draw()
		x = int32(g.InvincibleText.X + (g.InvincibleText.Width))
		y = int32(g.InvincibleText.Y)

	case ENGINEBLAST:
		g.EngineBlastText.Draw()
		x = int32(g.EngineBlastText.X + (g.EngineBlastText.Width))
		y = int32(g.EngineBlastText.Y)

	case SLOWDOWN:
		g.SlowdownText.Draw()
		x = int32(g.SlowdownText.X + (g.SlowdownText.Width))
		y = int32(g.SlowdownText.Y)
	}

	timeout := formatTime(int(g.Ship.StateTimeout), false)
	if g.Ship.StateTimeout > 1000 {
		g.Resource.DrawText(timeout, x, y, engine.FONT_SMALL)
	} else {
		g.Resource.DrawText(timeout, x, y, engine.FONT_SMALL_RED)
	}
}

// Draws game
func (g *Game) Draw() {
	// Draw objects
	g.Dust.Draw()
	g.Fog.Draw()
	g.Dots.Draw()
	g.Rocks.Draw()
	g.Powups.Draw()
	g.Ship.Draw()

	// Draw score
	g.DrawScore()

	// Draw lives
	g.DrawLives()

	// Draw ship state
	if g.Ship.State != PLAIN && g.State != GameOver {
		g.DrawState()
	}

	// Draw FPS
	if g.Cfg.ShowFps && g.State != GamePause {
		g.DrawFPS()
	}

	// Draw paused text
	if g.State == GamePause {
		g.PausedText.Draw()
	}

	// Draw gameover text
	if g.State == GameOver {
		g.GameOverText.Draw()
	}

	if g.State == GameQuit {
		// Change state to menu
		g.Engine.State.Change(NewScores(g.Engine, g.Resource, g.Score, true))
	}
}
