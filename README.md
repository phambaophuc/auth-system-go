# Authentication & Authorization System

The authentication and authorization system is built in Go, following SOLID principles and Clean Architecture.

## 🏗️ Architecture

The system is designed with Clean Architecture and consists of 4 main layers:

- **Domain Layer**: Entities, Repositories, Services interfaces
- **Application Layer**: Business logic, DTOs, Use cases
- **Infrastructure Layer**: Database, External services, Security
- **Interface Layer**: HTTP handlers, Middleware, Routes

## 🚀 Features

### Authentication

- ✅ User registration
- ✅ Login with email/password
- ✅ JWT Access & Refresh tokens
- ✅ Logout
- ✅ Refresh token

### Authorization

- ✅ Role-based access control (RBAC)
- ✅ Permission-based authorization
- ✅ Authentication middleware
- ✅ Permission checking middleware

### User Management

- ✅ View user profile
- ✅ Update personal information
- ✅ Change password
- ✅ Assign/Remove roles to users (Admin)

## 🛠️ Technologies Used

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt
- **Config Management**: godotenv

## 📁 Project Structure

```
auth-system/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities & interfaces
│   ├── application/     # Business logic & DTOs
│   ├── infrastructure/  # External dependencies
│   └── interfaces/      # HTTP handlers & routes
├── pkg/                 # Reusable packages
└── .env.example        # Environment template
```

## 🚀 Installation and Running

### 1. Clone repository

```bash
git clone <repository-url>
cd auth-system
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure the database

```bash
# Create database PostgreSQL
createdb auth_db

# Copy and edit environment file
cp .env.example .env
# Update database info in .env file
```

### 4. Run the application

```bash
go run cmd/server/main.go
```

The server will run on port 8080 (or the port configured in .env).

## 📚 API Documentation

### Authentication Endpoints

#### POST /api/v1/auth/register

Register a new user

```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### POST /api/v1/auth/login

Login

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### POST /api/v1/auth/refresh

Refresh token

```json
{
  "refresh_token": "your_refresh_token_here"
}
```

#### POST /api/v1/auth/logout

Logout (requires authentication)

### User Endpoints

#### GET /api/v1/user/profile

Get user profile (requires authentication)

#### PUT /api/v1/user/profile

Update profile (requires authentication)

```json
{
  "first_name": "John",
  "last_name": "Doe"
}
```

#### POST /api/v1/user/change-password

Change password (requires authentication)

```json
{
  "current_password": "old_password",
  "new_password": "new_password123"
}
```

### Admin Endpoints

#### POST /api/v1/admin/users/:id/roles

Assign a role to a user (requires "users.write" permission)

```json
{
  "role_id": 1
}
```

#### DELETE /api/v1/admin/users/:id/roles/:roleId

Remove a role from a user (requires "users.write" permission)

## 🔐 Authentication & Authorization

### JWT Tokens

- **Access Token**: Short-lived (15 minutes), used to authenticate API calls
- **Refresh Token**: Long-lived (7 days), used to get new access tokens

### Headers

To access protected endpoints, add this header:

```
Authorization: Bearer <your_access_token>
```

### Default Roles & Permissions

The system automatically creates:

- **Role "admin"**: Full permissions
- **Role "user"**: Read-only access to user info
- **Permissions**: users.read, users.write, users.delete, roles.read, roles.write, roles.delete

## 🧪 Testing with curl

#### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### Get Profile (with token)

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <your_access_token>"
```

## 🔧 Applied SOLID Principles

1. **Single Responsibility**: Each struct/package has a single responsibility
2. **Open/Closed**: Use interfaces to extend functionality
3. **Liskov Substitution**: Implementations can substitute their interfaces
4. **Interface Segregation**: Small and focused interfaces
5. **Dependency Inversion**: High-level modules do not depend on low-level modules

## 📝 License

MIT License
