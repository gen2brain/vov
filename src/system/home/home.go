// +build !android,!windows

package home

import (
	"os"
)

func Dir() string {
	return os.Getenv("HOME")
}
