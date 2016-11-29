package core

import (
	log "github.com/Sirupsen/logrus"
)

type SentryFormatter struct {
	log.TextFormatter
}

func NewSentryFormatter() *SentryFormatter {
	return &SentryFormatter{
		log.TextFormatter{},
	}
}

func (formatter *SentryFormatter) Format(entry *log.Entry) ([]byte, error) {
	if buf, err := formatter.TextFormatter.Format(entry); err != nil {
		return nil, err
	} else {
		buf = append([]byte("[SENTRY] "), buf...)
		return buf, nil
	}
}
