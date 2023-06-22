BINARY_NAME=rss
.PHONY: default build run clean 


default: build

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME}-linux main.go
	docker build -t rss:latest ./build

run: build
	docker-compose -f ./build/docker-compose.yml up -d

clean:
	go clean
	docker-compose -f ./build/docker-compose.yml down
	rm ./build/${BINARY_NAME}-linux