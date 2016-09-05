package core

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

type Tick struct {
	config *Config
	shell  *shellCommander
}

func NewTick(config *Config) *Tick {
	log.Debugf("path : %v", config)

	return &Tick{
		config: config,
		shell:  newShellCommander(config.Shell),
	}
}

func (tick *Tick) Run() {
	err := tick.watch(func() {
		if tick.config.KillOnRestart {
			tick.shell.stop()
		}
		tick.shell.exec(tick.config.Command)
	})
	if err != nil {
		log.Fatal(err)
	}
}
func (tick *Tick) watch(callback func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Infof("event: %s", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debugln("modified file:", event.Name)
				}
				callback()
			case err := <-watcher.Errors:
				log.Debugln("error:", err)
			}
		}
	}()

	list := findAllDir(tick.config.Path)
	log.Infof("watching '%s'", tick.config.Path)
	err = watcher.Add(tick.config.Path)
	for _, path := range list {
		log.Infof("watching '%s'", path)
		err = watcher.Add(path)
	}

	if err != nil {
		log.Fatal(err)
		return err
	}
	<-done
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
