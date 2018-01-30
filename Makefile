.PHONY: test

test:
	govendor test +local -cover -race
