run:
	go run main.go

docker-run:
	docker run --publish 3000:3000 docker-build:latest

build:
	docker build --tag docker-build:latest .

test:
	go test -v ./...

docs:
	swag init -g ./main.go -o ./swagger

test_data:
	curl -X POST -H "Content-Type: application/json" -d @./nba_players.json http://localhost:3001/api/document

prod_data:
	curl -X POST -H "Content-Type: application/json" -d @./nba_players.json \
 https://bp-search-engine.herokuapp.com/api/document

# Commands for Heroku deployment with Heroku container stack
# Note: You'll need to build the Docker image before pushing
# it to Heroku
heroku-login:
	heroku container:login

heroku-push:
	heroku container:push web --app bp-search-engine

heroku-release:
	heroku container:release web --app bp-search-engine