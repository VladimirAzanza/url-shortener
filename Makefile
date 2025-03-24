.PHONY: docs, gen, test, cover

test: gen
	go test -v ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
