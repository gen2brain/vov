// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/system/rumble"
)

// Ship states
const (
	PLAIN = iota
	INVINCIBLE
	ENGINEBLAST
	SHIELDS
	ATTACK
	SLOWDOWN
)

// Ship structure
type Ship struct {
	Sprite

	Game *Game
	Cfg  *engine.Config

	Glow *Sprite

	Moving      bool
	Transparent bool

	StateTimeout  float64
	TranspTimeout float64

	LifePowText        *Sprite
	ShieldsPowText     *Sprite
	AttackPowText      *Sprite
	InvinciblePowText  *Sprite
	EngineBlastPowText *Sprite
	SlowdownPowText    *Sprite

	PowupTextScale   float64
	PowupTextTimeout float64
	PowupCurrent     int
	PowupTextActive  bool
}

// Returns new ship
func NewShip(g *Game) *Ship {
	s := &Ship{}
	s.Game = g
	s.Engine = g.Engine
	s.Cfg = g.Engine.Cfg

	return s
}

// Initialize ship
func (s *Ship) Init() {
	s.Type = SHIP
	s.Lives = 4
	s.Flags = MOVE | DRAW | COLLIDE
	s.State = PLAIN

	s.Texture = s.Game.Resource.Ship
	s.Query()

	s.Surface = s.Game.Resource.ShipSurf

	s.X = s.Cfg.WinWidth / 2.2
	s.Y = s.Cfg.WinHeight/2 - s.Width/2
	s.DX = s.Cfg.BarrierSpeed
	s.DY = 0.0

	s.Exp1 = NewSprite(s.Game.Engine, s.Game.Resource.Explosion1)

	s.Glow = NewSprite(s.Game.Engine, s.Game.Resource.ShipGlow)
	s.Glow.Texture.SetBlendMode(sdl.BLENDMODE_ADD)

	s.LifePowText = NewSprite(s.Game.Engine, s.Game.Resource.LifePowText)
	s.ShieldsPowText = NewSprite(s.Game.Engine, s.Game.Resource.ShieldsPowText)
	s.AttackPowText = NewSprite(s.Game.Engine, s.Game.Resource.AttackPowText)
	s.InvinciblePowText = NewSprite(s.Game.Engine, s.Game.Resource.InvinciblePowText)
	s.EngineBlastPowText = NewSprite(s.Game.Engine, s.Game.Resource.EngineBlastPowText)
	s.SlowdownPowText = NewSprite(s.Game.Engine, s.Game.Resource.SlowdownPowText)

	s.Width /= float64(6)
}

// Kills ship
func (s *Ship) Kill() {
	s.ExpActive = true

	// Play explosion sound
	s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion1, 2, 0)

	// Play rumble
	if s.Cfg.HapticEnabled {
		rumble.RumblePlay(1.0, 1000)
	}

	// Take life
	s.Lives -= 1

	s.FadeSound()
	s.Game.Direction.Motion = false
	s.Game.Direction.SetStates(false)
	s.PowupTextActive = false

	s.Cfg.GameSpeed = 1.0
	s.Cfg.EngineDots = 1000

	if s.Lives == 0 {
		s.Game.State = GameOver
		s.Game.StateTimeout = s.Cfg.GameOverLength

		s.Flags = 0

		// Scrolling is based on the ship speed, so we need to reset it
		s.DX = s.Cfg.BarrierSpeed
		s.DY = 0
	} else {
		s.Game.State = DeadPause
		s.Game.StateTimeout = s.Cfg.DeadPauseLength

		// Want ship to be invisible, but keep drifting at sqrt(speed)
		// to leave it in the middle of the space from the explosion
		s.Flags = MOVE

		if s.DX < 0 {
			s.DX = -math.Sqrt(-s.DX)
		} else {
			s.DX = math.Sqrt(s.DX)
		}

		if s.DY < 0 {
			s.DY = -math.Sqrt(-s.DY)
		} else {
			s.DY = math.Sqrt(s.DY)
		}

		if s.DX < s.Cfg.BarrierSpeed {
			s.DX = s.Cfg.BarrierSpeed
		}
	}

	s.Game.Dots.NewShipBangDots()
}

// Plays engine sound
func (s *Ship) PlaySound() {
	s.Moving = true
	if mix.Playing(1) == 0 {
		if s.State == ENGINEBLAST {
			s.Game.Resource.PlaySound(s.Game.Resource.SoundEngine2, 1, 0)
		} else {
			s.Game.Resource.PlaySound(s.Game.Resource.SoundEngine1, 1, 0)
		}
	}
}

// Fades engine sound
func (s *Ship) FadeSound() {
	s.Moving = false
	if mix.Playing(1) == 1 {
		mix.FadeOutChannel(1, 100)
	}
}

