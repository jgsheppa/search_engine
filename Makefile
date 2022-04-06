run:
	docker run --publish 3000:3000 docker-build:latest
build:
	docker build --tag docker-build:latest .