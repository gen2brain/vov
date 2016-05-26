// VoV engine
package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"

	"github.com/gen2brain/vov/src/system/log"
)

const (
	fontMain  = "OrbitronMedium.ttf"
	fontSmall = "OrbitronLight.ttf"
	fontTitle = "OrbitronBold.ttf"

	fontMainSize   = 24
	fontSmallSize  = 14
	fontTitleSize  = 52
	fontMediumSize = 18

	soundClick      = "click.ogg"
	soundBounce     = "bounce.ogg"
	soundEngine1    = "engine1.ogg"
	soundEngine2    = "engine2.ogg"
	soundEngine3    = "engine3.ogg"
	soundPowup0     = "powup0.ogg"
	soundPowup1     = "powup1.ogg"
	soundPowup2     = "powup2.ogg"
	soundPowup3     = "powup3.ogg"
	soundPowup4     = "powup4.ogg"
	soundPowup5     = "powup5.ogg"
	soundExplosion1 = "exp1.ogg"
	soundExplosion2 = "exp2.ogg"

	musicMenu = "factory-on-mercury.ogg"
	musicGame = "world-of-automatons.ogg"

	imageShip        = "ship.png"
	imageShipGlow    = "ship_glow.png"
	imageLife        = "life.png"
	imagePowup       = "powup.png"
	imagePowupGlow   = "powup_glow.png"
	imageBackground1 = "bg1.png"
	imageBackground2 = "bg2.png"
	imageBackground3 = "bg3.png"
	imageExplosion1  = "exp1.png"
	imageExplosion2  = "exp2.png"

	controllerMappings = "gamecontrollerdb.txt"

	// Export for window icon
	ImageIcon = "icon.png"
)

var (
	red   sdl.Color = sdl.Color{255, 36, 0, 255}
	green sdl.Color = sdl.Color{124, 252, 0, 255}
	brown sdl.Color = sdl.Color{224, 128, 26, 255}
	white sdl.Color = sdl.Color{255, 250, 250, 255}
)

const (
	FONT_SMALL = iota
	FONT_SMALL_RED
	FONT_MEDIUM
	FONT_LARGE
)

// Resource structure
type Resource struct {
	Engine *Engine

	DataDir string

	Mappings []string

	FontMain   *ttf.Font
	FontSmall  *ttf.Font
	FontTitle  *ttf.Font
	FontMedium *ttf.Font

	SoundClick      *mix.Chunk
	SoundBounce     *mix.Chunk
	SoundEngine1    *mix.Chunk
	SoundEngine2    *mix.Chunk
	SoundEngine3    *mix.Chunk
	SoundPowup0     *mix.Chunk
	SoundPowup1     *mix.Chunk
	SoundPowup2     *mix.Chunk
	SoundPowup3     *mix.Chunk
	SoundPowup4     *mix.Chunk
	SoundPowup5     *mix.Chunk
	SoundExplosion1 *mix.Chunk
	SoundExplosion2 *mix.Chunk

	MusicMenu *mix.Music
	MusicGame *mix.Music

	Ship     *sdl.Texture
	ShipSurf *sdl.Surface
	ShipGlow *sdl.Texture

	Powup     *sdl.Texture
	PowupSurf *sdl.Surface
	PowupGlow *sdl.Texture

	Life        *sdl.Texture
	Background1 *sdl.Texture
	Background2 *sdl.Texture
	Background3 *sdl.Texture
	Explosion1  *sdl.Texture
	Explosion2  *sdl.Texture

	LoadingText      *sdl.Texture
	TitleText        *sdl.Texture
	HiScoreText      *sdl.Texture
	HiScoreEnterText *sdl.Texture

	StartText     *sdl.Texture
	StartTextHi   *sdl.Texture
	ScoresText    *sdl.Texture
	ScoresTextHi  *sdl.Texture
	OptionsText   *sdl.Texture
	OptionsTextHi *sdl.Texture
	CreditsText   *sdl.Texture
	CreditsTextHi *sdl.Texture

	ProgrammingText          *sdl.Texture
	ProgrammingCreditText    *sdl.Texture
	MusicAndSoundsText       *sdl.Texture
	MusicAndSoundsCreditText *sdl.Texture
	GraphicsText             *sdl.Texture
	GraphicsCreditText       *sdl.Texture
	FontText                 *sdl.Texture
	FontCreditText           *sdl.Texture
	BasedText                *sdl.Texture
	BasedCreditText          *sdl.Texture
	SDLText                  *sdl.Texture
	SDLCreditText            *sdl.Texture
	GoText                   *sdl.Texture
	GoCreditText             *sdl.Texture
	VoVText                  *sdl.Texture
	VoVCreditText            *sdl.Texture

	MusicText           *sdl.Texture
	MusicTextHi         *sdl.Texture
	SoundsText          *sdl.Texture
	SoundsTextHi        *sdl.Texture
	AccelerometerText   *sdl.Texture
	AccelerometerTextHi *sdl.Texture
	HapticText          *sdl.Texture
	HapticTextHi        *sdl.Texture
	ShowFpsText         *sdl.Texture
	ShowFpsTextHi       *sdl.Texture

	YesText *sdl.Texture
	NoText  *sdl.Texture

	FpsText  *sdl.Texture
	TimeText *sdl.Texture

	ShieldsText     *sdl.Texture
	AttackText      *sdl.Texture
	InvincibleText  *sdl.Texture
	EngineBlastText *sdl.Texture
	SlowdownText    *sdl.Texture

	LifePowText        *sdl.Texture
	ShieldsPowText     *sdl.Texture
	AttackPowText      *sdl.Texture
	InvinciblePowText  *sdl.Texture
	EngineBlastPowText *sdl.Texture
	SlowdownPowText    *sdl.Texture

	PausedText   *sdl.Texture
	GameOverText *sdl.Texture

	Rocks     []*sdl.Texture
	RocksSurf []*sdl.Surface

	Glyphs           []string
	GlyphMapSmall    map[string]*Glyph
	GlyphMapSmallRed map[string]*Glyph
	GlyphMapMedium   map[string]*Glyph
	GlyphMapLarge    map[string]*Glyph
}

