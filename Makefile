# Docker settings
PACKAGE_NAME := knowledge-capsule-api
DOCKER_USERNAME := shahadathhs
PACKAGE_VERSION := latest
APP_IMAGE := $(DOCKER_USERNAME)/$(PACKAGE_NAME):$(PACKAGE_VERSION)
COMPOSE_FILE := compose.yaml

# Go / build
BINARY_NAME := server
BUILD_DIR := tmp
BUILD_OUT := $(BUILD_DIR)/$(BINARY_NAME)
GO := go
GOFLAGS ?=
LDFLAGS ?=
BUILD_FLAGS ?=

# Local tools (install into ./.bin by default)
GOBIN ?= $(CURDIR)/.bin
BIN_DIR := $(GOBIN)

# Convenience
.PHONY: all help install hooks run build-local build push \
	clean fmt vet test up down restart logs containers volumes networks images

all: build-local

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  make all / build-local   Build local binary (default)"
	@echo "  make build               Build Docker image ($(APP_IMAGE))"
	@echo "  make push                Push Docker image to Docker Hub"
	@echo "  make up                  Start containers (docker compose -f $(COMPOSE_FILE) up)"
	@echo "  make down                Stop containers"
	@echo "  make restart             Restart containers"
	@echo "  make logs                Follow logs for the app container"
	@echo "  make run                 Run dev server with live reload (air)"
	@echo "  make hooks               Install git hooks (lefthook)"
	@echo "  make install             Install dev tools to $(GOBIN)"
	@echo "  make fmt                 Run go fmt ./..."
	@echo "  make vet                 Run go vet ./..."
	@echo "  make test                Run go test ./..."
	@echo "  make clean               Remove build artifacts and optionally docker images"
	@echo "  make containers/volumes/networks/images   Inspect docker compose resources"

# -------------------------
# Dev tools
# -------------------------
install:
	@echo "📦 Creating bin dir: $(GOBIN)"
	@mkdir -p "$(GOBIN)"
	@echo "⬇️  Installing dev tools into $(GOBIN)..."
	@GOBIN="$(GOBIN)" $(GO) install github.com/air-verse/air@latest
	@GOBIN="$(GOBIN)" $(GO) install github.com/evilmartians/lefthook@latest
	@echo "✅ Installed (air, lefthook) to $(GOBIN). Add $(GOBIN) to PATH to run them globally."

hooks: install
	@echo "🔧 Installing git hooks..."
	@$(GOBIN)/lefthook install || lefthook install

# run uses local .bin/air if present, otherwise falls back to system 'air'
run: install
	@echo "🚀 Starting API (with live reload)..."
	@$(GOBIN)/air || air

# -------------------------
# Build & test
# -------------------------
build-local:
	@echo "🔨 Building local binary -> $(BUILD_OUT)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_OUT) $(LDFLAGS) $(BUILD_FLAGS) ./main.go
	@echo "✅ Built: $(BUILD_OUT)"

fmt:
	@echo "🧹 Formatting code..."
	@$(GO) fmt ./...

vet:
	@echo "🔍 Running go vet..."
	@$(GO) vet ./...

test:
	@echo "🧪 Running tests..."
	@$(GO) test ./...

# -------------------------
# Docker / compose
# -------------------------
build:
	@echo "🐳 Building Docker image: $(APP_IMAGE)"
	@docker build -t $(APP_IMAGE) .

push: build
	@echo "📤 Pushing Docker image: $(APP_IMAGE)"
	@docker push $(APP_IMAGE)

up:
	@echo "🐳 Starting Docker Compose..."
	@docker compose -f $(COMPOSE_FILE) up --build

down:
	@echo "🛑 Stopping Docker Compose..."
	@docker compose -f $(COMPOSE_FILE) down

restart: down up

# follow logs, try common service names (compose service name or container name)
logs:
	@echo "📜 Following logs..."
	@docker compose -f $(COMPOSE_FILE) logs -f $(PACKAGE_NAME) || docker compose -f $(COMPOSE_FILE) logs -f $(APP_NAME)

containers:
	@docker compose -f $(COMPOSE_FILE) ps

volumes:
	@docker compose -f $(COMPOSE_FILE) volume ls

networks:
	@docker compose -f $(COMPOSE_FILE) network ls

images:
	@docker compose -f $(COMPOSE_FILE) images

# -------------------------
# Cleanup
# -------------------------
clean:
	@echo "🧼 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(BIN_DIR)
	@echo "🧽 Removing docker image (if exists): $(APP_IMAGE)"
	-@docker rmi -f $(APP_IMAGE) 2>/dev/null || true
	@echo "✅ Clean complete"
