u := $(if $(update),-u)

.PHONY: deps
deps:
	go get ${u}
	go mod tidy

.PHONY: test
test:
	go test

.PHONY: prepare-release
prepare-release: devel-deps
	go mod tidy
	git update-index --add --remove -- go.mod go.sum
