// VoV game
package game

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
)

// Fog structure
type Fog struct {
	Sprite

	Engine   *engine.Engine
	Resource *engine.Resource

	ScrollOffset float64

	Backgrounds map[int]*sdl.Texture
}

// Returns new fog
func NewFog(e *engine.Engine, r *engine.Resource) *Fog {
	b := &Fog{}
	b.Engine = e
	b.Resource = r

	return b
}

// Initializes fog
func (b *Fog) Init() {
	b.Backgrounds = make(map[int]*sdl.Texture)
	b.Backgrounds[0] = b.Resource.Background1
	b.Backgrounds[1] = b.Resource.Background2
	b.Backgrounds[2] = b.Resource.Background3

	// Random background
	b.Texture = b.Backgrounds[rnd(0, 3)]
	b.Query()

	// Set texture transparency
	b.Texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	b.Texture.SetAlphaMod(30)

	b.Flags = DRAW
}

// Updates fog
func (b *Fog) Update() {
	// Scroll
	b.X -= b.Engine.ScreenDX * b.Engine.TFrame / 2.0
	if b.X < -b.Width {
		b.X = 0
	}

	b.ScrollOffset -= b.Engine.ScreenDX * b.Engine.TFrame / 2.0
	if b.ScrollOffset < -b.Width {
		b.ScrollOffset = 0
	}
}

// Draws fog
func (b *Fog) Draw() {
	if b.Flags&DRAW != 0 {
		src := b.Rect()
		b.Engine.Renderer.Copy(b.Texture, src, nil)

		src = &sdl.Rect{int32(b.X + b.Width), 0, int32(b.Width), int32(b.Height)}
		b.Engine.Renderer.Copy(b.Texture, src, nil)
	}
}
