// VoV engine
package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
)

var (
	// Paused boolean
	Paused bool

	// Started boolean
	Started bool

	// WasPaused boolean
	WasPaused bool

	// The time when the timer started
	StartTicks uint32

	// The ticks stored when the timer was paused
	PausedTicks uint32
)

// Starts timer
func StartTimer() {
	// Start the timer
	Started = true

	// Unpause the timer
	Paused = false

	// Get the current time
	StartTicks = sdl.GetTicks()
	PausedTicks = 0
}

// Stops timer
func StopTimer() {
	// Stop the timer
	Started = false

	// Unpause the timer
	Paused = false

	// Clear tick variables
	StartTicks = 0
	PausedTicks = 0
}

// Pauses timer
func Pause() {
	if Started && !Paused {
		// Pause the timer
		Paused = true

		// Calculate the paused ticks
		PausedTicks = sdl.GetTicks() - StartTicks
		StartTicks = 0

		// Pause music
		mix.PauseMusic()
	}
}

// Unpauses timer
func Unpause() {
	if Started && Paused {
		// Unpause the timer
		Paused = false

		// Reset the starting ticks
		StartTicks = sdl.GetTicks() - PausedTicks

		// Reset the paused ticks
		PausedTicks = 0

		// Resume music
		mix.ResumeMusic()
	}
}

// Returns timer's time
func GetTicks() uint32 {
	// The actual timer time
	var time uint32 = 0

	if Started {
		if Paused {
			// The number of ticks when the timer was paused
			time = PausedTicks
		} else {
			// The current time minus the start time
			time = sdl.GetTicks() - StartTicks
		}
	}

	return time
}
