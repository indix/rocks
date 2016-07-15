APPNAME = rocks
VERSION = 0.1.0-dev

setup:
	glide install

build-all: build-mac build-linux

build:
	go build -ldflags "-X main.Version=${VERSION}" -v -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static' -X main.Version=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

ci:
	APPNAME=${APPNAME} bin/ci-run.sh

clean:
	rm -f ${APPNAME}
	rm -f ${APPNAME}-linux-amd64
	rm -f ${APPNAME}-darwin-amd64

all:
	setup
	build
	install

test:
	go test -v github.com/ind9/rocks
	go test -v github.com/ind9/rocks/cmd/backup
	go test -v github.com/ind9/rocks/cmd/restore
	go test -v github.com/ind9/rocks/cmd/statistics
	go test -v github.com/ind9/rocks/cmd/consistency

test-only:
	go test -v github.com/ind9/rocks/${name}

install: build
	sudo install -d /usr/local/bin
	sudo install -c ${APPNAME} /usr/local/bin/${APPNAME}

uninstall:
	sudo rm /usr/local/bin/${APPNAME}
