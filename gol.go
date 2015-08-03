/*
 Conways game of life, in a console window
 by Telecoda - 2015
*/

package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

// Function main initializes termbox, renders the view, and starts
// handling events.
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := NewGame()
	g.render()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Key == termbox.KeyEsc:
					return
				}
			}
		default:
			//g.updateState()
			g.update()
			g.render()
			time.Sleep(animationSpeed)
		}
	}
}
