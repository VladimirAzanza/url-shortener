.PHONY: test, cover, docs

test:
	go test -v ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

docs:
	rm -rf docs || true
	swag init -g cmd/shortener/main.go
