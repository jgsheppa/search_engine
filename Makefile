run:
	docker run --publish 3000:3000 docker-build:latest
build:
	docker build --tag docker-build:latest .
test:
	go test -v ./...

heroku-login:
	heroku container:login
	
heroku-push:
	heroku container:push web --app bp-search-engine

heroku-release:
	heroku container:release web --app bp-search-engine

docs:
	swag init -g ./main.go -o ./swagger