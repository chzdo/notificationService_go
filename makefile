build:
	@go build  -o=./dist/ ./cmd/

run:
	@nodemon --exec go run ./cmd/.  --ext .go