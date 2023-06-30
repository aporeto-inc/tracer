export GO111MODULE = on
export GOPRIVATE = go.aporeto.io,github.com/aporeto-inc

lint:
	golangci-lint run \
		--deadline=5m \
		--disable-all \
		--exclude-use-default=false \
		--enable=errcheck \
		--enable=goimports \
		--enable=ineffassign \
		--enable=revive \
		--enable=unused \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=misspell \
		--enable=prealloc \
		--enable=nakedret \
		--enable=typecheck \
		./...

test: lint
	go test ./... -race -cover -covermode=atomic -coverprofile=unit_coverage.cov

build_linux: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

build: test
	go build

# tag it first, then run...
release:
	unset GITLAB_TOKEN && goreleaser check && goreleaser release --clean

.PHONY: build
