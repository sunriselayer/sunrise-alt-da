GITCOMMIT ?= $(shell git rev-parse HEAD)
GITDATE ?= $(shell git show -s --format='%ct')
VERSION ?= v0.0.0

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.Version=$(VERSION)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

install:
	@echo "--> Installing Sunrise-Alt-DA"
	@env GO111MODULE=on GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -v $(LDFLAGS) -o ./bin/da-server ./cmd/da-server
	@go install ./cmd/da-server
.PHONY: install

clean:
	@echo "--> Cleaning up"
	@rm bin/da-server 
.PHONY: clean

test:
	@echo "--> Running tests"
	@go test -v ./...
.PHONY: test
