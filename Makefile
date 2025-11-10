APP_NAME=knowledge-capsule-api
DOCKER_IMAGE=knowledge-capsule-api

.PHONY: run build clean test fmt vet docker-up docker-down hooks

run:
	@echo "🚀 Starting API (with live reload)..."
	air

build:
	@echo "🔨 Building binary..."
	go build -o server main.go

fmt:
	@echo "🧹 Formatting code..."
	go fmt ./...

vet:
	@echo "🔍 Running go vet..."
	go vet ./...

test:
	@echo "🧪 Running tests..."
	go test ./...

clean:
	@echo "🧼 Cleaning build files..."
	rm -f server tmp/server

docker-up:
	@echo "🐳 Starting Docker Compose..."
	docker-compose up --build

docker-down:
	@echo "🛑 Stopping containers..."
	docker-compose down

hooks:
	@echo "🔧 Installing git hooks..."
	lefthook install
