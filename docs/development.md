# Development Guide

## Local Development Setup

### Prerequisites

- Go 1.21+ ([Download](https://go.dev/dl/))
- Node.js 18+ ([Download](https://nodejs.org/))
- PostgreSQL 15+ ([Download](https://www.postgresql.org/download/))
- Git

### 1. Clone Repository

```bash
git clone https://github.com/dgentikian/rent-receipt-generator.git
cd rent-receipt-generator
```

### 2. Setup PostgreSQL Database

```bash
# Create database and user
sudo -u postgres psql

CREATE DATABASE rent_receipts;
CREATE USER rent_user WITH ENCRYPTED PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE rent_receipts TO rent_user;
ALTER DATABASE rent_receipts OWNER TO rent_user;
\q

# Initialize schema
psql -U rent_user -d rent_receipts -f database/schema.sql
```

### 3. Configure Environment

Create `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=rent_user
DB_PASSWORD=dev_password
DB_NAME=rent_receipts
DB_SSLMODE=disable

JWT_SECRET=your_development_jwt_secret_min_32_chars
JWT_EXPIRY=24h

UPLOADS_DIR=./uploads
FRONTEND_URL=http://localhost:3000
BACKEND_URL=http://localhost:8080

CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 4. Backend Setup

```bash
cd backend

# Download dependencies
go mod download

# Run backend (with auto-reload using air - optional)
go run cmd/server/main.go
```

The backend will be available at `http://localhost:8080`

#### Auto-reload with Air (Optional)

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air
```

### 5. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Project Structure

```
rent-receipt-generator/
├── backend/              # Go backend
│   ├── cmd/             # Application entrypoints
│   ├── internal/        # Private application code
│   │   ├── api/         # HTTP handlers and routes
│   │   ├── models/      # Data models
│   │   ├── repository/  # Database access layer
│   │   ├── service/     # Business logic
│   │   ├── config/      # Configuration
│   │   └── database/    # Database connection
│   └── pkg/             # Public libraries
│
├── frontend/            # React frontend
│   └── src/
│       ├── components/  # React components
│       ├── pages/       # Page components
│       ├── services/    # API services
│       ├── context/     # React context
│       ├── types/       # TypeScript types
│       └── styles/      # CSS styles
│
├── database/            # Database schema
├── deployment/          # Deployment configs
└── docs/               # Documentation
```

## Making Changes

### Backend Changes

1. Make changes to Go files
2. Code is automatically compiled when running `go run`
3. Or use `air` for hot-reload

### Frontend Changes

1. Make changes to React/TypeScript files
2. Vite automatically hot-reloads changes
3. Check browser console for errors

### Database Changes

1. Update `database/schema.sql`
2. Apply changes:
   ```bash
   psql -U rent_user -d rent_receipts -f database/schema.sql
   ```

## Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm test
```

## Building for Production

### Build Everything

```bash
make build
```

### Build Backend Only

```bash
make build-backend
```

### Build Frontend Only

```bash
make build-frontend
```

## Common Tasks

### Create New API Endpoint

1. Add model in `backend/internal/models/`
2. Add repository methods in `backend/internal/repository/`
3. Add service logic in `backend/internal/service/`
4. Add handler in `backend/internal/api/handlers/`
5. Register route in `backend/internal/api/router.go`

### Create New Frontend Page

1. Add page component in `frontend/src/pages/`
2. Add route in `frontend/src/App.tsx`
3. Add menu item in `frontend/src/components/layout/Sidebar.tsx`

## Code Style

### Backend (Go)

Follow standard Go conventions:
- Use `gofmt` for formatting
- Use `golint` for linting
- Follow [Effective Go](https://go.dev/doc/effective_go)

### Frontend (TypeScript/React)

- Use Prettier for formatting
- Use ESLint for linting
- Follow React best practices
- Use functional components with hooks

## Debugging

### Backend Debugging

Using VS Code:
1. Install Go extension
2. Add launch configuration:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/cmd/server"
    }
  ]
}
```

### Frontend Debugging

Use React DevTools browser extension

### Database Debugging

```bash
# Connect to database
psql -U rent_user -d rent_receipts

# List tables
\dt

# View table structure
\d receipts

# Query data
SELECT * FROM receipts;
```

## API Testing

### Using curl

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Get properties
curl http://localhost:8080/api/v1/properties \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Using Postman

Import the API endpoints from `docs/api.md`

## Environment Variables

Required:
- `DB_PASSWORD`: Database password
- `JWT_SECRET`: JWT secret (min 32 chars)

Optional:
- `PORT`: Backend port (default: 8080)
- `ENV`: Environment (development/production)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_NAME`
- `UPLOADS_DIR`: Upload directory path
- `CORS_ALLOWED_ORIGINS`: CORS origins

## Troubleshooting

### "Connection refused" errors

- Check if PostgreSQL is running: `sudo systemctl status postgresql`
- Verify database credentials in `.env`
- Check if backend is running on correct port

### Frontend can't connect to backend

- Check CORS settings in backend
- Verify `FRONTEND_URL` in `.env`
- Check proxy settings in `vite.config.ts`

### Module not found errors

Backend:
```bash
cd backend && go mod tidy
```

Frontend:
```bash
cd frontend && npm install
```

## Contributing

1. Create a feature branch
2. Make changes
3. Test locally
4. Commit with clear message
5. Push and create pull request

## Resources

- [Go Documentation](https://go.dev/doc/)
- [React Documentation](https://react.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Vite Documentation](https://vitejs.dev/)
