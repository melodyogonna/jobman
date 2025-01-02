
docker:
	docker compose up -d

test: docker
	go test ./...
