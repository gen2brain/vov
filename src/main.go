package main

import "C"

import (
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/gen2brain/vov/src/engine"
	"github.com/gen2brain/vov/src/game"
)

func run() {
	// Initialize random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// Data directory
	dataDir := ""
	if runtime.GOOS != "android" {
		currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		dataDir = filepath.Join(currDir, "assets")
	}

	// Initialize SDL engine
	e := engine.NewEngine(engine.NewConfig())
	err := e.Init()
	if err != nil {
		// Show error message
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", err.Error(), nil)
		return
	}

	// Resources
	r := engine.NewResource(e, dataDir)

	// Set window icon
	if runtime.GOOS != "android" {
		e.SetIcon(r.LoadSurface(engine.ImageIcon))
	}

	// Set controller
	e.SetController(r.Mappings)

	// Set accelerometer
	if runtime.GOOS == "android" {
		// TODO
		// e.SetAccelerometer()
	}

	// Set haptic
	if e.Cfg.HapticEnabled {
		e.SetHaptic()
	}

	// Clear screen
	e.Clear()

	// Show loading message
	loading := game.NewSprite(e, r.LoadingText)
	loading.X = e.Cfg.WinWidth/2 - (loading.Width / 2)
	loading.Y = e.Cfg.WinHeight/2 - (loading.Height / 2)
	loading.Draw()

	// Update screen
	e.Renderer.Present()

	// Load resources
	r.Load()

	// Change state to menu
	e.State.Change(game.NewMenu(e, r))

	// Main loop
	for e.Running {

		// Calculate frame start
		e.StartFrame()

		// Handle events
		e.State.HandleEvents()

		// Clear screen
		e.Clear()

		// Update state
		e.State.Update()

		// Draw
		e.State.Draw()

		// Update screen
		e.Renderer.Present()

		// Calculate frame end
		e.EndFrame()
	}

	// Free resources
	r.Free()

	// Destroy SDL engine
	e.Destroy()
}

//export main2
func main2() {
	run()
}

// Go main function
func main() {
	run()
}
