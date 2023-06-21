BINARY_NAME=rss
.PHONY: default build run clean 


default: build

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

run: build
	env
	./${BINARY_NAME}-linux

clean:
	go clean
	rm ${BINARY_NAME}-linux