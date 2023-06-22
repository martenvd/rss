BINARY_NAME=rss
.PHONY: default build run clean 


default: build

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

run: build
	docker-compose -f ./build/docker-compose.yml up -d
	./${BINARY_NAME}-linux

clean:
	go clean
	docker-compose -f ./build/docker-compose.yml down
	rm ${BINARY_NAME}-linux