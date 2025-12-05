GOPATH = $(shell go env GOPATH)
VERSION = $(shell git describe --tags )
TODAY = $(shell date -u +'%Y-%m-%d')
build:
	go build -ldflags "-X main.version=${VERSION} -X main.date=${TODAY}" -o activitymonitor main.go

run: 
	go run main.go

install-build-tools-darwin:
		@if [ "$(shell uname)" = "Darwin" ]; then \
				echo "macOS build"; \
				brew install filosottile/musl-cross/musl-cross; \
				brew install mingw-w64; \
				brew install --cask xquartz; \
		else \
				echo "Other OS build"; \
		fi


install-build-tools:
	go install fyne.io/tools/cmd/fyne@latest
	go install github.com/fyne-io/fyne-cross@latest


clean:
	rm -f activitymonitor

install: build
	install -m 755 activitymonitor ${GOPATH}/bin/activitymonitor

package: install-build-tools install-build-tools-darwin
	fyne package -os darwin -icon pkg/gui/gopher.jpeg
	env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc"	fyne package -os windows -icon pkg/gui/gopher.jpeg

