package log

import (
	"os"
	"sync"
	"testing"

	logging "gopkg.in/op/go-logging.v1"
)

func TestInstance(t *testing.T) {
	t.Run("Log Level", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			log := Instance()
			if !log.IsEnabledFor(logLevelDefault) {
				t.Error("Log level mismatched. Expected logging level to be ", logLevelDefault)
			}
		})

		t.Run("Env Var", func(t *testing.T) {
			defer os.Unsetenv(logLevelEnvVar)

			logLevels := []string{
				"CRITICAL",
				"ERROR",
				"WARNING",
				"NOTICE",
				"INFO",
				"DEBUG",
			}

			for _, level := range logLevels {
				if err := os.Setenv(logLevelEnvVar, level); err != nil {
					t.Fatal(err)
				}
				log := Instance()

				expected, err := logging.LogLevel(level)
				if err != nil {
					t.Fatal(err)
				}

				if !log.IsEnabledFor(expected) {
					t.Error("Log level mismatched. Expected logging level to be ", expected)
				}
			}
		})

		t.Run("No race", func(t *testing.T) {
			os.Setenv(logLevelEnvVar, "WARNING")
			defer os.Unsetenv(logLevelEnvVar)

			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				log := Instance()
				if log == nil {
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				log := Instance()
				if log == nil {
				}
			}()

			wg.Wait()
		})
	})

}
