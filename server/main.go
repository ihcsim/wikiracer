package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ihcsim/wikiracer"
	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/crawler"
	"github.com/ihcsim/wikiracer/internal/validator"
	"github.com/ihcsim/wikiracer/internal/wiki/wikipedia"
	"github.com/ihcsim/wikiracer/log"

	_ "net/http/pprof"
)

const (
	queryParameterOrigin      = "origin"
	queryParameterDestination = "destination"

	serverPort = "8080"
	pprofPort  = "6060"
)

var timeout = 180 * time.Second

func main() {
	go func() {
		log.Instance().Infof("Starting profiling server at port %s...", pprofPort)
		if err := http.ListenAndServe(":"+pprofPort, nil); err != nil {
			log.Instance().Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)
	go func() {
		if signal := <-interrupt; signal == os.Interrupt {
			log.Instance().Info("Stopping server...")
			os.Exit(0)
		}
	}()

	log.Instance().Infof("Starting up server at port %s...", serverPort)
	http.HandleFunc("/wikiracer", timedFindPath)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Instance().Fatal(err)
	}
}

func timedFindPath(w http.ResponseWriter, req *http.Request) {
	wiki, err := wikipedia.NewClient()
	if err != nil {
		response(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	var (
		crawler   = crawler.NewForward(wiki)
		validator = validator.NewInputValidator(wiki)
	)

	racer := wikiracer.New(crawler, validator)

	var (
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		origin      = req.URL.Query().Get(queryParameterOrigin)
		destination = req.URL.Query().Get(queryParameterDestination)
	)
	defer cancel()

	log.Instance().Infof("%q -> %q: Starting...", origin, destination)
	if origin == "" || destination == "" {
		err := errors.InvalidEmptyInput{origin, destination}.Error()
		log.Instance().Errorf("%q -> %q: Failed. Reason: %q", origin, destination, err)
		response(w, http.StatusBadRequest, []byte(err))
		return
	}

	result := racer.TimedFindPath(ctx, origin, destination)
	if result.Err != nil {
		err := result.Err.Error()
		log.Instance().Errorf("%q -> %q: Failed. Reason: %q", origin, destination, err)
		response(w, http.StatusInternalServerError, []byte(err))
		return
	}

	log.Instance().Infof("%q -> %q: SUCCESS. %s", origin, destination, result)
	response(w, http.StatusOK, []byte(fmt.Sprintf("%s", result)))
}

func response(w http.ResponseWriter, status int, content []byte) {
	w.WriteHeader(status)
	w.Write(content)
}
