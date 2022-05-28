# Start the go server
run:
	go run main.go

# Start redis and postgres with docker-compose
docker-run:
	docker-compose up

# Build a docker image from the Dockerfile
build:
	docker build --tag docker-build:latest .

# Run the go tests
test:
	go test -v ./...

# If you update the swagger documentation, run this command to update the swagger UI
docs:
	swag init -g ./main.go -o ./swagger

# You will need to add the Bearer token to the header for this to work
test_data:
	curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer _" -d @./nba_players.json http://localhost:3001/api/document


