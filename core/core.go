package core

import (
	log "github.com/Sirupsen/logrus"
)

type Sentry struct {
	config  *Config
	shell   *shellCommander
	watcher *fsWatcher
}

func NewSentry(config *Config) *Sentry {
	log.Debugf("path : %v", config)

	return &Sentry{
		config:  config,
		shell:   newShellCommander(config.Shell),
		watcher: newFsWatcher(config.Path),
	}
}

func (sentry *Sentry) Run() {
	err := sentry.watcher.watch(func() {
		if sentry.config.KillOnRestart {
			sentry.shell.stop()
		}
		sentry.shell.exec(sentry.config.Command)
	})
	if err != nil {
		log.Fatal(err)
	}
}
