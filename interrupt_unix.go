// +build !windows

package main

import (
	"os"
	"os/signal"
	"golang.org/x/sys/unix"
)

func CreateInterruptChannel() <-chan os.Signal {
	c := make(chan os.Signal, 5)
	signal.Notify(c, unix.SIGTERM, unix.SIGINT, unix.SIGQUIT)
	return c
}
