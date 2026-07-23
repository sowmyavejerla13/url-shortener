# URL Shortener Service - Implementation Guide

> This document explains the complete implementation process of the URL Shortener Service. It includes the project setup, architectural decisions, implementation steps, common issues encountered during development, and important backend concepts learned throughout the project.

---

# Project Overview

The URL Shortener Service is a production-ready REST API built using Go and Gin. It allows authenticated users to create and manage shortened URLs while following Clean Architecture principles.

## Features

* User Registration
* User Login
* JWT Authentication
* Protected APIs
* Create Short URLs
* Redirect Short URLs
* Click Tracking
* List User URLs
* Delete User URLs
* PostgreSQL Database
* Dockerized Development Environment
* Database Migrations
* Swagger Documentation
* Unit Testing

---

# Project Initialization

Initialize the project.

```bash
go mod init github.com/sowmyavejerla13/url-shortener
```

The module name should match the GitHub repository.

---

# Initial Git Commit Plan

The project was implemented incrementally using meaningful commits.

```
chore: initialize project structure

chore: add application configuration

feat: create HTTP server

feat: add health check endpoint

feat: implement graceful shutdown

feat: add PostgreSQL connection

feat: implement user registration

feat: implement JWT authentication

feat: create URL shortening endpoint

feat: add redirect endpoint

feat: add click analytics

test: add authentication unit tests

docs: update API documentation
```

Using small commits makes the project history easier to understand and maintain.

---

# Step 1 — Install Dependencies

Install the required packages.

```bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get github.com/spf13/viper
go get github.com/google/uuid
```

After installing dependencies:

```bash
go mod tidy
```

This downloads required packages and generates the `go.sum` file.

---

# Step 2 — Application Configuration

Configuration is handled using:

* Viper
* godotenv

Environment variables are loaded from the project root using the `.env` file.

---

## Important Note

Initially an error occurred because there was no `main` package.

The application's entry point is:

```
cmd/api/main.go
```

Run the application from the project root:

```bash
go run ./cmd/api
```

Do **not** execute the command from inside the `cmd/api` directory.

---

## Why?

`godotenv.Load()` searches for `.env` in the **current working directory**.

Correct:

```
url-shortener/
    .env
    cmd/
```

Running from the project root loads the environment correctly.

If executed inside `cmd/api`, Go searches for:

```
cmd/api/.env
```

which does not exist.

---

# HTTP Server

The server is responsible for:

* Loading configuration
* Connecting to PostgreSQL
* Registering routes
* Starting Gin
* Supporting graceful shutdown

---

# Graceful Shutdown

Graceful shutdown ensures that active requests finish before the server exits.

Benefits:

* Prevents interrupted requests
* Properly closes database connections
* Safe for production deployments

---

# PostgreSQL Setup

Instead of installing PostgreSQL locally, Docker is used.

Advantages:

* Consistent environment
* Easy setup
* No local installation
* Easy cleanup
* Portable across machines

---

## Docker Compose

Database container:

```
url-shortener-db
```

---

## Why use a named volume?

```
postgres_data:
```

Without a volume:

```
docker compose down
```

would remove all stored data.

Using a named volume keeps the database persistent.

---

## Port Mapping

```
5432:5432
```

Meaning:

```
Local Machine
localhost:5432

↓

Docker Container

↓

PostgreSQL
```

The Go application connects using:

```
localhost:5432
```

---

## Start PostgreSQL

```bash
docker compose up -d
```

Verify running containers:

```bash
docker ps
```

---

# PostgreSQL Driver

Install pgx.

```bash
go get github.com/jackc/pgx/v5
```

The project uses **pgx** instead of the standard `database/sql` driver.

---

# Why use a Connection Pool?

The application returns:

```go
*pgxpool.Pool
```

instead of:

```go
pgx.Conn
```

or

```go
sql.DB
```

Reason:

A backend server handles many concurrent requests.

If 100 users access the application simultaneously, creating a new database connection for every request would be inefficient.

A connection pool:

* Reuses connections
* Improves performance
* Reduces database overhead
* Is the standard approach in production applications

---

