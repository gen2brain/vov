// VoV game
package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
)

// Dot types
const (
	BANGDOT = iota
	ENGINEDOT
)

// Dot structure
type Dot struct {
	Game *Game

	X      float64
	Y      float64
	DX     float64
	DY     float64
	Mass   float64
	Decay  float64 // Rate at which to reduce mass
	Type   int
	Active bool
}

// Moves dot
func (d *Dot) Move() {
	g := d.Game
	c := d.Game.Engine.Cfg

	if d.Active {
		d.X += (d.DX - g.Engine.ScreenDX) * g.Engine.TFrame
		d.Y += (d.DY - g.Engine.ScreenDY) * g.Engine.TFrame
		d.Mass -= g.Engine.TFrame * d.Decay

		if d.Mass < 0 || fclip(d.X, c.WinWidth) || fclip(d.Y, c.WinHeight) {
			d.Active = false
		} else {
			d.Collisions()
		}
	}
}

// Checks for dot collision
func (d *Dot) Collide(b *Sprite) bool {
	if !b.Collides() {
		return false
	}

	aRect := &sdl.Rect{int32(d.X), int32(d.Y), 1, 1}
	bRect := b.Rect()

	if aRect.HasIntersection(bRect) {
		return true
	}

	return false
}

// Checks dot collisions
func (d *Dot) Collisions() {
	g := d.Game
	c := d.Game.Engine.Cfg

	for i := 0; i < len(g.Powups.Powups); i++ {
		if g.Powups.Powups[i].Active {
			if d.Collide(g.Powups.Powups[i]) {
				// Push the powup
				g.Powups.Powups[i].DX += c.DotMassUnit * d.Mass * (d.DX - g.Powups.Powups[i].DX) / g.Powups.Powups[i].Mass()
				g.Powups.Powups[i].DY += c.DotMassUnit * d.Mass * (d.DY - g.Powups.Powups[i].DY) / g.Powups.Powups[i].Mass()
			}
		}
	}

	for i := 0; i < c.MaxRocks; i++ {
		if g.Rocks.Rocks[i] != nil && g.Rocks.Rocks[i].Active {

			if d.Collide(g.Rocks.Rocks[i]) {
				//d.Active = false

				if d.Type == BANGDOT {
					// Fry rock with bang dots
					g.Rocks.Rocks[i].Life -= int((d.DX-g.Rocks.Rocks[i].DX)*(d.DX-g.Rocks.Rocks[i].DX) +
						(d.DY-g.Rocks.Rocks[i].DY)*(d.DY-g.Rocks.Rocks[i].DY))

					// Push the rock
					g.Rocks.Rocks[i].DX += c.DotMassUnit * d.Mass * (d.DX - g.Rocks.Rocks[i].DX) / g.Rocks.Rocks[i].Mass()
					g.Rocks.Rocks[i].DY += c.DotMassUnit * d.Mass * (d.DY - g.Rocks.Rocks[i].DY) / g.Rocks.Rocks[i].Mass()
				} else if d.Type == ENGINEDOT && g.Ship.State == ENGINEBLAST {
					// Fry rocks with engine dots
					g.Rocks.Rocks[i].Life -= int((d.DX-g.Rocks.Rocks[i].DX)*(d.DX-g.Rocks.Rocks[i].DX) +
						(d.DY-g.Rocks.Rocks[i].DY)*(d.DY-g.Rocks.Rocks[i].DY)*1.5)

					// Push the rock
					g.Rocks.Rocks[i].DX += c.DotMassUnit * d.Mass * (d.DX - g.Rocks.Rocks[i].DX) / g.Rocks.Rocks[i].Mass()
					g.Rocks.Rocks[i].DY += c.DotMassUnit * d.Mass * (d.DY - g.Rocks.Rocks[i].DY) / g.Rocks.Rocks[i].Mass()
				} else if d.Type == ENGINEDOT {
					// Push the rock
					g.Rocks.Rocks[i].DX += c.DotMassUnit * d.Mass * (d.DX - g.Rocks.Rocks[i].DX) / g.Rocks.Rocks[i].Mass()
					g.Rocks.Rocks[i].DY += c.DotMassUnit * d.Mass * (d.DY - g.Rocks.Rocks[i].DY) / g.Rocks.Rocks[i].Mass()
				}

				if g.Rocks.Rocks[i].Life < 0 {
					// Kill rock if out of life
					g.Resource.PlaySound(g.Resource.SoundExplosion2, -1, 0)
					g.Rocks.Rocks[i].Kill()

					// Bang dots
					g.Dots.NewBangDots(g.Rocks.Rocks[i])
				}
			}
		}
	}
}

