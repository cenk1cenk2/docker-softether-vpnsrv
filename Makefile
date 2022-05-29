GO_VERSION=1.18

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run .
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_VENDOR=$(GO_CMD) mod vendor

GO_OPTION_C=0

install:
	$(GO_VENDOR)

update:
	$(GO_GET) -u all
	$(GO_VENDOR)
	$(GO_CMD) mod tidy -compat=$(GO_VERSION)

tidy:
	$(GO_CMD) mod tidy -compat=$(GO_VERSION)
	$(GO_VENDOR)

BINARY_FOLDER=dist
BINARY_NAME=pipe

all: test build

test:
	CGO_ENABLED=$(GO_OPTION_C) $(GO_TEST) -v -p 1 ./...

clean:
	$(GO_CLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)*

.PHONY: build

# Cross compilation

build: build-linux-amd64

build-linux-amd64:
	CGO_ENABLED=$(GO_OPTION_C) GOOS=linux GOARCH=amd64 $(GO_BUILD) -mod=readonly -o $(BINARY_FOLDER)/$(BINARY_NAME)

dev:
	CGO_ENABLED=$(GO_OPTION_C) $(GO_RUN) --log-level debug $(ARGS)

docs:
	CGO_ENABLED=$(GO_OPTION_C) $(GO_RUN) --log-level debug docs

help:
	CGO_ENABLED=$(GO_OPTION_C) $(GO_RUN) --help