# Summary

At this stage, the project has:

* Go module initialized
* Dependencies installed
* Configuration management added
* HTTP server created
* Graceful shutdown implemented
* PostgreSQL running in Docker
* Database connection pool configured
* Initial Git commit history established

The next phase covers **database migrations and user authentication**.


# URL Shortener Service - Implementation Guide (Part 2)

# Database Migrations

Instead of creating database tables manually using pgAdmin or SQL clients, the project uses **database migrations**.

This approach ensures that every developer working on the project has the exact same database schema.

---

## Install golang-migrate

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify installation:

```bash
migrate -version
```

---

## Create a Migration

Generate a migration file:

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

This creates two files:

```text
000001_create_users_table.up.sql

000001_create_users_table.down.sql
```

The **up** migration creates the table.

The **down** migration rolls it back.

---

## Users Table

The users table stores:

* UUID as Primary Key
* Name
* Email
* Password Hash
* Created Timestamp

Passwords are **never stored in plain text**.

Only the bcrypt hash is stored.

---

## Apply Migration

Run:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable" up
```

---

## Verify Database

Open PostgreSQL inside Docker:

```bash
docker exec -it url-shortener-db psql -U postgres -d url_shortener
```

List tables:

```sql
\dt
```

You should see:

```text
users
```

---

# Why Database Migrations?

Many beginners create tables manually.

Professional teams don't.

Typical workflow:

Developer A

↓

Creates migration

↓

Pushes to GitHub

↓

Developer B

↓

Pulls repository

↓

Runs

```bash
migrate up
```

↓

Everyone gets the exact same database schema.

This keeps environments consistent and avoids "works on my machine" problems.

---

# SQL Queries

## Why QueryRow()?

Since email is unique:

```sql
SELECT * FROM users WHERE email=$1;
```

can return at most one row.

Therefore:

```go
QueryRow()
```

is the correct choice.

Instead of:

```go
Query()
```

---

# SQL Placeholders

Instead of:

```go
fmt.Sprintf(...)
```

or string concatenation,

always use placeholders:

```sql
WHERE email = $1
```

Benefits:

* Prevents SQL Injection
* Allows PostgreSQL to cache execution plans
* Standard production practice

---

# Password Hashing

Install bcrypt.

```bash
go get golang.org/x/crypto/bcrypt
```

Passwords should **never** be stored directly.

Registration flow:

```text
Password

↓

bcrypt.GenerateFromPassword()

↓

Hash

↓

Store Hash in Database
```

Login flow:

```text
Password

↓

bcrypt.CompareHashAndPassword()

↓

Success / Failure
```

Even if the database is compromised, attackers cannot immediately recover user passwords.

---

# User Registration Flow

The registration process follows Clean Architecture.

```text
HTTP Request

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL
```

---

## Handler Responsibility

The handler only performs three tasks.

### 1. Read JSON

```go
c.ShouldBindJSON(&req)
```

Converts:

```json
{
  "name":"Sowmya",
  "email":"sowmya@example.com",
  "password":"password123"
}
```

into:

```go
req.Name

req.Email

req.Password
```

---

### 2. Call Service

```go
h.userService.Register(...)
```

Notice:

No SQL.

No bcrypt.

No JWT.

Those responsibilities belong to other layers.

---

### 3. Return Response

Success:

```json
{
    "message":"User registered successfully"
}
```

Failure:

```json
{
    "error":"email already exists"
}
```

---

# Validation

Install validator.

```bash
go get github.com/go-playground/validator/v10
```

Validation ensures:

* Required fields
* Email format
* Password length

Invalid requests never reach the service layer.

---

# JWT Authentication

Install JWT package.

```bash
go get github.com/golang-jwt/jwt/v5
```

---

## JWT Payload

Generated token contains:

```json
{
    "user_id":"...",
    "exp":1752900000,
    "iat":1752813600
}
```

Claims:

### user_id

Identifies the authenticated user.

### exp

Token expiration.

Current implementation:

24 hours.

### iat

Token issued timestamp.

---

## Signing JWT

The token is signed using:

```text
JWT_SECRET
```

stored in the environment file.

---

# Login Flow

```text
Email

