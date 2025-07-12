run-local:
	go run ./src/main.go local

buid-image:
	docker buildx build --platform linux/amd64 -t rinha-2025-gtiburcio .