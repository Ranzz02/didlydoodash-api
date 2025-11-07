# Load .env if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

DB_URL = $(POSTGRES_DSN)
MIGRATIONS_DIR := internal/db/migrations
SQLC_CONFIG := internal/db/sqlc.yaml

ifeq ($(OS),Windows_NT)
	SQLC_OUT := internal\db\repository
else
	SQLC_OUT := internal/db/repository
endif

# --------------- migration commands ---------------

## Create a new migrations
# !ALERT: Needs to be run in a git bash terminal on windows
new-migration:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a name: make new-migration name=<name>"; \
	else \
		echo "📦 Creating migration: $(name)"; \
		migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name); \
	fi

# Migrate database to latest version
migrate-up:
	@echo "Applying all up migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Migrate database one version down
migrate-down:
	@echo "⬇️  Rolling back last migration..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Drop current migrations from database
# WARNING: Will delete all data in database
migrate-drop:
	echo "🧨 Dropping all tables..."; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f; \

# Migration status
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# ---------- sqlc ----------

sqlc-generate:
ifeq ($(OS),Windows_NT)
	@if exist $(SQLC_OUT) ( rmdir /S /Q $(SQLC_OUT) )
else
	rm -rf $(SQLC_OUT)
endif
	@echo "🧩 Regenerating fresh SQLC code..."
	sqlc generate -f $(SQLC_CONFIG)

sqlc-verify:
	@echo "🔍 Verifying SQLC config..."
	sqlc vet -f $(SQLC_CONFIG)

# Development workflow commands

# Build docker container
dev-build:
	@echo "🔨 Building development containers..."
	@docker compose -f docker-compose.dev.yml build

# Run development environment
dev-up:
	@echo "🚀 Starting development containers..."
	@docker compose -f docker-compose.dev.yml up -d

# Stop development environment
dev-down: 
	@echo "🧹 Stopping and removing development containers..."
	@docker compose -f docker-compose.dev.yml down

# Read logs from development environment
dev-logs:
	@docker compose -f docker-compose.dev.yml logs -f

# Start development environment
dev: dev-build dev-up
	@echo "✅ Dev environment is up and running!"