package core

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

type fsWatcher struct {
	path    string
	watcher *fsnotify.Watcher
	delayer *delayer
	done    chan bool
}

func newFsWatcher(path string, delay int32) *fsWatcher {
	return &fsWatcher{
		path:    path,
		watcher: nil,
		delayer: newDelayer(delay),
		done:    make(chan bool),
	}
}

func (fswatcher *fsWatcher) handleEvent(callback func()) {
	for {
		select {
		case event := <-fswatcher.watcher.Events:
			if event.Op == fsnotify.Chmod {
				//macos bug?
				continue
			}
			log.Infof("event: %s", event)

			fswatcher.delayer.Do(callback)
		case err := <-fswatcher.watcher.Errors:
			log.Debugln("error:", err)
		}
	}
}

func (fswatcher *fsWatcher) watch(callback func()) error {
	// TODO
	// * recursive listening
	// * event filtering

	var err error
	fswatcher.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer fswatcher.watcher.Close()
	go fswatcher.handleEvent(callback)

	list := findAllDir(fswatcher.path)
	list = append(list, fswatcher.path)
	for _, path := range list {
		log.Infof("watching '%s'", path)
		err = fswatcher.watcher.Add(path)

		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	<-fswatcher.done
	return nil
}

func findAllDir(path string) []string {
	list := []string{}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !strings.HasPrefix(path, ".") {
			list = append(list, path)
		}

		return nil
	})
	return list
}
