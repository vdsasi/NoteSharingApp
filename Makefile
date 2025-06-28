.PHONY: run build clean docker-up docker-down

run:
	go run main.go

build:
	go build -o bin/notes-app main.go

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

install:
	go mod tidy

dev: docker-up
	sleep 5
	go run main.go
