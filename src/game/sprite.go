// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/system/log"
)

// Sprite types
const (
	BASE = iota
	SHIP
	ROCK
	POWUP
)

// Sprite flags
const (
	MOVE    = 1
	DRAW    = 2
	COLLIDE = 4
)

// Sprite structure
type Sprite struct {
	Engine *engine.Engine

	Type   int
	Flags  int
	X      float64
	Y      float64
	DX     float64
	DY     float64
	Width  float64
	Height float64
	Frame  uint32
	Active bool

	Texture *sdl.Texture
	Surface *sdl.Surface

	// Explosion
	Exp1      *Sprite
	Exp2      *Sprite
	ExpActive bool
	ExpFrame  int

	// Rock extras
	Life      int
	Direction int

	// Ship extras
	Lives int
	Jets  int
	State int
}

// Returns new sprite
func NewSprite(e *engine.Engine, texture *sdl.Texture) (s *Sprite) {
	s = &Sprite{}
	s.Engine = e

	s.Flags = DRAW
	s.Texture = texture
	s.Query()

	return
}

// Queries sprite texture dimensions
func (s *Sprite) Query() {
	_, _, w, h, err := s.Texture.Query()
	if err != nil {
		log.Error("Query: %s\n", err)
		return
	}

	s.Width = float64(w)
	s.Height = float64(h)
}

// Returns sprite rect
func (s *Sprite) Rect() *sdl.Rect {
	return &sdl.Rect{int32(s.X), int32(s.Y), int32(s.Width), int32(s.Height)}
}

// Returns sprite quadrant
func (s *Sprite) Quad() int {
	w := s.Engine.Cfg.WinWidth / 2
	h := s.Engine.Cfg.WinHeight / 2

	top := s.Y < h && s.Y+s.Height < h
	bottom := s.Y > h

	if s.X < w && s.X+s.Width < w {
		if top {
			return 1
		} else if bottom {
			return 2
		}
	} else if s.X > w {
		if top {
			return 0
		} else if bottom {
			return 3
		}
	}

	return -1
}

// Returns mass of sprite
func (s *Sprite) Mass() (m float64) {
	if s.Type == SHIP || s.Type == POWUP {
		m = float64(s.Width * s.Height)
	} else if s.Type == ROCK {
		m = float64(3 * s.Width * s.Height)
	} else {
		m = 0
	}
	return
}

// Checks collide flag
func (s *Sprite) Collides() bool {
	return s.Flags&COLLIDE != 0
}

// Kills sprite
func (s *Sprite) Kill() {
	s.Flags = 0
	s.ExpActive = true
}

// Checks for sprite collision
func (s *Sprite) Collide(b *Sprite) bool {
	if !s.Collides() || !b.Collides() {
		return false
	}

	aRect := s.Rect()
	bRect := b.Rect()

	if aRect.HasIntersection(bRect) {
		iRect, _ := aRect.Intersect(bRect)

		aPix := s.Surface.Pixels()
		bPix := b.Surface.Pixels()

		var x, y int32

		for x = 0; x < (s.Surface.Pitch / int32(s.Surface.BytesPerPixel())); x++ {
			for y = 0; y < iRect.H; y++ {

				x1 := iRect.X - int32(s.X) + x
				y1 := iRect.Y - int32(s.Y) + y

				if s.Type == ROCK {
					if s.Direction == 1 {
						x1 = x1 + int32(uint32(s.Width)*(s.Engine.Cfg.NFrames-1)-s.Frame*uint32(s.Width))
					} else {
						x1 = x1 + (int32(s.Width) * int32(s.Frame))
					}
				}

				i := (y1 * (s.Surface.Pitch / int32(s.Surface.BytesPerPixel()))) + x1

				_, _, _, alpha := sdl.GetRGBA(readUint32(aPix[i:i+4]), s.Surface.Format)

				if alpha == 255 {
					return true
				}
			}
		}

		for x = 0; x < (b.Surface.Pitch / int32(b.Surface.BytesPerPixel())); x++ {
			for y = 0; y < iRect.H; y++ {

				x1 := iRect.X - int32(b.X) + x
				y1 := iRect.Y - int32(b.Y) + y

				if b.Type == ROCK {
					if b.Direction == 1 {
						x1 = x1 + int32(uint32(b.Width)*(s.Engine.Cfg.NFrames-1)-b.Frame*uint32(b.Width))
					} else {
						x1 = x1 + (int32(b.Width) * int32(b.Frame))
					}
				}

				i := (y1 * (b.Surface.Pitch / int32(b.Surface.BytesPerPixel()))) + x1

				_, _, _, alpha := sdl.GetRGBA(readUint32(bPix[i:i+4]), b.Surface.Format)

				if alpha == 255 {
					return true
				}
			}
		}
	}

	return false
}

