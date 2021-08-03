APP?=sqre
RELEASE?=$(shell python version.py get)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell powershell get-date -format "{yyyy-mm-dd_HH:mm:ss}")
PROJECT?=github.com/Jarover/sqre

clean:
	rm -f ${APP}
	rm -f ${APP}.exe

build:	clean
	
	python version_json.py inc-patch ${COMMIT}
	GOOS=linux go build -o ${APP}

run:	build
	./${APP} -f dev.json

test:
	go test -v -race ./...