// Dots structure
type Dots struct {
	Game *Game
	Cfg  *engine.Config

	Bd int
	S  [4]float64

	// Slice of heat colors
	HeatColors []sdl.Color

	ShipDots []*Dot
	BangDots []*Dot

	ShipDotsColors []sdl.Color
	BangDotsColors []sdl.Color

	ShipDotsPoints map[int][]sdl.Point
	BangDotsPoints map[int][]sdl.Point
}

// Returns new dots
func NewDots(g *Game) (d *Dots) {
	d = &Dots{}
	d.Game = g
	d.Cfg = g.Engine.Cfg
	return
}

// Initializes dots
func (d *Dots) Init() {
	d.ShipDots = make([]*Dot, d.Cfg.MaxShipDots)
	d.BangDots = make([]*Dot, d.Cfg.MaxBangDots)

	d.ShipDotsColors = make([]sdl.Color, d.Cfg.NShipDotsArray)
	d.BangDotsColors = make([]sdl.Color, d.Cfg.NBangDotsArray)

	d.ShipDotsPoints = make(map[int][]sdl.Point)
	d.BangDotsPoints = make(map[int][]sdl.Point)

	d.S = [4]float64{2, 1, 0, 1}

	d.InitColors()

	for i := 0; i < d.Cfg.NShipDotsArray; i++ {
		points := make([]sdl.Point, d.Cfg.MaxShipDots/d.Cfg.NShipDotsArray)

		idx := urnd() % uint32(len(d.HeatColors))
		color := d.HeatColors[idx]

		for n := 0; n < d.Cfg.MaxShipDots/d.Cfg.NShipDotsArray; n++ {
			points[n] = sdl.Point{-200, -200}

			dot := &Dot{}
			dot.Game = d.Game
			dot.Active = false
			dot.Type = ENGINEDOT

			x := i + (n * d.Cfg.NShipDotsArray)

			d.ShipDots[x] = dot
		}

		d.ShipDotsColors[i] = color
		d.ShipDotsPoints[i] = points
	}

	for i := 0; i < d.Cfg.NBangDotsArray; i++ {
		points := make([]sdl.Point, d.Cfg.MaxBangDots/d.Cfg.NBangDotsArray)

		idx := urnd() % uint32(len(d.HeatColors))
		color := d.HeatColors[idx]

		for n := 0; n < d.Cfg.MaxBangDots/d.Cfg.NBangDotsArray; n++ {
			points[n] = sdl.Point{-200, -200}

			dot := &Dot{}
			dot.Game = d.Game
			dot.Active = false
			dot.Type = BANGDOT

			x := i + (n * d.Cfg.NBangDotsArray)

			d.BangDots[x] = dot
		}

		d.BangDotsColors[i] = color
		d.BangDotsPoints[i] = points
	}
}