↓

Find User

↓

Compare Password

↓

Generate JWT

↓

Return Token
```

---

# JWT Middleware

Every protected request sends:

```text
Authorization: Bearer <JWT_TOKEN>
```

The middleware performs:

```text
Read Header

↓

Extract Token

↓

Verify Signature

↓

Validate Expiration

↓

Extract user_id

↓

Store user_id in Gin Context

↓

Call Next Handler
```

After middleware executes, handlers can simply access:

```go
c.GetString("userID")
```

without parsing the JWT again.

---

# Why Keep Handlers Thin?

A handler should only:

* Read Request
* Call Service
* Return Response

Business logic belongs in the Service layer.

Database operations belong in the Repository layer.

This separation makes the code:

* Easier to test
* Easier to maintain
* Easier to extend

---

# Summary

At this stage the application supports:

* Database migrations
* User table creation
* Password hashing
* User registration
* Login
* JWT generation
* JWT middleware
* Request validation
* Clean Architecture separation

The next phase implements the URL Shortening functionality.



# URL Shortener Service - Implementation Guide (Part 2)

# Database Migrations

Instead of creating database tables manually using pgAdmin or SQL clients, the project uses **database migrations**.

This approach ensures that every developer working on the project has the exact same database schema.

---

## Install golang-migrate

```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Verify installation:

```bash
migrate -version
```

---

## Create a Migration

Generate a migration file:

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

This creates two files:

```text
000001_create_users_table.up.sql

000001_create_users_table.down.sql
```

The **up** migration creates the table.

The **down** migration rolls it back.

---

## Users Table

The users table stores:

* UUID as Primary Key
* Name
* Email
* Password Hash
* Created Timestamp

Passwords are **never stored in plain text**.

Only the bcrypt hash is stored.

---

## Apply Migration

Run:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable" up
```

---

## Verify Database

Open PostgreSQL inside Docker:

```bash
docker exec -it url-shortener-db psql -U postgres -d url_shortener
```

List tables:

```sql
\dt
```

You should see:

```text
users
```

---

# Why Database Migrations?

Many beginners create tables manually.

Professional teams don't.

Typical workflow:

Developer A

↓

Creates migration

↓

Pushes to GitHub

↓

Developer B

↓

Pulls repository

↓

Runs

```bash
migrate up
```

↓

Everyone gets the exact same database schema.

This keeps environments consistent and avoids "works on my machine" problems.

---

# SQL Queries

## Why QueryRow()?

Since email is unique:

```sql
SELECT * FROM users WHERE email=$1;
```

can return at most one row.

Therefore:

```go
QueryRow()
```

is the correct choice.

Instead of:

```go
Query()
```

---

# SQL Placeholders

Instead of:

```go
fmt.Sprintf(...)
```

or string concatenation,

always use placeholders:

```sql
WHERE email = $1
```

Benefits:

* Prevents SQL Injection
* Allows PostgreSQL to cache execution plans
* Standard production practice

---

# Password Hashing

Install bcrypt.

```bash
go get golang.org/x/crypto/bcrypt
```

Passwords should **never** be stored directly.

Registration flow:

```text
Password

↓

bcrypt.GenerateFromPassword()

↓

Hash

↓

Store Hash in Database
```

Login flow:

```text
Password

↓

bcrypt.CompareHashAndPassword()

↓

Success / Failure
```

Even if the database is compromised, attackers cannot immediately recover user passwords.

---

# User Registration Flow

The registration process follows Clean Architecture.

```text
HTTP Request

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL
```

---

## Handler Responsibility

The handler only performs three tasks.

### 1. Read JSON

```go
c.ShouldBindJSON(&req)
```

Converts:

```json
{
  "name":"Sowmya",
  "email":"sowmya@example.com",
  "password":"password123"
}
```

into:

```go
req.Name

req.Email

