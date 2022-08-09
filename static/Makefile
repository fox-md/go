GIT_HASH=`git rev-parse --short HEAD || echo 0000`
BUILD_DATE=`date +%FT%T%z`

VERSION="1.0.0"
TAG="static"
PORT="8080"

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
	docker build --build-arg Version="${VERSION}" -t "${TAG}:${VERSION}" --no-cache=true .

run:
	docker run --rm -p "${PORT}"\:"${PORT}" "${TAG}"

stop:
	docker rm -f "${TAG}" || true

vet:
	go vet ./...

lint: vet
	golangci-lint run

clean:
	[ -f "${BIN}" ] && rm "${BIN}"

login:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/agu

tag: image
	docker tag "${TAG}:${VERSION}" public.ecr.aws/agu/misc:"${TAG}-${VERSION}"

push: tag login
	docker push public.ecr.aws/agu/misc:"${TAG}-${VERSION}"

pod:
	kubectl replace -f pod.yaml --force
