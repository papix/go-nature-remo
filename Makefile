DEP ?= $(GOPATH)/bin/dep

deps:
	$(DEP) ensure -vendor-only

test:
	go test -v ./...