req.Password
```

---

### 2. Call Service

```go
h.userService.Register(...)
```

Notice:

No SQL.

No bcrypt.

No JWT.

Those responsibilities belong to other layers.

---

### 3. Return Response

Success:

```json
{
    "message":"User registered successfully"
}
```

Failure:

```json
{
    "error":"email already exists"
}
```

---

# Validation

Install validator.

```bash
go get github.com/go-playground/validator/v10
```

Validation ensures:

* Required fields
* Email format
* Password length

Invalid requests never reach the service layer.

---

# JWT Authentication

Install JWT package.

```bash
go get github.com/golang-jwt/jwt/v5
```

---

## JWT Payload

Generated token contains:

```json
{
    "user_id":"...",
    "exp":1752900000,
    "iat":1752813600
}
```

Claims:

### user_id

Identifies the authenticated user.

### exp

Token expiration.

Current implementation:

24 hours.

### iat

Token issued timestamp.

---

## Signing JWT

The token is signed using:

```text
JWT_SECRET
```

stored in the environment file.

---

# Login Flow

```text
Email

↓

Find User

↓

Compare Password

↓

Generate JWT

↓

Return Token
```

---

# JWT Middleware

Every protected request sends:

```text
Authorization: Bearer <JWT_TOKEN>
```

The middleware performs:

```text
Read Header

↓

Extract Token

↓

Verify Signature

↓

Validate Expiration

↓

Extract user_id

↓

Store user_id in Gin Context

↓

Call Next Handler
```

After middleware executes, handlers can simply access:

```go
c.GetString("userID")
```

without parsing the JWT again.

---

# Why Keep Handlers Thin?

A handler should only:

* Read Request
* Call Service
* Return Response

Business logic belongs in the Service layer.

Database operations belong in the Repository layer.

This separation makes the code:

* Easier to test
* Easier to maintain
* Easier to extend

---

# Summary

At this stage the application supports:

* Database migrations
* User table creation
* Password hashing
* User registration
* Login
* JWT generation
* JWT middleware
* Request validation
* Clean Architecture separation

The next phase implements the URL Shortening functionality.



# URL Shortener Service - Implementation Guide (Part 3)

# URL Shortening Module

After implementing authentication, the next step is allowing authenticated users to create and manage shortened URLs.

Each shortened URL belongs to exactly one authenticated user.

---

# Database Design

A new table named **urls** is created.

Important fields:

* ID (UUID)
* User ID
* Original URL
* Short Code
* Click Count
* Created At

---

# Foreign Key Relationship

The `user_id` column references the `users` table.

```sql id="hyjlwm"
CONSTRAINT fk_urls_users
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE
```

---

## Why ON DELETE CASCADE?

Imagine a user owns several URLs.

```text id="4szh5l"
John

↓

google

youtube

github

stackoverflow
```

If John deletes his account, should these URLs remain?

Usually **no**.

Using:

```sql id="31ej6m"
ON DELETE CASCADE
```

PostgreSQL automatically removes every URL belonging to that user.

No additional application logic is required.

---

# URL Creation Flow

```text id="l0vm3l"
HTTP Request

↓

JWT Middleware

↓

Handler

↓

Service

↓

Repository

↓

Database
```

---

# URL Validation

Before creating a shortened URL, the request is validated.

Checks include:

* Required field
* Valid URL format

Invalid URLs never reach the repository.

---

# Random Short Code Generation

Every shortened URL receives a randomly generated code.

Character set:

```text id="t6eqn6"
abcdefghijklmnopqrstuvwxyz

ABCDEFGHIJKLMNOPQRSTUVWXYZ

0123456789
```

Total characters:

```text id="v9fv31"
26 lowercase

26 uppercase

10 digits

=

62 characters
```

---

## Total Possible Combinations

Current implementation:

```text id="v4zx7q"
Length = 6
```

Possible combinations:

```text id="3lv74u"
62^6

=

56,800,235,584
```

More than **56 billion** possible short codes.

Collisions are extremely unlikely.

---

# Why crypto/rand?

Instead of `math/rand`, the project uses:

```go id="4lw4w8"
crypto/rand
```

Reason:

It generates cryptographically secure random numbers.

Better suited for identifiers exposed publicly.

---

# Why big.NewInt?

`crypto/rand.Int()` expects:

```go id="ppam1x"
*big.Int
```

Therefore:

```go id="4v22ek"
big.NewInt(int64(len(charset)))
```

creates a random range:

```text id="d7dsi8"
0

