# Run server
start:
	@echo "Starting server"
	@go run .

# Run tests
test:
	@go test -v ./...

# Start docker container
docker-start:
	@docker-compose up -d

# Stop docker container
docker-stop:
	@docker-compose down -v

# Run tests
docker-test:
	@docker-compose exec app go test -v ./...