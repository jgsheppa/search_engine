# A Search Engine Built With Redisearch and Go
This is a microservice written in Go, which allows developers to index and search for documents with Redis. 

## Demo
This microservice uses the player list from [basketball reference](https://www.basketball-reference.com/) as a dataset
to show the functionality of the search capabilities of this Redisearch/Go API. The UI for this API can be found 
[here](https://search-bar-jade.vercel.app/). Test it out by searching for your favorite basketball players.

## Setup
The easiest way to get up and running is to create a `config.yaml` file with your environment variables, run 
`redis-stack` with `docker-compose`, and start the Go API with `make run`. The values from the `config-example.yaml`
should work if you copy and paste them into a `config.yaml` file, but feel free to customize them for your needs.

### Makefile commands
The `makefile` contains the main commands you will need to get up and running. I prefer to run `make docker-run` and 
then `make run` to start docker and the Go server separately. This allows me to play with the Go code without shutting 
down Docker.

### Redisearch Schema
In order to customize this API for your own dataset, you will need to change the schema in `models/schema.go`. The 
`Document` type, the variables denoting the schema values, and the schema in `models.CreateSchema` will need to be 
updated if you use your own dataset. In the current example, I only use `redisearch.NewTextFieldOptions` but Redisearch 
also supports numeric and geo fields. Check 
[the docs](https://pkg.go.dev/github.com/RediSearch/redisearch-go/redisearch@v1.1.1#Field) for more information.

