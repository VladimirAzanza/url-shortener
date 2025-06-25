.PHONY: install_tools, test, cover, docs, gen

test:
	go test -v ./...

cover:
	CGO_ENABLED=1 go test -short -count=1 -race -coverpkg=./internal/controller/...,./internal/services/...,./internal/middleware/... -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

docs:
	rm -rf docs || true
	swag init -g cmd/shortener/main.go

gen:
	rm -rf mocks || true
	mockgen -source=internal/services/interfaces.go -destination=mocks/services_mock.go -package=mocks
	mockgen -source=internal/repo/interfaces.go -destination=mocks/repo_mock.go -package=mocks

install_tools:
	@echo "Installing necessary tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install go.uber.org/mock/mockgen@latest
	@echo "Tools installed successfully."