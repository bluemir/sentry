package core

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type delayer struct {
	timer *time.Timer
	delay time.Duration

	event chan func()
}

func newDelayer(delay int32) *delayer {
	this := &delayer{
		time.NewTimer(time.Duration(0)),
		time.Duration(delay) * time.Millisecond,
		make(chan func()),
	}

	go this.handleEvent()

	return this
}
func (this *delayer) handleEvent() {
	var callback func()
	for {
		select {
		case <-this.timer.C:
			// do function
			if callback != nil {
				log.Debug("Runnging callback...")
				callback()
			} else {
				log.Debug("Callback is nil")
			}
		case cb := <-this.event:
			log.Debugf("New event arrived. Wating %s...", this.delay.String())
			//register callback
			callback = cb
			//reset timer
			this.resetTimer()
		}
	}
}
func (this *delayer) resetTimer() {
	//if !this.timer.Stop() {
	//	<-this.timer.C
	//}
	this.timer.Reset(this.delay)
}
func (this *delayer) Do(callback func()) {
	this.event <- callback
}
