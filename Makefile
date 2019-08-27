# Package variables
URL := "https://github.com/maxmouchet/goasn"
NAME := goasn
LICENSE := MIT
VERSION := 0.0.1
MAINTAINER := "<max@maxmouchet.com>"
DESCRIPTION := "Fast offline lookup of IP addresses to IXP and AS numbers from RIB archives."

FPMFLAGS := --description ${DESCRIPTION} --license ${LICENSE} --maintainer ${MAINTAINER} --url ${URL} -n ${NAME} -v ${VERSION}

.PHONY: all
all: goasn

.PHONY: clean
clean:
	rm -f goasn goasn.exe *.deb *.rpm *.pkg

.PHONY: release
release:
	GOARCH=amd64 GOOS=linux go build -o goasn main.go
	fpm -s dir -t deb -f --prefix /usr/bin ${FPMFLAGS} goasn
	fpm -s dir -t rpm -f --prefix /usr/bin ${FPMFLAGS} goasn
	GOARCH=amd64 GOOS=windows go build -o goasn.exe main.go

goasn: $(shell find . -name '*.go')
	go build -o goasn main.go