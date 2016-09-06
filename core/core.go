package core

import (
	log "github.com/Sirupsen/logrus"
)

type Tick struct {
	config  *Config
	shell   *shellCommander
	watcher *fsWatcher
}

func NewTick(config *Config) *Tick {
	log.Debugf("path : %v", config)

	return &Tick{
		config:  config,
		shell:   newShellCommander(config.Shell),
		watcher: newFsWatcher(config.Path),
	}
}

func (tick *Tick) Run() {
	err := tick.watcher.watch(func() {
		if tick.config.KillOnRestart {
			tick.shell.stop()
		}
		tick.shell.exec(tick.config.Command)
	})
	if err != nil {
		log.Fatal(err)
	}
}
