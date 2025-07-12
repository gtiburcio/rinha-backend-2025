run-local:
	go run ./src/main.go local

build-image:
	docker buildx build --platform linux/amd64 -t guilhermetiburcio/rinha-2025-gtiburcio .

start-compose:
	docker-compose down && docker-compose up
