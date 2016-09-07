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
		watcher: newFsWatcher(config),
	}
}

func (sentry *Sentry) Run() {
	err := sentry.watcher.watch(func() {
		if sentry.config.KillOnRestart {
			sentry.shell.stop()
		}
		log.Infof("excute command '%s'", sentry.config.Command)
		sentry.shell.exec(sentry.config.Command)
	})
	if err != nil {
		log.Fatal(err)
	}
}
