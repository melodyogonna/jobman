
docker:
	docker compose up -d
clean:
	docker compose down

test: docker
	go test ./...