// Explodes sprite
func (s *Sprite) Explode() {
	if !s.ExpActive || s.Exp1 == nil {
		return
	}

	if s.ExpFrame == int(s.Engine.Cfg.NFrames)-1 {
		// Last explosion frame
		if s.Type == ROCK {
			s.Flags = MOVE
		}

		s.ExpFrame = 0
		s.ExpActive = false
	} else {
		if s.Type == ROCK {
			if s.Width >= 48 {
				// Large explosion
				width := int32(s.Exp1.Width / float64(s.Engine.Cfg.NFrames))
				src := &sdl.Rect{width * int32(s.ExpFrame), 0, int32(width), int32(s.Exp1.Height)}
				dest := &sdl.Rect{int32(s.X) + int32(s.Width)/2 - width/2, int32(s.Y) + int32(s.Height)/2 - width/2, int32(width), int32(s.Exp1.Height)}
				s.Engine.Renderer.Copy(s.Exp1.Texture, src, dest)
			} else {
				// Small explosion
				width := int32(s.Exp2.Width / float64(s.Engine.Cfg.NFrames))
				src := &sdl.Rect{width * int32(s.ExpFrame), 0, int32(width), int32(s.Exp2.Height)}
				dest := &sdl.Rect{int32(s.X) + int32(s.Width)/2 - width/2, int32(s.Y) + int32(s.Height)/2 - width/2, int32(width), int32(s.Exp2.Height)}
				s.Engine.Renderer.Copy(s.Exp2.Texture, src, dest)
			}
		} else {
			width := int32(s.Exp1.Width / float64(s.Engine.Cfg.NFrames))
			src := &sdl.Rect{width * int32(s.ExpFrame), 0, int32(width), int32(s.Exp1.Height)}
			dest := &sdl.Rect{int32(s.X) + int32(s.Width)/2 - width/2, int32(s.Y) + int32(s.Height)/2 - width/2, int32(width), int32(s.Exp1.Height)}
			s.Engine.Renderer.Copy(s.Exp1.Texture, src, dest)
		}
	}

	if !engine.Paused {
		// Increment explosion frame
		s.ExpFrame++
	}
}

// Bounces sprites
func (s *Sprite) Bounce(b *Sprite) {
	var x, y, n float64 // (x, y) is unit vector from a to b
	var va, vb float64  // va, vb are balls' speeds along (x, y)
	var ma, mb float64  // ma, mb are the balls' masses
	var vc float64      // vc is the "center of momentum"

	// (x, y) is unit vector pointing from A's center to B's center
	x = (b.X + b.Width/2) - (s.X + s.Width/2)
	y = (b.Y + b.Height/2) - (s.Y + s.Height/2)
	n = math.Sqrt(x*x + y*y)

	x /= n
	y /= n

	// velocities along (x, y)
	va = x*s.DX + y*s.DY
	vb = x*b.DX + y*b.DY

	if vb-va > 0 {
		// don't bounce if we're already moving away
		return
	}

	// get masses and compute "center" speed
	ma = s.Mass()
	mb = b.Mass()
	vc = (va*ma + vb*mb) / (ma + mb)

	// bounce off the center speed
	s.DX += 2 * x * (vc - va)
	s.DY += 2 * y * (vc - va)

	b.DX += 2 * x * (vc - vb)
	b.DY += 2 * y * (vc - vb)
}

// Updates sprite
func (s *Sprite) Update() {
	if s.Flags&MOVE != 0 {
		s.X += (s.DX - s.Engine.ScreenDX) * s.Engine.TFrame
		s.Y += (s.DY - s.Engine.ScreenDY) * s.Engine.TFrame
	}
}

// Draws sprite
func (s *Sprite) Draw() {
	if s.Flags&DRAW == 0 {
		return
	}

	if s.Type == ROCK {
		s.Frame = (s.Engine.StartTicks / 50) % s.Engine.Cfg.NFrames

		dest := s.Rect()
		src := &sdl.Rect{int32(uint32(s.Width) * s.Frame), 0, int32(s.Width), int32(s.Height)}

		if s.Direction == 1 {
			// reversed sprite frame direction
			src = &sdl.Rect{int32(uint32(s.Width)*(s.Engine.Cfg.NFrames-1) - s.Frame*uint32(s.Width)), 0, int32(s.Width), int32(s.Height)}
		}

		s.Engine.Renderer.Copy(s.Texture, src, dest)

	} else if s.Type == SHIP || s.Type == POWUP {
		dest := s.Rect()
		src := &sdl.Rect{int32(uint32(s.Width) * uint32(s.State)), 0, int32(s.Width), int32(s.Height)}

		s.Engine.Renderer.Copy(s.Texture, src, dest)
	} else {
		dest := s.Rect()
		s.Engine.Renderer.Copy(s.Texture, nil, dest)
	}
}
