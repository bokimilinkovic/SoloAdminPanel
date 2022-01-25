package log

import (
	"errors"

	glog "github.com/labstack/gommon/log"
	log "github.com/sirupsen/logrus"
)

const (
	logTypeJSON = "json"
	logTypeText = "text"
)

type Logger struct {
	*log.Logger
}

// New creates a new logger
func New(cfg Config) (*Logger, error) {
	engine := log.New()

	switch cfg.Type {
	case logTypeJSON:
		engine.Formatter = &log.JSONFormatter{}
	case logTypeText:
		engine.Formatter = &log.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
		}
	default:
		return nil, errors.New("bad log type provided ; supported log types are: [json,text]; got: " + cfg.Type)
	}

	lvl, err := log.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.New("bad log level provided:" + err.Error())
	}

	engine.SetLevel(lvl)
	return &Logger{engine}, nil
}

func (l *Logger) Debugj(j glog.JSON) {

}

// TestLogger is used in unit & integration tests.
func TestLogger() *Logger {
	l, _ := New(Config{
		Level: "DEBUG",
		Type:  logTypeText,
	})

	return l
}