// Checks ship collisions
func (s *Ship) Collisions() {
	for i := 0; i < s.Cfg.MaxPowups; i++ {
		if s.Game.Powups.Powups[i] != nil && s.Game.Powups.Powups[i].Active {

			if s.Collide(s.Game.Powups.Powups[i]) {
				s.Game.Powups.Powups[i].Active = false

				// Restore default config
				if s.State == SLOWDOWN && s.Game.Powups.Powups[i].State != PLAIN {
					s.Game.Engine.Cfg.GameSpeed = 1.0
				}
				if s.State == ENGINEBLAST && s.Game.Powups.Powups[i].State != PLAIN {
					s.Game.Engine.Cfg.EngineDots = 1000
				}

				switch s.Game.Powups.Powups[i].State {
				case PLAIN:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup0, -1, 0)
				case INVINCIBLE:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup1, -1, 0)
				case ENGINEBLAST:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup2, -1, 0)
					s.Game.Engine.Cfg.EngineDots = 1500
				case SHIELDS:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup3, -1, 0)
				case ATTACK:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup4, -1, 0)
				case SLOWDOWN:
					s.Game.Resource.PlaySound(s.Game.Resource.SoundPowup5, -1, 0)
					s.Game.Engine.Cfg.GameSpeed = 0.50
				}

				s.PowupTextActive = true
				s.PowupTextScale = s.Game.Engine.Cfg.PowupTextScale
				s.PowupTextTimeout = s.Game.Engine.Cfg.PowupTextTimeout
				s.PowupCurrent = s.Game.Powups.Powups[i].State

				if s.Game.Powups.Powups[i].State != PLAIN {
					if s.State == s.Game.Powups.Powups[i].State {
						s.StateTimeout += s.Game.Engine.Cfg.PowupStateTimeout
					} else {
						s.StateTimeout = s.Game.Engine.Cfg.PowupStateTimeout
					}

					s.State = s.Game.Powups.Powups[i].State
				} else {
					// Extra life
					s.Lives += 1
				}
			}
		}
	}

	for i := 0; i < s.Cfg.MaxRocks; i++ {
		if s.Game.Rocks.Rocks[i] != nil && s.Game.Rocks.Rocks[i].Active {

			if s.Collide(s.Game.Rocks.Rocks[i]) {
				switch s.State {
				case PLAIN, SLOWDOWN:
					if mix.Playing(2) == 0 {
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					} else {
						mix.FadeOutChannel(2, 10)
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					}

					// Kill ship
					s.Kill()

					// Kill rock
					s.Game.Rocks.Rocks[i].Kill()

				case SHIELDS:
					if mix.Playing(2) == 0 {
						s.Game.Resource.PlaySound(s.Game.Resource.SoundBounce, 2, 0)
					} else {
						mix.FadeOutChannel(2, 10)
						s.Game.Resource.PlaySound(s.Game.Resource.SoundBounce, 2, 0)
					}

					// Bounce ship
					s.Bounce(s.Game.Rocks.Rocks[i])

				case ATTACK:
					if mix.Playing(2) == 0 {
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					} else {
						mix.FadeOutChannel(2, 10)
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					}

					// Bounce ship
					s.Bounce(s.Game.Rocks.Rocks[i])

					// Kill rock
					s.Game.Rocks.Rocks[i].Kill()

					// New bang dots
					s.Game.Dots.NewBangDots(s.Game.Rocks.Rocks[i])

				case INVINCIBLE:
					// Set alpha transparency
					s.Texture.SetAlphaMod(100)
					s.Transparent = true
					s.TranspTimeout = 200

					s.Game.Resource.PlaySoundTimed(s.Game.Resource.SoundEngine3, 2, 0, 200)

				case ENGINEBLAST:
					if mix.Playing(2) == 0 {
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					} else {
						mix.FadeOutChannel(2, 10)
						s.Game.Resource.PlaySound(s.Game.Resource.SoundExplosion2, 2, 0)
					}

					// Kill ship
					s.Kill()

					// Kill rock
					s.Game.Rocks.Rocks[i].Kill()
				}
			}
		}
	}
}

// Updates ship state
func (s *Ship) UpdateState() {
	if s.StateTimeout > 0 {
		s.StateTimeout -= float64(s.Game.Engine.FrameDelta)
		return
	}

	if s.State != PLAIN {
		// Restore default config
		if s.State == SLOWDOWN {
			s.Game.Engine.Cfg.GameSpeed = 1.0
		}
		if s.State == ENGINEBLAST {
			s.Game.Engine.Cfg.EngineDots = 1000
		}

		// Restore default state
		s.State = PLAIN
	}
}

// Updates ship transparency
func (s *Ship) UpdateTransp() {
	if s.TranspTimeout > 0 {
		s.TranspTimeout -= float64(s.Game.Engine.FrameDelta)
		return
	}

	if s.Transparent {
		s.Texture.SetAlphaMod(255)
		s.Transparent = false
	}
}

