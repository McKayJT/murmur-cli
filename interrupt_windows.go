// +build windows

package main

import (
	"os"
	"os/signal"
	"golang.org/x/sys/windows"
)

func CreateInterruptChannel() <-chan os.Signal {
	c := make(chan os.Signal, 5)
	signal.Notify(c, windows.SIGTERM, windows.SIGINT, windows.SIGQUIT)
	return c
}
