// VoV game
package game

import (
	"github.com/gen2brain/vov/src/engine"
)

// Rocks structure
type Rocks struct {
	Engine   *engine.Engine
	Cfg      *engine.Config
	Resource *engine.Resource

	Rocks      []*Sprite
	Prototypes []*Sprite

	Nrocks         int
	InitialRocks   int
	FinalRocks     int
	NrocksTimer    float64
	NrocksIncTicks float64
	CurrentRock    int

	Ti       [4]float64
	Rtimers  [4]float64
	SpeedMin [4]float64
	SpeedMax [4]float64
}

// Returns new rocks
func NewRocks(e *engine.Engine, res *engine.Resource) (r *Rocks) {
	r = &Rocks{}
	r.Engine = e
	r.Cfg = e.Cfg
	r.Resource = res

	r.Nrocks = r.Cfg.InitialRocks
	r.InitialRocks = r.Cfg.InitialRocks
	r.FinalRocks = r.Cfg.FinalRocks
	return
}

// Initializes rocks
func (r *Rocks) Init() {
	r.Rocks = make([]*Sprite, r.Cfg.MaxRocks)
	r.Prototypes = make([]*Sprite, r.Cfg.NRocks)

	for i := 0; i < r.Cfg.NRocks; i++ {
		s := NewSprite(r.Engine, r.Resource.Rocks[i])

		s.Type = ROCK
		s.Flags = MOVE | DRAW | COLLIDE

		s.Surface = r.Resource.RocksSurf[i]

		s.Exp1 = NewSprite(r.Engine, r.Resource.Explosion1)
		s.Exp2 = NewSprite(r.Engine, r.Resource.Explosion2)

		s.Width /= float64(r.Cfg.NFrames)

		r.Prototypes[i] = s
	}

	r.Reset()
}

// Resets rocks
func (r *Rocks) Reset() {
	r.Nrocks = r.InitialRocks
	r.NrocksIncTicks = float64(2*60*20) / float64(r.FinalRocks-r.InitialRocks)
	r.NrocksTimer = 0
	r.CurrentRock = 0
}

// Computes the number of rocks/tick that should be coming from each side,
// and the speed ranges of rocks coming from each side
func (r *Rocks) Sides() {
	var dx0, dx1, dy0, dy1 float64
	var hfactor, vfactor float64

	for i := 0; i < 4; i++ {
		r.Ti[i] = 0
		r.SpeedMin[i] = 0
		r.SpeedMax[i] = 0
	}

	hfactor = float64(r.Nrocks) / r.Cfg.KH
	vfactor = float64(r.Nrocks) / r.Cfg.KV

	dx0 = -r.Cfg.RDX - r.Engine.ScreenDX
	dx1 = r.Cfg.RDX - r.Engine.ScreenDX
	dy0 = -r.Cfg.RDY - r.Engine.ScreenDY
	dy1 = r.Cfg.RDY - r.Engine.ScreenDY

	if dx0 < 0 {
		r.SpeedMax[RIGHT] = -dx0
		if dx1 < 0 {
			// Rocks moving left only. So the RIGHT side of the screen
			r.SpeedMin[RIGHT] = -dx1
			r.Ti[RIGHT] = -(dx0 + dx1) / 2
		} else {
			// Rocks moving left and right
			r.SpeedMax[LEFT] = dx1
			r.Ti[RIGHT] = -dx0 / 2
			r.Ti[LEFT] = dx1 / 2
		}
	} else {
		// Rocks moving right only. So the LEFT side of the screen
		r.SpeedMin[LEFT] = dx0
		r.SpeedMax[LEFT] = dx1
		r.Ti[LEFT] = (dx0 + dx1) / 2
	}

	r.Ti[LEFT] *= hfactor
	r.Ti[RIGHT] *= hfactor

	if dy0 < 0 {
		r.SpeedMax[DOWN] = -dy0
		if dy1 < 0 {
			// Rocks moving up only. So the BOTTOM of the screen
			r.SpeedMin[DOWN] = -dy1
			r.Ti[DOWN] = -(dy0 + dy1) / 2
		} else {
			// Rocks moving up and down
			r.SpeedMax[UP] = dy1
			r.Ti[DOWN] = -dy0 / 2
			r.Ti[UP] = dy1 / 2
		}
	} else {
		// Rocks moving down only. so the TOP of the screen
		r.SpeedMin[UP] = dy0
		r.SpeedMax[UP] = dy1
		r.Ti[UP] = (dy0 + dy1) / 2
	}

	r.Ti[UP] *= vfactor
	r.Ti[DOWN] *= vfactor
}

