package core

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

type Sentry struct {
	config  *Config
	shell   *shellCommander
	watcher *fsWatcher
}

func NewSentry(config *Config) *Sentry {
	log.Debugf("config : %v", config)

	return &Sentry{
		config:  config,
		shell:   newShellCommander(config.Shell),
		watcher: newFsWatcher(config),
	}
}

func (sentry *Sentry) Run() {
	log.Infof("execute command '%s'", sentry.config.Command)
	sentry.shell.exec(sentry.config.Command)

	sentry.registerSignal()

	err := sentry.watcher.watch(func() {
		if sentry.config.KillOnRestart {
			sentry.shell.stop()
		}
		log.Infof("execute command '%s'", sentry.config.Command)
		sentry.shell.exec(sentry.config.Command)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (sentry *Sentry) registerSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Warnf("recive %v", sig)
			sentry.shell.stop()
			log.Info("Exiting....")
			os.Exit(0)
		}
	}()
}
