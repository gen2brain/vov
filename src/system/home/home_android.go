// +build android

package home

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Dir() string {
	return sdl.AndroidGetInternalStoragePath()
}