↓

61
```

One random character is selected each iteration until the short code reaches six characters.

---

# Collision Handling

Even though collisions are highly unlikely, production systems should still verify uniqueness.

Typical flow:

```text id="bnx1yn"
Generate Code

↓

Check Database

↓

Exists?

↓

Generate Again
```

The probability is extremely low but should never be ignored.

---

# Repository Responsibilities

Repository methods communicate directly with PostgreSQL.

Examples include:

* Create URL
* Find by Short Code
* Get User URLs
* Delete URL
* Increment Click Count

Repositories should contain SQL only.

Business logic belongs in the Service layer.

---

# PostgreSQL Operations

Different database operations require different methods.

---

## QueryRow()

Used when expecting one row.

Examples:

```text id="1q84ht"
Find User

Find URL by Short Code
```

---

## Query()

Used when expecting multiple rows.

Example:

```text id="tgcj2j"
List all URLs for a user
```

---

## Exec()

Used for operations that don't return rows.

Examples:

```text id="vjlwmv"
UPDATE

DELETE
```

Good mental model:

```text id="2vjlwm"
INSERT + RETURNING

↓

QueryRow()

------------------

SELECT one

↓

QueryRow()

------------------

SELECT many

↓

Query()

------------------

UPDATE

DELETE

↓

Exec()
```

---

# Redirect Flow

Redirecting is straightforward.

```text id="0g98h4"
Client

↓

GET /abc123

↓

Repository

↓

Original URL

↓

HTTP Redirect
```

The browser automatically opens the destination.

---

# Click Analytics

Every redirect increases the click count.

SQL:

```sql id="h9oziv"
UPDATE urls
SET click_count = click_count + 1
WHERE id = $1;
```

---

## Why Update Directly?

Avoid this pattern:

```text id="hm4eoq"
SELECT click_count

↓

click_count++

↓

UPDATE
```

This introduces race conditions.

Instead PostgreSQL performs:

```sql id="tlw8ur"
click_count = click_count + 1
```

atomically.

Advantages:

* Thread-safe
* Faster
* Production standard

---

# Listing User URLs

Authenticated users can retrieve all URLs they own.

Flow:

```text id="eh9xg6"
JWT

↓

Middleware

↓

userID

↓

Repository

↓

[]URL

↓

JSON Response
```

Only URLs belonging to the authenticated user are returned.

---

# Delete URL

Deletion follows authorization checks.

```text id="zcfzhh"
Request

↓

Authenticated User

↓

Verify Ownership

↓

Delete

↓

204 No Content
```

If the URL belongs to another user:

```text id="xljlwm"
403 Forbidden
```

If it doesn't exist:

```text id="ppjlwm"
404 Not Found
```

---

# Error Handling

Custom application errors are used.

Examples:

```text id="84o5k2"
ErrURLNotFound

ErrForbidden
```

Handlers translate these into proper HTTP status codes.

---

# Clean Architecture Review

At this stage the application follows:

```text id="99h4s9"
HTTP Request

↓

Router

↓