// Initialize colors
func (d *Dots) InitColors() {
	var r, g, b int

	d.HeatColors = make([]sdl.Color, 0)

	for i := 0; i < d.Cfg.W*3; i++ {
		if i < d.Cfg.W {
			r = i * d.Cfg.M / d.Cfg.W
		} else {
			r = d.Cfg.M
		}

		if i < d.Cfg.W {
			g = 0
		} else {
			if i < 2*d.Cfg.W {
				g = (i - d.Cfg.W) * d.Cfg.M / d.Cfg.W
			} else {
				g = d.Cfg.M
			}
		}

		if i < 2*d.Cfg.W {
			b = 0
		} else {
			b = (i - d.Cfg.W) * d.Cfg.M / d.Cfg.W
		}

		d.HeatColors = append(d.HeatColors, sdl.Color{uint8(r), uint8(g), uint8(b), 255})
	}
}

// Generates new ship engine dots
func (d *Dots) NewShipDots() {

	n := d.Cfg.EngineDots

	for dir := 0; dir < 4; dir++ {
		if d.Game.Ship.Jets&(1<<uint(dir)) == 0 {
			continue
		}

		for i := 0; i < n; i++ {

			if !d.ShipDots[i].Active {
				a := frnd()*math.Pi + float64(dir-1)*(math.Pi/2) // angle
				r := math.Sin(frnd() * math.Pi)                  // random length

				dx := r * math.Cos(a)
				dy := r * -math.Sin(a) // Screen y is "backwards"

				d.ShipDots[i].Decay = 3.5

				// Dot was created at a random time during the time span
				time := frnd() * d.Game.Engine.TFrame // This is how long ago

				// Calculate how fast the ship was going when this engine dot was created (as if it had a smooth acceleration).
				// This is used in determining the velocity of the dots, but not their starting location.
				accelh := float64((d.Game.Ship.Jets>>2)&1) - float64(d.Game.Ship.Jets&1)
				accelh *= d.Cfg.ThrusterStrength * time
				pastShipDX := d.Game.Ship.DX - accelh

				accelv := float64((d.Game.Ship.Jets>>1)&1) - float64((d.Game.Ship.Jets>>3)&1)
				accelv *= d.Cfg.ThrusterStrength * time
				pastShipDY := d.Game.Ship.DY - accelv

				// The starting position (not speed) of the dot is calculated as though the ship were traveling at a constant speed for this TFrame.
				d.ShipDots[i].X = (d.Game.Ship.X - (d.Game.Ship.DX-d.Game.Engine.ScreenDX)*time) + d.S[dir]*(d.Game.Ship.Width/2)
				d.ShipDots[i].Y = (d.Game.Ship.Y - (d.Game.Ship.DY-d.Game.Engine.ScreenDY)*time) + d.S[(dir+1)&3]*(d.Game.Ship.Height/2)

				if dir&1 != 0 {
					d.ShipDots[i].DX = pastShipDX + 2*dx
					d.ShipDots[i].DY = pastShipDY + 20*dy
					d.ShipDots[i].Mass = 60 * math.Abs(dy)
				} else {
					d.ShipDots[i].DX = pastShipDX + 20*dx
					d.ShipDots[i].DY = pastShipDY + 2*dy
					d.ShipDots[i].Mass = 60 * math.Abs(dx)
				}

				// Move the dot as though it were created in the past
				d.ShipDots[i].X += (d.ShipDots[i].DX - d.Game.Engine.ScreenDX) * time
				d.ShipDots[i].Y += (d.ShipDots[i].DY - d.Game.Engine.ScreenDY) * time

				if !fclip(d.ShipDots[i].X, d.Cfg.WinWidth) && !fclip(d.ShipDots[i].Y, d.Cfg.WinHeight) {
					d.ShipDots[i].Active = true
				}
			}
		}
	}
}

// Generates new bang dots
func (d *Dots) NewShipBangDots() {
	n := 10
	for i := 0; i < n; i++ {
		for y := 0; y < int(d.Game.Ship.Height); y++ {
			for x := 0; x < int(d.Game.Ship.Width); x++ {
				theta := frnd() * math.Pi * 2

				r := frnd()
				r = 1 - r*r

				d.BangDots[d.Bd].DX = 45*r*math.Cos(theta) + d.Game.Ship.DX
				d.BangDots[d.Bd].DY = 45*r*math.Sin(theta) + d.Game.Ship.DY
				d.BangDots[d.Bd].X = float64(x) + d.Game.Ship.X
				d.BangDots[d.Bd].Y = float64(y) + d.Game.Ship.Y

				d.BangDots[d.Bd].Mass = frnd() * 99
				d.BangDots[d.Bd].Decay = frnd()*4.5 + 0.5
				d.BangDots[d.Bd].Active = true

				d.Bd = (d.Bd + 1) % d.Cfg.MaxBangDots
			}
		}
	}
}

// Generates new bang dots
func (d *Dots) NewBangDots(s *Sprite) {
	n := 5
	for i := 0; i < n; i++ {
		for y := 0; y < int(s.Height); y++ {
			for x := 0; x < int(s.Width); x++ {
				theta := frnd() * math.Pi * 2

				r := frnd()
				r = 1 - r*r

				d.BangDots[d.Bd].DX = 45*r*math.Cos(theta) + s.DX
				d.BangDots[d.Bd].DY = 45*r*math.Sin(theta) + s.DY
				d.BangDots[d.Bd].X = float64(x) + s.X
				d.BangDots[d.Bd].Y = float64(y) + s.Y

				d.BangDots[d.Bd].Mass = frnd() * 99
				d.BangDots[d.Bd].Decay = frnd()*1.5 + 0.5
				d.BangDots[d.Bd].Active = true

				d.Bd = (d.Bd + 1) % d.Cfg.MaxBangDots
			}
		}
	}
}

// Updates dots
func (d *Dots) Update() {
	for i := 0; i < d.Cfg.NShipDotsArray; i++ {
		for n := 0; n < d.Cfg.MaxShipDots/d.Cfg.NShipDotsArray; n++ {
			x := i + (n * d.Cfg.NShipDotsArray)
			d.ShipDots[x].Move()

			if d.ShipDots[x].Active {
				d.ShipDotsPoints[i][n].X = int32(d.ShipDots[x].X)
				d.ShipDotsPoints[i][n].Y = int32(d.ShipDots[x].Y)
			} else {
				d.ShipDotsPoints[i][n].X = -200
				d.ShipDotsPoints[i][n].Y = -200
			}
		}
	}

	for i := 0; i < d.Cfg.NBangDotsArray; i++ {
		for n := 0; n < d.Cfg.MaxBangDots/d.Cfg.NBangDotsArray; n++ {
			x := i + (n * d.Cfg.NBangDotsArray)
			d.BangDots[x].Move()

			if d.BangDots[x].Active {
				d.BangDotsPoints[i][n].X = int32(d.BangDots[x].X)
				d.BangDotsPoints[i][n].Y = int32(d.BangDots[x].Y)
			} else {
				d.BangDotsPoints[i][n].X = -200
				d.BangDotsPoints[i][n].Y = -200
			}
		}
	}
}

// Draws dots
func (d *Dots) Draw() {
	for i := 0; i < d.Cfg.NShipDotsArray; i++ {
		active := false

		for n := 0; n < d.Cfg.MaxShipDots/d.Cfg.NShipDotsArray; n++ {
			x := i + (n * d.Cfg.NShipDotsArray)

			if d.ShipDots[x].Active {
				active = true
				break
			}
		}

		if active {
			color := d.ShipDotsColors[i]

			d.Game.Engine.Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			d.Game.Engine.Renderer.DrawPoints(d.ShipDotsPoints[i])
		}
	}

	for i := 0; i < d.Cfg.NBangDotsArray; i++ {
		active := false

		for n := 0; n < d.Cfg.MaxBangDots/d.Cfg.NBangDotsArray; n++ {
			x := i + (n * d.Cfg.NBangDotsArray)

			if d.BangDots[x].Active {
				active = true
				break
			}
		}

		if active {
			color := d.BangDotsColors[i]

			d.Game.Engine.Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			d.Game.Engine.Renderer.DrawPoints(d.BangDotsPoints[i])
		}
	}
}
