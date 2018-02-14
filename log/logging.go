package log

import (
	"io/ioutil"
	"os"
	"sync"

	logging "gopkg.in/op/go-logging.v1"
)

const (
	logModule        = "racer"
	logLevelDefault  = logging.INFO
	logLevelEnvVar   = "WIKIRACER_LOG_LEVEL"
	logFormatDefault = "%{color}[%{level:.4s}] - %{id} (%{shortpkg}/%{shortfile}) ▶ %{color:reset}%{message}"
)

var (
	log    *logging.Logger
	logMux sync.Mutex

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

	// read the log level from env var, if specified.
	// otherwise, use the default log level.
	logLevel, exist := os.LookupEnv(logLevelEnvVar)
	if exist {
		logMux.Lock()
		if err := setLogLevel(logLevel); err != nil {
			log.Warning("Can't change log level. ", err.Error())
		}
		logMux.Unlock()
	}

	return log
}

func setLogLevel(l string) error {
	logLevel, err := logging.LogLevel(l)
	if err != nil {
		return err
	}
	logging.SetLevel(logLevel, logModule)
	return nil
}
