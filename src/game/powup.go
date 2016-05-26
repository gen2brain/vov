// VoV game
package game

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
)

// Powups structure
type Powups struct {
	Game   *Game
	Engine *engine.Engine

	Powups []*Sprite

	Glow *Sprite

	Timeout int

	SpeedMin [4]float64
	SpeedMax [4]float64
}

// Returns new powups
func NewPowups(g *Game) (r *Powups) {
	r = &Powups{}
	r.Game = g
	r.Engine = g.Engine
	return
}

// Initializes powups
func (r *Powups) Init() {
	r.Powups = make([]*Sprite, r.Engine.Cfg.MaxPowups)

	r.Glow = NewSprite(r.Game.Engine, r.Game.Resource.PowupGlow)
	r.Glow.Texture.SetBlendMode(sdl.BLENDMODE_ADD)
	r.Glow.Width /= float64(r.Engine.Cfg.NFrames)

	for i := 0; i < r.Engine.Cfg.MaxPowups; i++ {
		pow := NewSprite(r.Engine, r.Game.Resource.Powup)
		pow.Width /= float64(6)

		pow.Surface = r.Game.Resource.PowupSurf

		pow.Type = POWUP
		pow.Flags = MOVE | DRAW | COLLIDE
		pow.State = rnd(0, 6)

		pow.Exp1 = NewSprite(r.Engine, r.Game.Resource.Explosion2)

		r.Powups[i] = pow
	}
}

// Computes the speed ranges of powups coming from each side
func (r *Powups) Sides() {
	var dx0, dx1, dy0, dy1 float64

	for i := 0; i < 4; i++ {
		r.SpeedMin[i] = 0
		r.SpeedMax[i] = 0
	}

	dx0 = -r.Engine.Cfg.RDX - r.Engine.ScreenDX
	dx1 = r.Engine.Cfg.RDX - r.Engine.ScreenDX
	dy0 = -r.Engine.Cfg.RDY - r.Engine.ScreenDY
	dy1 = r.Engine.Cfg.RDY - r.Engine.ScreenDY

	if dx0 < 0 {
		r.SpeedMax[RIGHT] = -dx0
		if dx1 < 0 {
			// Powups moving left only. So the RIGHT side of the screen
			r.SpeedMin[RIGHT] = -dx1
		} else {
			// Powups moving left and right
			r.SpeedMax[LEFT] = dx1
		}
	} else {
		// Powups moving right only. So the LEFT side of the screen
		r.SpeedMin[LEFT] = dx0
		r.SpeedMax[LEFT] = dx1
	}

	if dy0 < 0 {
		r.SpeedMax[DOWN] = -dy0
		if dy1 < 0 {
			// Powups moving up only. So the BOTTOM of the screen
			r.SpeedMin[DOWN] = -dy1
		} else {
			// Powups moving up and down
			r.SpeedMax[UP] = dy1
		}
	} else {
		// Powups moving down only. so the TOP of the screen
		r.SpeedMin[UP] = dy0
		r.SpeedMax[UP] = dy1
	}
}

// Generates new powup
func (r *Powups) New() {
	if r.Timeout < r.Engine.Cfg.PowupsTimeout {
		return
	}

	i := urnd() % uint32(r.Engine.Cfg.MaxPowups)

	if r.Powups[i].Active {
		return
	}

	r.Sides()

	direction := rnd(0, 5)
	r.Powups[i].X = 0
	r.Powups[i].Y = 0

	switch direction {
	case RIGHT:
		r.Powups[i].X = r.Engine.Cfg.WinWidth
		r.Powups[i].Y = frnd() * (r.Engine.Cfg.WinHeight + r.Powups[i].Height)

		r.Powups[i].DX = -weightedRndRange(r.SpeedMin[direction], r.SpeedMax[direction]) + r.Engine.ScreenDX
		r.Powups[i].DY = r.Engine.Cfg.RDY * crnd()
	case LEFT:
		r.Powups[i].X = -r.Powups[i].Width
		r.Powups[i].Y = frnd() * (r.Engine.Cfg.WinHeight + r.Powups[i].Height)

		r.Powups[i].DX = weightedRndRange(r.SpeedMin[direction], r.SpeedMax[direction]) + r.Engine.ScreenDX
		r.Powups[i].DY = r.Engine.Cfg.RDY * crnd()
	case DOWN:
		r.Powups[i].X = (frnd() * (r.Engine.Cfg.WinWidth + r.Powups[i].Width)) - r.Powups[i].Width
		r.Powups[i].Y = r.Engine.Cfg.WinHeight

		r.Powups[i].DX = r.Engine.Cfg.RDX * crnd()
		r.Powups[i].DY = -weightedRndRange(r.SpeedMin[direction], r.SpeedMax[direction]) + r.Engine.ScreenDY
	case UP:
		r.Powups[i].X = (frnd() * (r.Engine.Cfg.WinWidth + r.Powups[i].Width)) - r.Powups[i].Width
		r.Powups[i].Y = -r.Powups[i].Height

		r.Powups[i].DX = r.Engine.Cfg.RDX * crnd()
		r.Powups[i].DY = weightedRndRange(r.SpeedMin[direction], r.SpeedMax[direction]) + r.Engine.ScreenDY
	}

	r.Powups[i].Active = true

	if r.Timeout > r.Engine.Cfg.PowupsTimeout {
		r.Timeout = 0
	}
}

// Checks powups collisions
func (r *Powups) Collisions() {
	for i := 0; i < len(r.Powups); i++ {
		for n := 0; n < r.Engine.Cfg.MaxRocks; n++ {
			if r.Powups[i].Active && r.Game.Rocks.Rocks[n] != nil && r.Game.Rocks.Rocks[n].Active {
				//if r.Powups[i].Quad() != r.Game.Rocks.Rocks[n].Quad() {
				//continue
				//}

				if r.Powups[i].Collide(r.Game.Rocks.Rocks[n]) {
					// Bounce powup
					r.Powups[i].Bounce(r.Game.Rocks.Rocks[n])
				}
			}
		}
	}
}

// Updates powups
func (r *Powups) Update() {
	r.Timeout += int(r.Engine.FrameDelta)

	for i := 0; i < len(r.Powups); i++ {
		if r.Powups[i].Active {
			// Move powup
			r.Powups[i].Update()

			// Clip powup
			x := 100.0
			if r.Powups[i].X < -r.Powups[i].Width || r.Powups[i].X >= r.Engine.Cfg.WinWidth+x || r.Powups[i].Y < -(r.Powups[i].Height+x) || r.Powups[i].Y >= r.Engine.Cfg.WinHeight+x {
				r.Powups[i].Active = false
			}
		}
	}
}

// Draws powup glow
func (r *Powups) DrawGlow(i int) {
	n := (r.Engine.StartTicks / 75) % r.Engine.Cfg.NFrames
	src := &sdl.Rect{int32(uint32(r.Glow.Width) * n), 0, int32(r.Glow.Width), int32(r.Glow.Height)}
	dest := &sdl.Rect{int32(r.Powups[i].X - 8), int32(r.Powups[i].Y - 8), int32(r.Glow.Width), int32(r.Glow.Height)}
	r.Engine.Renderer.Copy(r.Glow.Texture, src, dest)
}

// Draws powups
func (r *Powups) Draw() {
	for i := 0; i < len(r.Powups); i++ {
		if r.Powups[i].Active {
			// Draw powup
			r.Powups[i].Draw()

			// Draw glow
			r.DrawGlow(i)

			// Explode powup
			r.Powups[i].Explode()
		}
	}
}
