package watcher

import (
	"os"
	"time"
)

/*
	This is a stupid file modification watcher, I expect to use fsnotify once it
	becomes consistent among all platforms.
*/
type Event struct {
	Name string
	// I don't really need any other event besides modification.
	isModify bool
}

type Watcher struct {
	Files    map[string]*WatcherFile
	Event    chan (*Event)
	t        time.Duration
	watching bool
}

type WatcherFile struct {
	Filemtime time.Time
}

func (self *Event) IsModify() bool {
	if self.isModify == true {
		self.isModify = false
		return true
	}
	return false
}

func (self *Watcher) RemoveWatch(file string) error {
	delete(self.Files, file)
	return nil
}

func (self *Watcher) Watch(file string) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	wf := &WatcherFile{
		Filemtime: stat.ModTime(),
	}
	self.Files[file] = wf
	return nil
}

func (self *Watcher) check() {
	for name, w := range self.Files {
		stat, err := os.Stat(name)
		if err == nil {
			mtime := stat.ModTime()
			if mtime != w.Filemtime {
				ev := &Event{
					Name:     name,
					isModify: true,
				}
				w.Filemtime = mtime
				self.Event <- ev
			}
		}
	}
}

func (self *Watcher) Close() {
	self.watching = false
}

func New() (*Watcher, error) {
	self := &Watcher{}
	self.t = time.Millisecond * 500
	self.Event = make(chan *Event)
	self.watching = true
	self.Files = make(map[string]*WatcherFile)

	go func() {
		for self.watching {
			self.check()
			time.Sleep(self.t)
		}
	}()

	return self, nil
}
