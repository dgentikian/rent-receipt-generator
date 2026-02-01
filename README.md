# Rent Receipt Generator

A full-stack web application for generating and managing rent receipts with PDF export capabilities.

## Features

- 🏠 Property management
- 👤 Tenant management
- 📄 Automatic receipt generation with PDF export
- 📊 Receipt history and tracking
- ✍️ Digital signature support
- 🔐 Secure authentication
- 💾 PostgreSQL database storage

## Tech Stack

**Backend:**
- Go 1.21+
- PostgreSQL 15+
- Gin Web Framework
- JWT Authentication
- gofpdf for PDF generation

**Frontend:**
- React 18
- TypeScript
- Vite
- TailwindCSS
- React Router
- React Query
- Axios

## Project Structure

```
rent-receipt-generator/
├── backend/          # Go backend API
├── frontend/         # React frontend
├── database/         # Database schema and migrations
├── deployment/       # Deployment scripts and configs
└── docs/            # Documentation
```

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 15 or higher
- Git

## Development Setup

### 1. Clone the repository

```bash
git clone https://github.com/dgentikian/rent-receipt-generator.git
cd rent-receipt-generator
```

### 2. Setup environment variables

**For Local Development:**
```bash
# Copy the local development template
cp env-local.example .env.local
# Edit .env.local with your local configuration
```

**For Production:**
```bash
# Copy the production template
cp env-prod.example .env.prod
# Edit .env.prod with your production values
# IMPORTANT: Change all passwords and secrets!
```

The backend will automatically load `.env` or you can specify which file to use.

### 3. Setup PostgreSQL database

```bash
# Create database and user
sudo -u postgres psql
CREATE DATABASE rent_receipts;
CREATE USER rent_user WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE rent_receipts TO rent_user;
\q

# Run migrations
psql -U rent_user -d rent_receipts -f database/schema.sql
```

### 4. Install dependencies

**Backend:**
```bash
cd backend
go mod download
```

**Frontend:**
```bash
cd frontend
npm install
```

### 5. Run in development mode

**Terminal 1 - Backend:**
```bash
make run-backend
```

**Terminal 2 - Frontend:**
```bash
make run-frontend
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## Building for Production

```bash
make build
```

This will create:
- Backend binary: `backend/bin/rent-receipt-server`
- Frontend build: `frontend/dist/`

## Deployment

See [deployment documentation](docs/deployment.md) for detailed deployment instructions.

Quick deploy:
```bash
./deployment/scripts/deploy.sh
```

## API Documentation

API documentation is available at [docs/api.md](docs/api.md)

## License

MIT

## Author

David