// Text glyph
type Glyph struct {
	Image  *sdl.Texture
	Width  float64
	Height float64
}

// Queries glyph texture dimensions
func (g *Glyph) Query() {
	_, _, w, h, err := g.Image.Query()
	if err != nil {
		log.Error("Query: %s\n", err)
		return
	}
	g.Width, g.Height = float64(w), float64(h)
}

// Returns new resource
func NewResource(e *Engine, d string) (r *Resource) {
	r = &Resource{}
	r.Engine = e
	r.DataDir = d

	r.Rocks = make([]*sdl.Texture, e.Cfg.NRocks)
	r.RocksSurf = make([]*sdl.Surface, e.Cfg.NRocks)

	r.GlyphMapSmall = make(map[string]*Glyph)
	r.GlyphMapSmallRed = make(map[string]*Glyph)
	r.GlyphMapMedium = make(map[string]*Glyph)
	r.GlyphMapLarge = make(map[string]*Glyph)

	r.Glyphs = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"`", "~", "!", "@", "#", "$", "%", "&", "*", "(", ")", "-", "_", "=", "+", "[", "]", "{", "}", ":", ";", "'", "\"", ".", ",", "<", ">", "/", "?", " ",
	}

	r.LoadMappings()

	r.FontMain = r.LoadFont(fontMain, fontMainSize)
	r.LoadingText = r.RenderText(r.FontMain, "L O A D I N G . . .", green, true, 0)

	return
}

