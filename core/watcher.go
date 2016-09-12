package core

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

type fsWatcher struct {
	watchPaths []string
	watcher    *fsnotify.Watcher
	delayer    *delayer
	done       chan bool
	filter     *fileNameFilter
}

func newFsWatcher(config *Config) *fsWatcher {
	return &fsWatcher{
		watchPaths: config.WatchPaths,
		watcher:    nil,
		delayer:    newDelayer(config.Delay),
		done:       make(chan bool),
		filter:     newFileNameFilter(config.Exclude),
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

			if fswatcher.filter.check(event.Name) {
				continue //skip exlude pattern
			}

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

	list := append([]string{}, fswatcher.watchPaths...)
	list = expand(list, expandPath)
	log.Debug(list)
	list = expand(list, findAllDir)
	log.Debug(list)

	for _, path := range list {
		log.Infof("watching '%s'", path)

		if fswatcher.filter.check(path) {
			continue //skip exclude pattern
		}

		err = fswatcher.watcher.Add(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	<-fswatcher.done
	return nil
}

func expand(seed []string, expandFunc func(string) []string) []string {
	result := []string{}
	for _, str := range seed {
		result = append(result, expandFunc(str)...)
	}
	return result
}

func expandPath(path string) []string {
	matches, err := filepath.Glob(path)
	if err != nil {
		log.Warn(err)
		return []string{}
	}
	return matches
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