// Generates new rocks
func (r *Rocks) New() {

	if r.Nrocks < r.FinalRocks {
		r.NrocksTimer += r.Engine.TFrame
		if r.NrocksTimer >= r.NrocksIncTicks {
			r.NrocksTimer -= r.NrocksIncTicks
			r.Nrocks++
		}
	}

	r.Sides()

	// Generate rocks
	for i := 0; i < 4; i++ {

		// Increment timers
		r.Rtimers[i] += r.Ti[i] * r.Engine.TFrame

		for r.Rtimers[i] >= 1 {
			r.Rtimers[i] -= 1

			if r.CurrentRock >= r.Cfg.MaxRocks {
				r.CurrentRock = 0
			}

			if r.Rocks[r.CurrentRock] == nil || !r.Rocks[r.CurrentRock].Active {

				p := &Sprite{}
				*p = *r.Prototypes[urnd()%uint32(r.Cfg.NRocks)]

				switch i {
				case RIGHT:
					p.X = r.Cfg.WinWidth
					p.Y = frnd() * (r.Cfg.WinHeight + p.Height)

					p.DX = -weightedRndRange(r.SpeedMin[i], r.SpeedMax[i]) + r.Engine.ScreenDX
					p.DY = r.Cfg.RDY * crnd()
				case LEFT:
					p.X = -p.Width
					p.Y = frnd() * (r.Cfg.WinHeight + p.Height)

					p.DX = weightedRndRange(r.SpeedMin[i], r.SpeedMax[i]) + r.Engine.ScreenDX
					p.DY = r.Cfg.RDY * crnd()
				case DOWN:
					p.X = (frnd() * (r.Cfg.WinWidth + p.Width)) - p.Width
					p.Y = r.Cfg.WinHeight

					p.DX = r.Cfg.RDX * crnd()
					p.DY = -weightedRndRange(r.SpeedMin[i], r.SpeedMax[i]) + r.Engine.ScreenDY
				case UP:
					p.X = (frnd() * (r.Cfg.WinWidth + p.Width)) - p.Width
					p.Y = -p.Height

					p.DX = r.Cfg.RDX * crnd()
					p.DY = weightedRndRange(r.SpeedMin[i], r.SpeedMax[i]) + r.Engine.ScreenDY
				}

				p.Active = true
				p.Direction = rnd(0, 2)
				p.Life = int(p.Width * p.Height * 300)
				p.Flags = MOVE | DRAW | COLLIDE

				// Add prototype to rocks
				r.Rocks[r.CurrentRock] = p

				r.CurrentRock++
			}

		}
	}
}

// Checks rocks collisions
func (r *Rocks) Collisions() {
	for i := 0; i < r.Cfg.MaxRocks; i++ {
		for n := 0; n < r.Cfg.MaxRocks; n++ {
			if r.Rocks[i] != nil && r.Rocks[n] != nil {
				if !r.Rocks[i].Active || !r.Rocks[n].Active || &r.Rocks[i] == &r.Rocks[n] {
					continue
				}

				//if r.Rocks[i].Quad() != r.Rocks[n].Quad() {
				//continue
				//}

				if r.Rocks[i].Collide(r.Rocks[n]) {
					// Bounce rocks
					r.Rocks[i].Bounce(r.Rocks[n])

					if r.Rocks[i].Mass() < r.Rocks[n].Mass() {
						// Change direction
						r.Rocks[i].Direction = 1 - r.Rocks[i].Direction
					}
				}
			}
		}
	}
}

// Updates rocks
func (r *Rocks) Update() {
	for i := 0; i < r.Cfg.MaxRocks; i++ {
		if r.Rocks[i] != nil && r.Rocks[i].Active {
			// Move rock
			r.Rocks[i].Update()

			// Clip rock
			if r.Rocks[i].X < -r.Rocks[i].Width || r.Rocks[i].X >= r.Cfg.WinWidth || r.Rocks[i].Y < -r.Rocks[i].Height || r.Rocks[i].Y >= r.Cfg.WinHeight {
				r.Rocks[i].Active = false
			}
		}
	}
}

// Draws rocks
func (r *Rocks) Draw() {
	for i := 0; i < r.Cfg.MaxRocks; i++ {
		if r.Rocks[i] != nil && r.Rocks[i].Active {
			// Draw rock
			r.Rocks[i].Draw()

			// Explode rock
			r.Rocks[i].Explode()
		}
	}
}