// Loads resources
func (r *Resource) Load() {
	r.FontSmall = r.LoadFont(fontSmall, fontSmallSize)
	r.FontTitle = r.LoadFont(fontTitle, fontTitleSize)
	r.FontMedium = r.LoadFont(fontMain, fontMediumSize)

	r.SoundClick = r.LoadSound(soundClick)
	r.SoundBounce = r.LoadSound(soundBounce)
	r.SoundEngine1 = r.LoadSound(soundEngine1)
	r.SoundEngine2 = r.LoadSound(soundEngine2)
	r.SoundEngine3 = r.LoadSound(soundEngine3)
	r.SoundPowup0 = r.LoadSound(soundPowup0)
	r.SoundPowup1 = r.LoadSound(soundPowup1)
	r.SoundPowup2 = r.LoadSound(soundPowup2)
	r.SoundPowup3 = r.LoadSound(soundPowup3)
	r.SoundPowup4 = r.LoadSound(soundPowup4)
	r.SoundPowup5 = r.LoadSound(soundPowup5)
	r.SoundExplosion1 = r.LoadSound(soundExplosion1)
	r.SoundExplosion2 = r.LoadSound(soundExplosion2)

	r.MusicMenu = r.LoadMusic(musicMenu)
	r.MusicGame = r.LoadMusic(musicGame)

	r.Background1 = r.LoadTexture(imageBackground1)
	r.Background2 = r.LoadTexture(imageBackground2)
	r.Background3 = r.LoadTexture(imageBackground3)

	r.TitleText = r.RenderText(r.FontTitle, "V o V", green, true, 1)
	r.HiScoreText = r.RenderText(r.FontMain, "New High Score!", green, true, 0)
	r.HiScoreEnterText = r.RenderText(r.FontSmall, "Enter Your Name:", green, true, 0)

	r.StartText = r.RenderText(r.FontMain, "S T A R T", brown, true, 0)
	r.StartTextHi = r.RenderText(r.FontMain, "S T A R T", white, true, 0)
	r.ScoresText = r.RenderText(r.FontMain, "H A L L  O F  F A M E", brown, true, 0)
	r.ScoresTextHi = r.RenderText(r.FontMain, "H A L L  O F  F A M E", white, true, 0)
	r.OptionsText = r.RenderText(r.FontMain, "O P T I O N S", brown, true, 0)
	r.OptionsTextHi = r.RenderText(r.FontMain, "O P T I O N S", white, true, 0)
	r.CreditsText = r.RenderText(r.FontMain, "C R E D I T S", brown, true, 0)
	r.CreditsTextHi = r.RenderText(r.FontMain, "C R E D I T S", white, true, 0)

	r.MusicText = r.RenderText(r.FontMain, "M U S I C :", brown, true, 0)
	r.MusicTextHi = r.RenderText(r.FontMain, "M U S I C :", white, true, 0)
	r.SoundsText = r.RenderText(r.FontMain, "S O U N D S :", brown, true, 0)
	r.SoundsTextHi = r.RenderText(r.FontMain, "S O U N D S :", white, true, 0)
	r.AccelerometerText = r.RenderText(r.FontMain, "A C C E L E R O M E T E R :", brown, true, 0)
	r.AccelerometerTextHi = r.RenderText(r.FontMain, "A C C E L E R O M E T E R :", white, true, 0)
	r.HapticText = r.RenderText(r.FontMain, "R U M B L E :", brown, true, 0)
	r.HapticTextHi = r.RenderText(r.FontMain, "R U M B L E :", white, true, 0)
	r.ShowFpsText = r.RenderText(r.FontMain, "S H O W  F P S :", brown, true, 0)
	r.ShowFpsTextHi = r.RenderText(r.FontMain, "S H O W  F P S :", white, true, 0)

	r.ProgrammingText = r.RenderText(r.FontSmall, "Programming", red, true, 0)
	r.ProgrammingCreditText = r.RenderText(r.FontMedium, "M i l a n  N i k o l i c  (github.com/gen2brain)", green, true, 0)
	r.MusicAndSoundsText = r.RenderText(r.FontSmall, "Music and Sound Effects", red, true, 0)
	r.MusicAndSoundsCreditText = r.RenderText(r.FontMedium, "E r i c  M a t y a s  (soundimage.org)", green, true, 0)
	r.GraphicsText = r.RenderText(r.FontSmall, "Rocks Graphics", red, true, 0)
	r.GraphicsCreditText = r.RenderText(r.FontMedium, "P h a e l a x  (OpenGameArt.Org)", green, true, 0)
	r.FontText = r.RenderText(r.FontSmall, "Orbitron Font", red, true, 0)
	r.FontCreditText = r.RenderText(r.FontMedium, "M a t t  M c I n e r n e y  (github.com/theleagueof)", green, true, 0)
	r.BasedText = r.RenderText(r.FontSmall, "Based on VoR (Variations on Rockdodger) by", red, true, 0)
	r.BasedCreditText = r.RenderText(r.FontMedium, "J a s o n  W o o f e n d e n  (sametwice.com/vor)", green, true, 0)
	r.SDLText = r.RenderText(r.FontSmall, "Powered by SDL", red, true, 0)
	r.SDLCreditText = r.RenderText(r.FontMedium, "S i m p l e  D i r e c t  M e d i a  L a y e r  (libsdl.org)", green, true, 0)
	r.GoText = r.RenderText(r.FontSmall, "Written in Go", red, true, 0)
	r.GoCreditText = r.RenderText(r.FontMedium, "G o  l a n g u a g e  (golang.org)", green, true, 0)
	r.VoVText = r.RenderText(r.FontSmall, "Website", red, true, 0)
	r.VoVCreditText = r.RenderText(r.FontMedium, "V o V  (github.com/gen2brain/vov)", green, true, 0)

	r.YesText = r.RenderText(r.FontMain, "ON", green, true, 0)
	r.NoText = r.RenderText(r.FontMain, "OFF", brown, true, 0)

	r.FpsText = r.RenderText(r.FontSmall, "FPS: ", green, true, 0)
	r.TimeText = r.RenderText(r.FontSmall, "TIME: ", green, true, 0)

	r.ShieldsText = r.RenderText(r.FontSmall, "SHIELDS: ", green, true, 0)
	r.AttackText = r.RenderText(r.FontSmall, "ATTACK: ", green, true, 0)
	r.InvincibleText = r.RenderText(r.FontSmall, "INVINCIBLE: ", green, true, 0)
	r.EngineBlastText = r.RenderText(r.FontSmall, "ENGINE BLAST: ", green, true, 0)
	r.SlowdownText = r.RenderText(r.FontSmall, "SLOWDOWN: ", green, true, 0)

	r.LifePowText = r.RenderText(r.FontMain, "EXTRA LIFE", green, true, 0)
	r.ShieldsPowText = r.RenderText(r.FontMain, "SHIELDS", green, true, 0)
	r.AttackPowText = r.RenderText(r.FontMain, "ATTACK", green, true, 0)
	r.InvinciblePowText = r.RenderText(r.FontMain, "INVINCIBLE", green, true, 0)
	r.EngineBlastPowText = r.RenderText(r.FontMain, "ENGINE BLAST", green, true, 0)
	r.SlowdownPowText = r.RenderText(r.FontMain, "SLOWDOWN", green, true, 0)

	r.PausedText = r.RenderText(r.FontMain, "P A U S E D", green, true, 0)
	r.GameOverText = r.RenderText(r.FontMain, "G A M E  O V E R", green, true, 0)

	r.Ship = r.LoadTexture(imageShip)
	r.ShipSurf = r.LoadSurface(imageShip)
	r.ShipGlow = r.LoadTexture(imageShipGlow)

	r.Powup = r.LoadTexture(imagePowup)
	r.PowupSurf = r.LoadSurface(imagePowup)
	r.PowupGlow = r.LoadTexture(imagePowupGlow)

	r.Life = r.LoadTexture(imageLife)

	r.Background1 = r.LoadTexture(imageBackground1)
	r.Background2 = r.LoadTexture(imageBackground2)
	r.Background3 = r.LoadTexture(imageBackground3)

	r.Explosion1 = r.LoadTexture(imageExplosion1)
	r.Explosion2 = r.LoadTexture(imageExplosion2)

	r.LoadRocks()
	r.LoadGlyphs()
}

