package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bluemir/sentry/utils"
	"github.com/fsnotify/fsnotify"
)

type fsWatcher struct {
	config     *Config
	watchPaths map[string]bool
	watcher    *fsnotify.Watcher
	done       chan bool
}

func newFsWatcher(config *Config) *fsWatcher {

	fw := &fsWatcher{
		config:     config,
		watchPaths: map[string]bool{},
		watcher:    nil,
		done:       make(chan bool),
	}

	return fw
}

func (fswatcher *fsWatcher) handleEvent(callback func()) {
	for {
		select {
		case event := <-fswatcher.watcher.Events:
			if event.Op == fsnotify.Chmod {
				//osx bug?
				continue
			}
			log.Infof("Event: %s", event)

			if !notMatch(fswatcher.config.Exclude)(event.Name) {
				log.Info("Skip... Matching with exclude pattern")
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

	list := utils.NewStrings(fswatcher.config.WatchPaths...).
		Expand(glob).
		Expand(findAllDir).
		Filter(notMatch(fswatcher.config.Exclude))
	sort.Strings(list)

	for _, path := range list {
		fswatcher.appendFile(path)
	}

	go fswatcher.handleEvent(callback)

	callback() // do first command

	<-fswatcher.done
	return nil
}

func (fw *fsWatcher) appendFile(path string) {
	if fw.watchPaths[path] {
		return // already exist
	}

	if err := fw.watcher.Add(path); err != nil {
		log.Fatal(err)
	} else {
		fw.watchPaths[path] = true
		log.Infof("watching '%s'", path)
	}
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
func glob(path string) []string {
	result, err := filepath.Glob(path)
	if err != nil {
		log.Warn(err)
		return []string{path}
	}

	return result
}
func notMatch(patterns []string) func(string) bool {
	return func(path string) bool {
		for _, pattern := range patterns {
			log.Debugf("pattern: %s, path: %s", pattern, path)
			if ok, err := filepath.Match(pattern, path); err != nil {
				log.Warn(err)
			} else if ok {
				return false
			}
		}
		return true
	}
}
