
GO = go
GOFLAGS = ""

all: build test

help:
	@echo "build     - go build"
	@echo "install   - go install"
	@echo "test      - go test"
	@echo "fmt       - go fmt"
	@echo "clean     - remove temp files"

build:
	$(GO) build -i $(GOFLAGS)

test:
	@$(GO) test -test.timeout 15s `go list ./... | grep -v '/vendor/'`
	@if [ $$? -eq 0 ] ; then \
		echo "All tests PASSED" ; \
	else \
		echo "Tests FAILED" ; \
	fi

fmt:
	gofmt -w `find . -type f -name '*.go' | grep -v vendor`

install:
	$(GO) install $(GOFLAGS)

clean:
	$(GO) clean
	rm -f adaptlb