// Frees resources
func (r *Resource) Free() {
	r.FontMain.Close()
	r.FontSmall.Close()
	r.FontTitle.Close()
	r.FontMedium.Close()

	r.SoundClick.Free()
	r.SoundBounce.Free()
	r.SoundEngine1.Free()
	r.SoundEngine2.Free()
	r.SoundEngine3.Free()
	r.SoundPowup0.Free()
	r.SoundPowup1.Free()
	r.SoundPowup2.Free()
	r.SoundPowup3.Free()
	r.SoundPowup4.Free()
	r.SoundPowup5.Free()
	r.SoundExplosion1.Free()
	r.SoundExplosion2.Free()

	r.MusicMenu.Free()
	r.MusicGame.Free()

	r.Background1.Destroy()
	r.Background2.Destroy()
	r.Background3.Destroy()

	r.LoadingText.Destroy()
	r.TitleText.Destroy()
	r.HiScoreText.Destroy()
	r.HiScoreEnterText.Destroy()

	r.StartText.Destroy()
	r.StartTextHi.Destroy()
	r.ScoresText.Destroy()
	r.ScoresTextHi.Destroy()
	r.OptionsText.Destroy()
	r.OptionsTextHi.Destroy()
	r.CreditsText.Destroy()
	r.CreditsTextHi.Destroy()

	r.ProgrammingText.Destroy()
	r.ProgrammingCreditText.Destroy()
	r.MusicAndSoundsText.Destroy()
	r.MusicAndSoundsCreditText.Destroy()
	r.GraphicsText.Destroy()
	r.GraphicsCreditText.Destroy()
	r.FontText.Destroy()
	r.FontCreditText.Destroy()
	r.BasedText.Destroy()
	r.BasedCreditText.Destroy()
	r.SDLText.Destroy()
	r.SDLCreditText.Destroy()
	r.GoText.Destroy()
	r.GoCreditText.Destroy()

	r.YesText.Destroy()
	r.NoText.Destroy()
	r.FpsText.Destroy()
	r.TimeText.Destroy()
	r.PausedText.Destroy()
	r.GameOverText.Destroy()

	r.Ship.Destroy()
	r.ShipSurf.Free()
	r.ShipGlow.Destroy()

	r.Powup.Destroy()
	r.PowupSurf.Free()
	r.PowupGlow.Destroy()

	r.Life.Destroy()

	r.ShieldsText.Destroy()
	r.AttackText.Destroy()
	r.InvincibleText.Destroy()
	r.EngineBlastText.Destroy()
	r.SlowdownText.Destroy()

	r.LifePowText.Destroy()
	r.ShieldsPowText.Destroy()
	r.AttackPowText.Destroy()
	r.InvinciblePowText.Destroy()
	r.EngineBlastPowText.Destroy()
	r.SlowdownPowText.Destroy()

	r.Background1.Destroy()
	r.Background2.Destroy()
	r.Background3.Destroy()

	r.Explosion1.Destroy()
	r.Explosion2.Destroy()

	r.FreeRocks()
	r.FreeGlyphs()
}

