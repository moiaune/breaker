run: build
	@./bin/breaker

build:
	@go build -o ./bin/breaker .

