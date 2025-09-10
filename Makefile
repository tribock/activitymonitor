GOPATH = $(shell go env GOPATH)
build:
	go build -o activitymonitor main.go

run: 
	go run main.go

clean:
	rm -f activitymonitor

install: build
	install -m 755 activitymonitor ${GOPATH}/bin/activitymonitor
