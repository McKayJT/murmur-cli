// +build windows

package main

import (
	"golang.org/x/sys/windows"
	"os"
	"os/signal"
)

func CreateInterruptChannel() <-chan os.Signal {
	c := make(chan os.Signal, 5)
	signal.Notify(c, windows.SIGTERM, windows.SIGINT, windows.SIGQUIT)
	return c
}
