// VoV engine
package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gen2brain/vov/src/system/home"
	"github.com/gen2brain/vov/src/system/log"
)

// Configuration structure
type Config struct {
	// Enable music
	MusicEnabled bool

	// Enable sounds
	SoundsEnabled bool

	// Enable accelerometer
	AccelerometerEnabled bool

	// Accelerometer threshold
	AccelThreshold float64

	// Touch threshold
	TouchThreshold float64

	// Enable haptic
	HapticEnabled bool

	// Show frames per second
	ShowFps bool

	// Maximum frames per second
	MaxFps int

	// Barrier speed
	BarrierSpeed float64

	// Ship bounciness from screen edges
	Bounciness float64

	// Number of rock prototypes
	NRocks int

	// Number of rock structs to allocate
	MaxRocks int

	// Number of powups structs to allocate
	MaxPowups int

	// How often to generate powups
	PowupsTimeout int

	// How much milliseconds powup lasts
	PowupStateTimeout float64

	// Scale powup text
	PowupTextScale float64

	// Powup text timeout
	PowupTextTimeout float64

	// Number of animation sprites
	NFrames uint32

	// Gamespeed
	GameSpeed float64

	// Initial rocks
	InitialRocks int

	// Final rocks
	FinalRocks int

	// 32s for a speed=1 rock to cross the screen horizontally
	KH float64
	// 24s for a speed=1 rock to cross the screen vertically
	KV float64
	// range for rock dx values (+/-)
	RDX float64
	// range for rock dy values (+/-)
	RDY float64

	// Dust color depth
	MaxDustDepth float64

	// Number of dust motes
	NDustMotes int

	// Number od dust arrays
	NDustArray int

	// Maximum ship engine dots
	MaxShipDots int

	// Number od ship dots arrays
	NShipDotsArray int

	// Maximum bang dots
	MaxBangDots int

	// Number od bang dots arrays
	NBangDotsArray int

	// Ship thruster strength
	ThrusterStrength float64

	// How many engine dots come out of each thruster
	EngineDots int

	// Dots heat colors
	W int
	M int

	// Determines how hard dots push the rocks. Set to 0 to disable pushing rocks
	DotMassUnit float64

	// Time (in 1/60ths of a seccond) between when you blow up, and when next ship appears
	DeadPauseLength float64

	// Time (in milliseconds) to be invincible after next ship appears
	InvinciblePauseLength float64

	// Game over state timeout
	GameOverLength float64

	// Time (in milliseconds) to show scores screen after hiscore/gameover
	ScoresLength float64

	// Number of scores
	NScores int

	// Window width
	WinWidth float64

	// Window height
	WinHeight float64

	// X scroll
	XScrollTo float64
	// Y scroll
	YScrollTo float64

	// Distance ahead
	DistAhead float64

	// Maximum distance ahead
	MaxDistAhead float64

	// Credits scroll speed
	ScrollSpeed float64
}

// Returns new config
func NewConfig() (c *Config) {
	c = &Config{}
	if c.Exists() {
		c.Load()
	} else {
		c.Default()
	}
	return
}

// Sets default config
func (c *Config) Default() {
	c.MusicEnabled = true
	c.SoundsEnabled = true
	c.AccelerometerEnabled = false
	c.HapticEnabled = false
	c.ShowFps = false
	c.MaxFps = 60
	c.AccelThreshold = 500
	c.TouchThreshold = 0.05
	c.BarrierSpeed = 7.5
	c.Bounciness = 0.50
	c.NRocks = 65
	c.MaxRocks = 130
	c.MaxPowups = 30
	c.PowupsTimeout = 3000
	c.PowupStateTimeout = 10000
	c.PowupTextScale = 2.5
	c.PowupTextTimeout = 1000
	c.NFrames = 16
	c.GameSpeed = 1.0
	c.InitialRocks = 8
	c.FinalRocks = 25
	c.KH = 32 * 20
	c.KV = 24 * 20
	c.RDX = 2.5
	c.RDY = 2.5
	c.MaxDustDepth = 1
	c.NDustMotes = 1500
	c.NDustArray = 15
	c.MaxShipDots = 1500
	c.NShipDotsArray = 5
	c.MaxBangDots = 1500
	c.NBangDotsArray = 10
	c.EngineDots = 1000
	c.ThrusterStrength = 1.2
	c.W = 100
	c.M = 255
	c.DotMassUnit = 0.07
	c.DeadPauseLength = 40.0
	c.GameOverLength = 250.0
	c.ScoresLength = 10000.0
	c.NScores = 8
	c.InvinciblePauseLength = 2000.0
	c.XScrollTo = c.WinWidth / 3
	c.YScrollTo = c.WinHeight / 2
	c.MaxDistAhead = c.WinWidth
	c.ScrollSpeed = 0.8
}

// Loads config from file
func (c *Config) Load() {
	file := filepath.Join(home.Dir(), ".vov", "config")
	js, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("ReadFile: %s\n", err)
	}

	err = json.Unmarshal(js, c)
	if err != nil {
		log.Error("Unmarshal: %s\n", err)
	}
}

// Saves config to file
func (c *Config) Save() {
	dir := filepath.Join(home.Dir(), ".vov")
	if _, err := os.Stat(dir); err != nil {
		os.Mkdir(dir, 0755)
	}

	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Error("Marshal: %s\n", err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, "config"), js, 0644)
	if err != nil {
		log.Error("WriteFile: %s\n", err)
	}
}

// Checks if config file exists
func (c *Config) Exists() bool {
	file := filepath.Join(home.Dir(), ".vov", "config")
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}