// Updates ship
func (s *Ship) Update() {
	// Update ship state
	s.UpdateState()

	// Update transparency
	s.UpdateTransp()

	// Scrolling
	tmp := (s.Y+s.Height/2+s.DY*s.Game.Engine.TFrame-s.Cfg.YScrollTo)/25 + (s.DY - s.Game.Engine.ScreenDY)
	s.Game.Engine.ScreenDY += tmp * s.Game.Engine.TFrame / 12
	tmp = (s.X+s.Width/2+s.DX*s.Game.Engine.TFrame-s.Cfg.XScrollTo)/25 + (s.DX - s.Game.Engine.ScreenDX)
	s.Game.Engine.ScreenDX += tmp * s.Game.Engine.TFrame / 12

	// Taper off so we don't hit the barrier abruptly. If we would hit in < 2 seconds, adjust to 2 seconds
	if s.Cfg.DistAhead+(s.Game.Engine.ScreenDX-s.Cfg.BarrierSpeed)*toTicks(2, s.Cfg.GameSpeed) < 0 {
		s.Game.Engine.ScreenDX = s.Cfg.BarrierSpeed - (s.Cfg.DistAhead / toTicks(2, s.Cfg.GameSpeed))
	}
	s.Cfg.DistAhead += (s.Game.Engine.ScreenDX - s.Cfg.BarrierSpeed) * s.Game.Engine.TFrame
	if s.Cfg.MaxDistAhead >= 0 {
		s.Cfg.DistAhead = math.Min(s.Cfg.DistAhead, s.Cfg.MaxDistAhead)
	}

	// Move sprite
	s.Sprite.Update()

	// Bounce off left or right edge of screen
	if s.X < 0 || s.X+s.Width > s.Cfg.WinWidth {
		s.X -= (s.DX - s.Game.Engine.ScreenDX) * s.Game.Engine.TFrame
		s.DX = s.Game.Engine.ScreenDX - (s.DX-s.Game.Engine.ScreenDX)*s.Cfg.Bounciness
		s.X = fconstrain(s.X, s.Cfg.WinWidth-s.Width)
	}

	// Bounce off top or bottom of screen
	if s.Y < 0 || s.Y+s.Height > s.Cfg.WinHeight {
		s.Y -= (s.DY - s.Game.Engine.ScreenDY) * s.Game.Engine.TFrame
		s.DY = s.Game.Engine.ScreenDY - (s.DY-s.Game.Engine.ScreenDY)*s.Cfg.Bounciness
		s.Y = fconstrain(s.Y, s.Cfg.WinHeight-s.Height)
	}
}

// Draws powup text
func (s *Ship) DrawPowupText() {
	if s.PowupTextActive && s.Game.State == GamePlay {
		var spr *Sprite

		switch s.PowupCurrent {
		case PLAIN:
			spr = s.LifePowText
		case INVINCIBLE:
			spr = s.InvinciblePowText
		case ENGINEBLAST:
			spr = s.EngineBlastPowText
		case SHIELDS:
			spr = s.ShieldsPowText
		case ATTACK:
			spr = s.AttackPowText
		case SLOWDOWN:
			spr = s.SlowdownPowText
		}

		ratio := spr.Width / spr.Height
		width := spr.Width / s.PowupTextScale
		height := width / ratio

		src := spr.Rect()

		x := s.Cfg.WinWidth/2 - (width / 2)
		y := s.Cfg.WinHeight/2 - (height / 2)

		dest := &sdl.Rect{int32(x), int32(y), int32(width), int32(height)}

		s.Game.Engine.Renderer.Copy(spr.Texture, src, dest)

		if s.PowupTextScale <= 1.0 {
			s.PowupTextTimeout -= float64(s.Game.Engine.FrameDelta)
			if s.PowupTextTimeout <= 255 {
				spr.Texture.SetAlphaMod(uint8(s.PowupTextTimeout))
			}
			if s.PowupTextTimeout <= 0 {
				s.PowupTextActive = false
				spr.Texture.SetAlphaMod(255)
			}
		} else {
			s.PowupTextScale -= 0.03
		}
	}
}

// Draws ship
func (s *Ship) Draw() {
	s.Sprite.Draw()

	// Draw glow
	if s.Moving {
		s.Glow.X = s.X - 4
		s.Glow.Y = s.Y - 4
		s.Glow.Draw()
	}

	// Blink if state is to expire
	if s.StateTimeout > 0 && s.StateTimeout < 1000 && !s.Transparent {
		// Set alpha transparency
		s.Texture.SetAlphaMod(80)
		s.Transparent = true
		s.TranspTimeout = 100
	}

	// Draw powup text
	s.DrawPowupText()

	// Draw explosion
	s.Explode()

	// Reset jets
	s.Jets = 0
}