Middleware

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL
```

Each layer has exactly one responsibility.

---

# Summary

After completing this phase, the application supports:

* URL creation
* Random short code generation
* Secure ownership using JWT
* Redirect endpoint
* Click tracking
* List user URLs
* Delete URLs
* Foreign key relationships
* Atomic SQL updates
* Proper repository design
* Clean Architecture


# Part 4 — Testing, Swagger & Final Notes

---

# Step 13 — Unit Testing

The project includes unit tests for:

- Service Layer
- Handler Layer

Repository layer is intentionally not unit tested because it directly interacts with PostgreSQL. Repository behavior is validated through integration testing.

---

## Install Testify

```bash
go get github.com/stretchr/testify
go mod tidy
```

Used packages:

```go
github.com/stretchr/testify/assert
```

---

## Mock Strategy

Instead of calling the real database, lightweight mocks are used.

Example:

```go
type UserRepositoryMock struct {
    GetByEmailFunc func(email string) (*model.User, error)
    CreateFunc     func(user *model.User) error
}
```

The mock simply calls the configured function.

This makes tests:

- Fast
- Independent
- Deterministic

---

## Dependency Injection for Testing

The service layer uses dependency injection to allow mocking.

Example:

```go
type UserService struct {
    repo repository.UserRepositoryInterface

    hashPassword    func([]byte, int) ([]byte, error)
    comparePassword func([]byte, []byte) error
    generateToken   func(string, string) (string, error)
}
```

During production:

```go
service.NewUserService(repo, cfg)
```

uses:

- bcrypt.GenerateFromPassword
- bcrypt.CompareHashAndPassword
- utils.GenerateToken

During tests these functions are replaced with fake implementations.

This makes password hashing and JWT generation fully testable.

---

## Table Driven Tests

All handlers and services use table-driven tests.

Example:

```go
tests := []struct {
    name string

    input string

    expectedStatus int

    expectedBody string
}{
    ...
}
```

Benefits:

- Easy to extend
- Cleaner code
- Standard Go testing style

---

## Running Tests

Run all tests:

```bash
go test ./...
```

Verbose mode:

```bash
go test -v ./...
```

Coverage:

```bash
go test -cover ./...
```

Coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

# Step 14 — Swagger Documentation

Swagger is used for API documentation.

---

## Install

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

```bash
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/swag
```

---

## API Metadata

After:

```go
package main
```

Add:

```go
// @title URL Shortener API
// @version 1.0
// @description Production-ready URL Shortener API built with Go, Gin and PostgreSQL.
// @termsOfService http://swagger.io/terms/

// @contact.name Sowmya Vejerla
// @contact.url https://github.com/sowmyavejerla13
// @contact.email your-email@example.com

// @license.name MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

---

## Generate Swagger Docs

From project root:

```bash
swag init -g cmd/api/main.go
```

Generated files:

```
docs/
    docs.go
    swagger.json
    swagger.yaml
```

---

## Register Swagger Route

Import:

```go
_ "github.com/sowmyavejerla13/url-shortener/docs"

swaggerFiles "github.com/swaggo/files"
ginSwagger "github.com/swaggo/gin-swagger"
```

Add route:

```go
router.GET(
    "/swagger/*any",
    ginSwagger.WrapHandler(swaggerFiles.Handler),
)
```

Open:

```
http://localhost:8080/swagger/index.html
```

---

## Document Endpoints

Example:

```go
// Register godoc
//
// @Summary Register a new user
// @Description Creates a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]string
// @Router /register [post]
```

After modifying annotations:

```bash
swag init -g cmd/api/main.go
```

Always regenerate documentation before running the server.

---

# Final Project Structure

```
cmd/
    api/

docs/

internal/
    apperrors/
    config/
    database/
    dto/
    handler/
    middleware/
    model/
    repository/
    routes/
    service/
    utils/

migrations/

tests/
```

---

# Clean Architecture Flow

```
HTTP Request
      │
      ▼
Gin Router
      │
      ▼
Middleware
      │
      ▼
Handler
      │
      ▼
Service
      │
      ▼
Repository
      │
      ▼
PostgreSQL
```

Responsibilities:

- Handler → HTTP only
- Service → Business logic
- Repository → Database operations
- Model → Database entities
- DTO → Request/Response payloads
- Middleware → Authentication
- Utils → Helper functions

---

# Production Practices Followed

- Clean Architecture
- Dependency Injection
- JWT Authentication
- Password Hashing with bcrypt
- PostgreSQL Connection Pool
- Database Migrations
- Dockerized Database
- UUID Primary Keys
- Swagger Documentation
- Unit Testing
- Table-Driven Tests
- Repository Pattern
- Service Layer Pattern
- Validation using go-playground/validator
- Graceful Shutdown
- Environment-based Configuration

---

# Lessons Learned

Throughout this project, the focus was not just on building a URL Shortener API but on understanding how production-grade backend applications are structured.

Key takeaways include:

