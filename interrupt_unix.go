// +build !windows

package main

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

func CreateInterruptChannel() <-chan os.Signal {
	c := make(chan os.Signal, 5)
	signal.Notify(c, unix.SIGTERM, unix.SIGINT, unix.SIGQUIT)
	return c
}
