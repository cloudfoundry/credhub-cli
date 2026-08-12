//go:build !windows
// +build !windows

package config

import (
	"os"
)

func userHomeDir() string {
	home := os.Getenv("CREDHUB_HOME")
	if home == "" {
		home = os.Getenv("HOME")
	}
	return home
}

func makeDirectory() error {
	return os.MkdirAll(ConfigDir(), 0755)
}
