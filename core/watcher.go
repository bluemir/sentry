package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bluemir/sentry/paths"
	"github.com/fsnotify/fsnotify"
)

type fsWatcher struct {
	watchPaths []string
	watcher    *fsnotify.Watcher
	done       chan bool
	filter     *fileNameFilter
}

func newFsWatcher(config *Config) *fsWatcher {

	filter := newFileNameFilter(config.Exclude)

	watchedFileList := paths.New(config.WatchPaths...).
		Glob().
		Expand(findAllDir).
		Filter(filter.check).
		Value()
	sort.Strings(watchedFileList)

	return &fsWatcher{
		watchPaths: watchedFileList,
		watcher:    nil,
		done:       make(chan bool),
		filter:     filter,
	}
}

func (fswatcher *fsWatcher) handleEvent(callback func()) {
	for {
		select {
		case event := <-fswatcher.watcher.Events:
			if event.Op == fsnotify.Chmod {
				//osx bug?
				continue
			}
			log.Infof("event: %s", event)

			if !fswatcher.filter.check(event.Name) {
				continue //skip exlude pattern
			}

			//if dir will add wather

			callback()
		case err := <-fswatcher.watcher.Errors:
			log.Debugln("error:", err)
		}
	}
}

func (fswatcher *fsWatcher) watch(callback func()) error {

	var err error
	fswatcher.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer fswatcher.watcher.Close()

	for _, path := range fswatcher.watchPaths {
		err = fswatcher.watcher.Add(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Infof("watching '%s'", path)
	}

	go fswatcher.handleEvent(callback)

	<-fswatcher.done
	return nil
}

func (fswatcher *fsWatcher) close() {
	fswatcher.done <- true
}

func findAllDir(path string) []string {
	list := []string{}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn(err)
			return nil //continue
		}
		if !strings.HasPrefix(path, ".") {
			list = append(list, path)
		}

		return nil
	})
	return list
}
