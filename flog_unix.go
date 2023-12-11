//go:build !windows
// +build !windows

package main

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/DataDog/datadog-go/v5/statsd"
)

// Run checks overwrite flag and generates logs with given options
func Run(option *Option) error {
	logDir := filepath.Dir(option.Output)
	oldMask := syscall.Umask(0000)
	if err := os.MkdirAll(logDir, 0766); err != nil {
		return err
	}
	syscall.Umask(oldMask)

	var statsdClient *statsd.Client
	var err error

	if option.Statsd != "" {
		statsdClient, err = statsd.New(option.Statsd)
		if err != nil {
			return errors.New("can't create the statsd client")
		}
	}

	if _, err = os.Stat(option.Output); err == nil && !option.Overwrite {
		return errors.New(option.Output + " already exists. You can overwrite with -w option")
	}

	return Generate(statsdClient, option)
}
