.PHONY: test

test:
	govendor test +local -cover -race

server:
	go run server/main.go

pprof:
	go run pprof/main.go
