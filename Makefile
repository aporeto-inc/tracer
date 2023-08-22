export GO111MODULE = on
export GOPRIVATE = go.aporeto.io,github.com/aporeto-inc

lint:
	golangci-lint run \
		--timeout=5m \
		--disable-all \
		--exclude-use-default=false \
		--exclude=package-comments \
		--exclude=unused-parameter \
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
		--enable=unparam \
		--enable=gosimple \
		--enable=nilerr \
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

# updating some special components
#  go get go4.org/unsafe/assume-no-moving-gc@latest
#  go get -u github.com/grafana/loki@main