// Loads rocks
func (r *Resource) LoadRocks() {
	rnd := func(min, max int) int {
		return rand.Intn(max-min) + min
	}

	for i := 0; i < r.Engine.Cfg.NRocks; i++ {
		rock := rnd(0, 13)
		file := filepath.Join("rocks", fmt.Sprintf("rock%02d.png", rock))

		s := r.LoadSurface(file)
		s.SetBlendMode(sdl.BLENDMODE_NONE)

		ratio := int(s.W / s.H)
		width := rnd(20, int(s.W)/ratio) * ratio
		height := width / ratio

		d, err := sdl.CreateRGBSurface(0, int32(width), int32(height), int32(s.Format.BitsPerPixel), s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask)
		if err != nil {
			log.Error("CreateRGBSurface: %s\n", err)
		}

		err = s.BlitScaled(nil, d, nil)
		if err != nil {
			log.Error("BlitScaled: %s\n", err)
		}

		s.Free()

		texture, err := r.Engine.Renderer.CreateTextureFromSurface(d)
		if err != nil {
			log.Error("CreateTextureFromSurface: %s\n", err)
		}

		r.Rocks[i] = texture
		r.RocksSurf[i] = d
	}
}

// Frees rocks
func (r *Resource) FreeRocks() {
	for _, t := range r.Rocks {
		t.Destroy()
	}
	for _, s := range r.RocksSurf {
		s.Free()
	}
}

// Loads glyphs
func (r *Resource) LoadGlyphs() {
	for _, g := range r.Glyphs {
		s := &Glyph{}
		s.Image = r.RenderText(r.FontSmall, g, green, true, 0)
		s.Query()

		r.GlyphMapSmall[g] = s

		s = &Glyph{}
		s.Image = r.RenderText(r.FontSmall, g, red, true, 0)
		s.Query()

		r.GlyphMapSmallRed[g] = s

		s = &Glyph{}
		s.Image = r.RenderText(r.FontMedium, g, brown, true, 0)
		s.Query()

		r.GlyphMapMedium[g] = s

		s = &Glyph{}
		s.Image = r.RenderText(r.FontMain, g, brown, true, 0)
		s.Query()

		r.GlyphMapLarge[g] = s
	}
}

// Frees glyphs
func (r *Resource) FreeGlyphs() {
	for _, g := range r.GlyphMapSmall {
		g.Image.Destroy()
	}

	for _, g := range r.GlyphMapSmallRed {
		g.Image.Destroy()
	}

	for _, g := range r.GlyphMapMedium {
		g.Image.Destroy()
	}

	for _, g := range r.GlyphMapLarge {
		g.Image.Destroy()
	}
}

// Loads texture
func (r *Resource) LoadTexture(filename string) (image *sdl.Texture) {
	var err error

	file := filepath.Join(r.DataDir, "images", filename)

	image, err = img.LoadTexture(r.Engine.Renderer, file)
	if err != nil {
		log.Error("LoadTexture: %s\n", err)
	}
	return
}

// Loads surface
func (r *Resource) LoadSurface(filename string) (image *sdl.Surface) {
	var err error

	file := filepath.Join(r.DataDir, "images", filename)

	image, err = img.Load(file)
	if err != nil {
		log.Error("LoadSurface: %s\n", err)
	}
	return
}