- Separating responsibilities using Clean Architecture.
- Writing maintainable and testable code through dependency injection.
- Using Docker to ensure consistent local development.
- Managing schema changes with database migrations.
- Securing APIs using JWT authentication and bcrypt password hashing.
- Documenting APIs using Swagger.
- Building confidence in Go testing using mocks and table-driven tests.

This project serves as a strong foundation for building larger distributed backend systems.


---

# Part 5 — Testing, API Documentation & Deployment

## Unit Testing

The project includes unit tests for:

- Service layer
- Handler layer

Repository tests were intentionally omitted because they primarily interact with PostgreSQL. Those are better suited as integration tests.

### Testing Packages

```bash
go get github.com/stretchr/testify
go mod tidy
```

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Generate coverage:

```bash
go test ./... -cover
```

Example:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Mocking

The project uses manually created mocks instead of third-party mocking libraries.

Example:

```
handler
        ↓
mock service

service
        ↓
mock repository
```

This keeps tests:

- fast
- deterministic
- independent of external systems

---

## Swagger API Documentation

Install Swagger CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Install dependencies

```bash
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/swag
```

---

### Swagger Metadata

Add API metadata above `func main()`.

Example:

```go
// @title URL Shortener API
// @version 1.0
// @description Production-ready URL Shortener API built with Go, Gin and PostgreSQL.

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

---

### Generate Documentation

```bash
swag init -g cmd/api/main.go
```

Generated directory:

```
docs/
    docs.go
    swagger.json
    swagger.yaml
```

---

### Import Swagger

```go
import (
    _ "github.com/sowmyavejerla13/url-shortener/docs"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)
```

---

### Register Swagger Route

```go
router.GET(
    "/swagger/*any",
    ginSwagger.WrapHandler(swaggerFiles.Handler),
)
```

Open

```
http://localhost:8080/swagger/index.html
```

---

### Important Note

Whenever Swagger annotations change, regenerate documentation.

```bash
swag init -g cmd/api/main.go
```

Restart the application after regeneration.

---

# Deployment

### Build the application

```bash
go build -o url-shortener ./cmd/api
```

---

### Run

```bash
./url-shortener
```

or

```bash
go run ./cmd/api
```

---

### Start PostgreSQL

```bash
docker compose up -d
```

Verify

```bash
docker ps
```

---

### Apply migrations

```bash
migrate \
-path migrations \
-database "<DATABASE_URL>" \
up
```

---

## Environment Variables

Required:

```
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=url_shortener

JWT_SECRET=your-secret-key
BASE_URL=http://localhost:8080
```

---

# Project Architecture

```
Request
    │
    ▼
Gin Router
    │
Middleware
    │
Handler
    │
Service
    │
Repository
    │
PostgreSQL
```

Each layer has a single responsibility.

---

# Engineering Practices Followed

- Clean Architecture
- Dependency Injection
- Repository Pattern
- Service Layer
- DTO Separation
- Configuration Management
- JWT Authentication
- Password Hashing
- Dockerized Database
- Database Migrations
- Swagger Documentation
- Unit Testing
- Graceful Shutdown
- Environment Configuration

---

# Common Issues Encountered

## "at least one main package"

Run from project root:

```bash
go run ./cmd/api
```

---

## `.env` not found

Always execute from the project root because `godotenv.Load()` searches the current working directory.

---

## Swagger not updating

Regenerate:

```bash
swag init -g cmd/api/main.go
```

Restart the application.

---

## Migration already applied

Check migration status before rerunning.

---

## Docker database not accessible

Verify container:

```bash
docker ps
```

Check logs:

```bash
docker logs url-shortener-db
```

---

# Lessons Learned

During this project I learned:

- Designing layered backend architecture
- Writing production-ready REST APIs
- Using PostgreSQL with pgx
- Managing schema changes with migrations
- JWT authentication
- Password hashing with bcrypt
- Docker fundamentals
- API documentation using Swagger
- Writing unit tests with mocks
- Dependency injection for testability
- Graceful shutdown patterns
- Environment-based configuration
- Production-oriented Go project structure

This project serves as a strong foundation for building larger backend systems using Go.

---