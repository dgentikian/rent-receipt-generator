.PHONY: help build build-backend build-frontend deploy clean test run-backend run-frontend dev

help:
	@echo "Available commands:"
	@echo "  make build           - Build both backend and frontend"
	@echo "  make build-backend   - Build Go backend"
	@echo "  make build-frontend  - Build React frontend"
	@echo "  make run-backend     - Run backend in development mode"
	@echo "  make run-frontend    - Run frontend in development mode"
	@echo "  make dev             - Run both backend and frontend in dev mode"
	@echo "  make deploy          - Deploy to production"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make test            - Run all tests"

build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd backend && chmod +x build.sh && ./build.sh

build-frontend:
	@echo "Building frontend..."
	cd frontend && chmod +x build.sh && ./build.sh

run-backend:
	@echo "Running backend in development mode..."
	cd backend && go run cmd/server/main.go

run-frontend:
	@echo "Running frontend in development mode..."
	cd frontend && npm run dev

dev:
	@echo "Starting development environment..."
	@echo "Run 'make run-backend' in one terminal and 'make run-frontend' in another"

deploy:
	@echo "Deploying application..."
	./deployment/scripts/deploy.sh

clean:
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	@echo "Cleaned build artifacts"

test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm test
