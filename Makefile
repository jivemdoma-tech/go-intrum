.DEFAULT_GOAL:=test

.PHONY: new-test
new-test:
	mkdir .test && cd .test && touch main.go

.PHONY: run
run:
	cd .test && go run .
