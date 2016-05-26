// +build !android

package rumble

import (
	"github.com/gen2brain/vov/src/engine"
)

func RumbleAvailable() bool {
	if engine.Haptic != nil && engine.Haptic.RumbleSupported() != 0 {
		return true
	}
	return false
}

func RumblePlay(strength float32, length uint32) {
	engine.Haptic.RumblePlay(strength, length)
}

func RumbleStop() {
	engine.Haptic.RumbleStop()
}
