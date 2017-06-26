
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
	@time $(GO) test -short -race -test.timeout 15s `go list ./... | grep -v '/vendor/'`
	@if [ $$? -eq 0 ] ; then \
		echo "All tests PASSED" ; \
	else \
		echo "Tests FAILED" ; \
	fi

fulltest:
	@time $(GO) test -race -test.timeout 15s `go list ./... | grep -v '/vendor/'`
	@if [ $$? -eq 0 ] ; then \
		echo "All tests PASSED" ; \
	else \
		echo "Tests FAILED" ; \
	fi

testloop:
	while $(GO) test -race -test.timeout 15s `go list ./... | grep -v '/vendor/'`; do :; done

testloopshort:
	while $(GO) test -race -short -test.timeout 15s `go list ./... | grep -v '/vendor/'`; do :; done

fmt:
	gofmt -w `find . -type f -name '*.go' | grep -v vendor`

commit: fmt fulltest
	git commit -a

install:
	$(GO) install $(GOFLAGS)

clean:
	$(GO) clean
	rm -f adaptlb
