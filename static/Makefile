GIT_HASH=`git rev-parse --short HEAD || echo 0000`
BUILD_DATE=`date +%FT%T%z`

TAG="static"

BIN=main
SRC=main.go

LDFLAGS=-w -s -linkmode external -extldflags -static -X main.GitHash=${GIT_HASH} -X main.BuildDate=${BUILD_DATE}

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

.PHONY: clean build

build: vet
        go build -ldflags "${LDFLAGS}" -o "${BIN}" "${SRC}"

image: vet
        docker build -t "${TAG}" --no-cache=true .

run:
        docker run --rm -p 443\:8443 "${TAG}"

stop:
        docker rm -f "${TAG}" || true

vet:
        go vet ./...

lint: vet
        golangci-lint run

clean:
        [ -f "${BIN}" ] && rm "${BIN}"
