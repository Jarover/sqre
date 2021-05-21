APP?=sqre
RELEASE?=$(shell python version.py get)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell powershell get-date -format "{yyyy-mm-dd_HH:mm:ss}")
PROJECT?=github.com/Jarover/sqre

clean:
	del ${APP}
	del ${APP}.exe

build:	clean
	
	python version.py inc-patch
	go build \
                -ldflags "-s -w -X ${PROJECT}/readconfig.Release=${RELEASE} \
                -X ${PROJECT}/readconfig.Commit=${COMMIT} -X ${PROJECT}/readconfig.BuildTime=${BUILD_TIME}" \
                -o ${APP}

run:	build
	./${APP} -f dev.json

test:
	go test -v -race ./...