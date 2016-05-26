// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
)

// Dust mote structure
type DustMote struct {
	X float64
	Y float64
	Z float64
}

// Dust structure
type Dust struct {
	Engine *engine.Engine
	Cfg    *engine.Config

	// Dust motes
	Motes map[sdl.Color][]*DustMote

	// Motes points
	MotesPoints map[sdl.Color][]sdl.Point
}

// Returns new dust
func NewDust(e *engine.Engine) (d *Dust) {
	d = &Dust{}
	d.Engine = e
	d.Cfg = e.Cfg
	return
}

// Initializes dust
func (d *Dust) Init() {
	d.Motes = make(map[sdl.Color][]*DustMote)
	d.MotesPoints = make(map[sdl.Color][]sdl.Point)

	for i := 0; i < d.Cfg.NDustArray; i++ {

		motes := make([]*DustMote, d.Cfg.NDustMotes/d.Cfg.NDustArray)
		points := make([]sdl.Point, d.Cfg.NDustMotes/d.Cfg.NDustArray)

		z := d.Cfg.MaxDustDepth * math.Sqrt(frnd())
		c := (d.Cfg.MaxDustDepth - z) * 255.0 / d.Cfg.MaxDustDepth
		color := sdl.Color{uint8(c), uint8(c), uint8(c), 255}

		for n := 0; n < d.Cfg.NDustMotes/d.Cfg.NDustArray; n++ {
			p := sdl.Point{}
			p.X = int32(frnd() * (d.Cfg.WinWidth - 5))
			p.Y = int32(frnd() * (d.Cfg.WinHeight - 5))
			points[n] = p

			m := &DustMote{}
			m.X = float64(p.X)
			m.Y = float64(p.Y)
			m.Z = d.Cfg.MaxDustDepth * math.Sqrt(frnd())
			motes[n] = m
		}

		d.Motes[color] = motes
		d.MotesPoints[color] = points
	}
}

// Updates dust
func (d *Dust) Update() {
	xscroll := d.Engine.ScreenDX * d.Engine.TFrame
	yscroll := d.Engine.ScreenDY * d.Engine.TFrame

	for color, _ := range d.MotesPoints {
		for n := 0; n < d.Cfg.NDustMotes/d.Cfg.NDustArray; n++ {

			x := float64(d.Motes[color][n].X)
			y := float64(d.Motes[color][n].Y)

			x -= xscroll / (1.3 + d.Motes[color][n].Z)
			x = fwrap(x, d.Cfg.WinWidth)

			y -= yscroll / (1.3 + d.Motes[color][n].Z)
			y = fwrap(y, d.Cfg.WinHeight)

			d.Motes[color][n].X = x
			d.Motes[color][n].Y = y

			d.MotesPoints[color][n].X = int32(x)
			d.MotesPoints[color][n].Y = int32(y)
		}
	}
}

// Draws dust
func (d *Dust) Draw() {
	for color, points := range d.MotesPoints {
		d.Engine.Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		d.Engine.Renderer.DrawPoints(points)
	}
}
