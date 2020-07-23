PATH := ${PWD}/bin:${PATH}
export PATH

.DEFAULT_GOAL := build

.PHONY: clean
clean:
	rm -rf ./bin/*

.PHONY: build
build:
	go build -v -o ./bin/organize-repo-folder ./cmd/organize-repo-folder
