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
	delay   *delayer
}

func NewSentry(config *Config) *Sentry {
	log.Debugf("config : %v", config)

	return &Sentry{
		config:  config,
		shell:   newShellCommander(config.Shell),
		watcher: newFsWatcher(config),
		delay:   newDelayer(config.Delay),
	}
}

func (sentry *Sentry) Run() {
	sentry.registerSignal()

	log.Infof("execute command '%s'", sentry.config.Command)
	sentry.shell.exec(sentry.config.Command)

	err := sentry.watcher.watch(func() {
		sentry.delay.Do(sentry.restartCommand)
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
			sentry.watcher.close()
			os.Exit(0)
		}
	}()
}

func (sentry *Sentry) restartCommand() {
	if sentry.config.KillOnRestart {
		sentry.shell.stop()
	}
	log.Infof("execute command '%s'", sentry.config.Command)
	sentry.shell.exec(sentry.config.Command)
}
