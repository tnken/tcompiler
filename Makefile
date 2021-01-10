GOCMD=go
GOTEST=$(GOCMD) test

build:
			$(GOCMD) build main.go
test:
			$(GOTEST) ./...
clean:
			rm main
fmt:
			$(GOCMD) fmt ./...

.PHONY: test clean build fmt
