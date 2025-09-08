GOPATH = $(shell go env GOPATH)
build:
	go build -o mousemover main.go

run: build
	./mousemover

clean:
	rm -f mousemover

install: build
	install -m 755 mousemover ${GOPATH}/bin/mousemover
