build: format test lint
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main .

test:
	go test ./...

format:
	go fmt ./...

clean:
	rm -rf ./bin

lint:
	golangci-lint run ./...

compose-up:
	docker-compose up -d redis debug && docker-compose up redis-populate && docker-compose up nylas-lite-gmail

compose-down:
	docker-compose down -v --remove-orphans

compose-exec:
	docker container exec -it gmail_debug_1 sh

docker-build:
	docker build -t nylas-lite-gmail -f docker/Dockerfile .