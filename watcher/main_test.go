package watcher

import (
	"fmt"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	w, _ := New()
	w.Watch("main.go")

	go func() {
		for {
			select {
			case ev := <-w.Event:
				if ev.IsModify() {
					fmt.Printf("File %s was modified.\n", ev.Name)
				}
			}
		}
	}()

	time.Sleep(120 * time.Second)

	w.Close()
}
