// Copyright (c) 2012-2014 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package watcher is a stupid time-based file modification watcher, I expect
// to use fsnotify once it becomes consistent among all platforms. You should
// not use it.
package watcher

import (
	"os"
	"time"
)

// Event is the struct that watches a file.
type Event struct {
	Name string
	// We don't need other event besides modification.
	isModify bool
}

// Watcher is the struct that handles a the list of file to watch.
type Watcher struct {
	Files    map[string]*WatcherFile
	Event    chan (*Event)
	t        time.Duration
	watching bool
}

// WatcherFile is the struct that handles the last known file properties.
type WatcherFile struct {
	Filemtime time.Time
}

// IsModify returns true if the event was a file modification, then it resets
// the modified flag.
func (ev *Event) IsModify() bool {
	if ev.isModify == true {
		ev.isModify = false
		return true
	}
	return false
}

// RemoveWatch deletes a file from the watching list.
func (w *Watcher) RemoveWatch(file string) error {
	delete(w.Files, file)
	return nil
}

// Watch adds a file to the watching list.
func (w *Watcher) Watch(file string) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	wf := &WatcherFile{
		Filemtime: stat.ModTime(),
	}
	w.Files[file] = wf
	return nil
}

// Check compares the last known state of a file with the current state and
// updates modification flags, if required.
func (w *Watcher) check() {
	for name, f := range w.Files {
		stat, err := os.Stat(name)
		if err == nil {
			mtime := stat.ModTime()
			if mtime != f.Filemtime {
				ev := &Event{
					Name:     name,
					isModify: true,
				}
				f.Filemtime = mtime
				w.Event <- ev
			}
		}
	}
}

// Close makes a watcher sleep.
func (w *Watcher) Close() {
	w.watching = false
}

// New allocates, returns a file watcher and starts the watching loop.
func New() (*Watcher, error) {
	w := &Watcher{}
	w.t = time.Millisecond * 500
	w.Event = make(chan *Event)
	w.watching = true
	w.Files = make(map[string]*WatcherFile)

	go func() {
		for w.watching {
			w.check()
			time.Sleep(w.t)
		}
	}()

	return w, nil
}
