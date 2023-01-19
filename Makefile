tidy:
	@go mod tidy -compat=1.17
generate:
	@go run cmd/main.go config
run: tidy
	@go run cmd/main.go api
run-race: tidy
	@go run -race cmd/main.go api
create-redis-stack:
	@docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest