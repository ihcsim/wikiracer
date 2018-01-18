package log

import (
	"io/ioutil"
	"sync"

	logging "gopkg.in/op/go-logging.v1"
)

const (
	logModule        = "racer"
	logLevelDefault  = logging.INFO
	logFormatDefault = "%{color}[%{level:.4s}] (%{shortpkg}/%{shortfile}) ▶ %{color:reset}%{message}"
)

var (
	log *logging.Logger

	once      = &sync.Once{}
	formatter = logging.MustStringFormatter(logFormatDefault)

	// QuietBackend discard all log messages
	QuietBackend = logging.AddModuleLevel(logging.NewLogBackend(ioutil.Discard, "", 0))
)

// Instance returns the logger instance.
func Instance() *logging.Logger {
	once.Do(func() {
		log = logging.MustGetLogger(logModule)
		logging.SetFormatter(formatter)
		logging.SetLevel(logLevelDefault, logModule)
	})
	return log
}