// Loads ttf font
func (r *Resource) LoadFont(filename string, size int) (font *ttf.Font) {
	var err error

	file := filepath.Join(r.DataDir, "fonts", filename)

	font, err = ttf.OpenFont(file, size)
	if err != nil {
		log.Error("LoadFont: %s\n", err)
	}
	return
}

// Loads music
func (r *Resource) LoadMusic(filename string) (music *mix.Music) {
	var err error

	file := filepath.Join(r.DataDir, "music", filename)

	music, err = mix.LoadMUS(file)
	if err != nil {
		log.Error("LoadMusic: %s\n", err)
	}
	return
}

// Loads sound
func (r *Resource) LoadSound(filename string) (sound *mix.Chunk) {
	var err error

	file := filepath.Join(r.DataDir, "sounds", filename)

	sound, err = mix.LoadWAV(file)
	if err != nil {
		log.Error("LoadSound: %s\n", err)
	}
	return
}

// Loads controllers mappings
func (r *Resource) LoadMappings() {
	r.Mappings = make([]string, 0)
	rw := sdl.RWFromFile(filepath.Join(r.DataDir, controllerMappings), "r")

	data := make([]byte, rw.RWsize())
	rw.RWread(unsafe.Pointer(&data[0]), 512, uint(rw.RWsize()/512))
	rw.RWclose()

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		r.Mappings = append(r.Mappings, scanner.Text())
	}
}

// Creates sdl texture from ttf font
func (r *Resource) RenderText(font *ttf.Font, text string, color sdl.Color, blended bool, outline int) (image *sdl.Texture) {
	var err error
	var surface *sdl.Surface

	if outline != 0 {
		font.SetOutline(outline)
	}

	if blended {
		surface, err = font.RenderUTF8_Blended(text, color)
	} else {
		surface, err = font.RenderUTF8_Solid(text, color)
	}

	if err != nil {
		log.Error("RenderText: %s\n", err)
		return
	}
	defer surface.Free()

	image, err = r.Engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Error("RenderText: %s\n", err)
	}

	return
}

// Draws text from glyph map
func (r *Resource) DrawText(text string, x, y int32, font int) {
	var img *sdl.Texture
	var dest *sdl.Rect = &sdl.Rect{}

	s := 0
	for i := 0; i < len(text); i++ {
		n := string(text[i])

		_, ok := r.GlyphMapSmall[n]
		if !ok {
			continue
		}

		dest.X = x + int32(s)
		dest.Y = y

		switch font {
		case FONT_SMALL:
			dest.W = int32(r.GlyphMapSmall[n].Width)
			dest.H = int32(r.GlyphMapSmall[n].Height)
			img = r.GlyphMapSmall[n].Image
			s += int(r.GlyphMapSmall[n].Width)
		case FONT_SMALL_RED:
			dest.W = int32(r.GlyphMapSmallRed[n].Width)
			dest.H = int32(r.GlyphMapSmallRed[n].Height)
			img = r.GlyphMapSmallRed[n].Image
			s += int(r.GlyphMapSmallRed[n].Width)
		case FONT_MEDIUM:
			dest.W = int32(r.GlyphMapMedium[n].Width)
			dest.H = int32(r.GlyphMapMedium[n].Height)
			img = r.GlyphMapMedium[n].Image
			s += int(r.GlyphMapMedium[n].Width)
		case FONT_LARGE:
			dest.W = int32(r.GlyphMapLarge[n].Width)
			dest.H = int32(r.GlyphMapLarge[n].Height)
			img = r.GlyphMapLarge[n].Image
			s += int(r.GlyphMapLarge[n].Width)
		}

		r.Engine.Renderer.Copy(img, nil, dest)
	}
}

// Plays sound
func (r *Resource) PlaySound(sound *mix.Chunk, channel int, loops int) {
	if r.Engine.Cfg.SoundsEnabled {
		_, err := sound.Play(channel, loops)
		if err != nil {
			log.Error("Play: %s\n", err)
		}
	}
}

// Plays sound timed
func (r *Resource) PlaySoundTimed(sound *mix.Chunk, channel int, loops int, ticks int) {
	if r.Engine.Cfg.SoundsEnabled {
		_, err := sound.PlayTimed(channel, loops, ticks)
		if err != nil {
			log.Error("PlayTimed: %s\n", err)
		}
	}
}

// Plays music
func (r *Resource) PlayMusic(music *mix.Music, loops int) {
	if r.Engine.Cfg.MusicEnabled {
		err := music.FadeIn(loops, 200)
		if err != nil {
			log.Error("Play: %s\n", err)
		}
	}
}
