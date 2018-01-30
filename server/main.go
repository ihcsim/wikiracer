package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ihcsim/wikiracer"
	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/crawler"
	"github.com/ihcsim/wikiracer/internal/validator"
	"github.com/ihcsim/wikiracer/log"
	"github.com/ihcsim/wikiracer/test"
)

const (
	queryParameterOrigin      = "origin"
	queryParameterDestination = "destination"

	serverPort = "8080"
)

var timeout = time.Second

func main() {
	log.Instance().Infof("Starting up server at port %s...", serverPort)
	http.HandleFunc("/", timedFindPath)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Instance().Fatal(err)
	}
}

func timedFindPath(w http.ResponseWriter, req *http.Request) {
	var (
		wiki      = test.NewMockWiki()
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
