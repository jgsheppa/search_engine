run:
	docker run --publish 3000:3000 docker-build:latest
build:
	docker build --tag docker-build:latest .

heroku-push:
	heroku container:push web --app bp-search-engine

heroku-release:
	heroku container:release web --app bp-search-engine
test:
	go test -v ./...