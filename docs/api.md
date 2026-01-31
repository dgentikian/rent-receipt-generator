# API Documentation

Base URL: `http://localhost:8080/api/v1` (development) or `https://yourdomain.com/api/v1` (production)

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### Register

```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "address": "123 Main St",
  "city": "Paris",
  "postal_code": "75001",
  "phone": "+33 1 23 45 67 89"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  ...
}
```

#### Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "landlord": {
    "id": 1,
    "email": "user@example.com",
    ...
  }
}
```

### Landlord

#### Get Profile

```http
GET /landlord/profile
Authorization: Bearer <token>
```

#### Update Profile

```http
PUT /landlord/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "address": "123 Main St",
  "city": "Paris",
  "postal_code": "75001",
  "phone": "+33 1 23 45 67 89"
}
```

#### Upload Signature

```http
POST /landlord/signature
Authorization: Bearer <token>
Content-Type: multipart/form-data

signature: <file>
```

### Properties

#### Create Property

```http
POST /properties
Authorization: Bearer <token>
Content-Type: application/json

{
  "address": "10 Rue de la Paix",
  "city": "Paris",
  "postal_code": "75001",
  "rent_amount": 1200.00,
  "charges_amount": 100.00,
  "syndic_name": "Syndic ABC",
  "syndic_address": "5 Rue du Commerce"
}
```

#### List Properties

```http
GET /properties
Authorization: Bearer <token>
```

#### Get Property

```http
GET /properties/:id
Authorization: Bearer <token>
```

#### Update Property

```http
PUT /properties/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "address": "10 Rue de la Paix",
  "city": "Paris",
  "postal_code": "75001",
  "rent_amount": 1250.00,
  "charges_amount": 100.00
}
```

#### Delete Property

```http
DELETE /properties/:id
Authorization: Bearer <token>
```

### Tenants

#### Create Tenant

```http
POST /tenants
Authorization: Bearer <token>
Content-Type: application/json

{
  "property_id": 1,
  "first_name": "Marie",
  "last_name": "Martin",
  "email": "marie@example.com",
  "phone": "+33 6 12 34 56 78",
  "move_in_date": "2024-01-01"
}
```

#### List Tenants by Property

```http
GET /tenants?property_id=1
Authorization: Bearer <token>
```

#### Get Tenant

```http
GET /tenants/:id
Authorization: Bearer <token>
```

#### Update Tenant

```http
PUT /tenants/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Marie",
  "last_name": "Martin",
  "email": "marie.new@example.com",
  "phone": "+33 6 12 34 56 78",
  "is_active": true
}
```

#### Delete Tenant

```http
DELETE /tenants/:id
Authorization: Bearer <token>
```

### Receipts

#### Create Receipt

```http
POST /receipts
Authorization: Bearer <token>
Content-Type: application/json

{
  "property_id": 1,
  "tenant_id": 1,
  "period_month": 1,
  "period_year": 2024,
  "rent_amount": 1200.00,
  "charges_amount": 100.00,
  "payment_method": "Virement bancaire",
  "payment_date": "2024-01-05",
  "notes": "Paiement reçu"
}
```

**Response:** PDF is automatically generated

#### List Receipts

```http
GET /receipts?year=2024&month=1&property_id=1&limit=50
Authorization: Bearer <token>
```

Query parameters:
- `property_id` (optional): Filter by property
- `tenant_id` (optional): Filter by tenant
- `year` (optional): Filter by year
- `month` (optional): Filter by month
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset

#### Get Receipt

```http
GET /receipts/:id
Authorization: Bearer <token>
```

#### Get Receipt with Details

```http
GET /receipts/:id/details
Authorization: Bearer <token>
```

Returns receipt with embedded landlord, property, and tenant information.

#### Download PDF

```http
GET /receipts/:id/pdf
Authorization: Bearer <token>
```

Returns the PDF file for download.

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

- `200 OK`: Success
- `201 Created`: Resource created
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Not authorized to access resource
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Rate Limiting

Currently no rate limiting is implemented. Consider adding it in production.

## CORS

CORS is configured based on `CORS_ALLOWED_ORIGINS` environment variable.

## Security

- All passwords are hashed using bcrypt
- JWT tokens expire based on `JWT_EXPIRY` setting (default: 24h)
- HTTPS should be used in production
- Input validation is performed on all endpoints
