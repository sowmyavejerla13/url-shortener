# URL Shortener Service

A production-ready URL Shortener REST API built with Go, Gin, PostgreSQL, JWT Authentication, and Clean Architecture.

## 🚀 Features

* User Registration
* User Login
* JWT Authentication
* Protected API Routes
* Create Short URLs
* Redirect Short URLs
* Click Tracking
* List User URLs
* Delete User URLs
* PostgreSQL Persistence
* Dockerized Database
* Database Migrations
* Clean Architecture (Handler → Service → Repository)

---

## 🛠 Tech Stack

* Go
* Gin
* PostgreSQL
* pgx/v5
* JWT
* bcrypt
* Docker
* golang-migrate

---

## 📂 Project Structure

```
cmd/
    api/

internal/
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
docs/
```

---

## 📌 API Endpoints

### Health Check

| Method | Endpoint  | Description      |
| ------ | --------- | ---------------- |
| GET    | `/health` | Check API health |

---

### Authentication

| Method | Endpoint                | Description           |
| ------ | ----------------------- | --------------------- |
| POST   | `/api/v1/auth/register` | Register a new user   |
| POST   | `/api/v1/auth/login`    | Login and receive JWT |

---

### User Profile (Protected)

| Method | Endpoint     | Description                        |
| ------ | ------------ | ---------------------------------- |
| GET    | `/api/v1/me` | Get authenticated user information |

---

### URL Management (Protected)

| Method | Endpoint           | Description                                  |
| ------ | ------------------ | -------------------------------------------- |
| GET    | `/api/v1/urls`     | Get all URLs for the authenticated user      |
| POST   | `/api/v1/shorten`  | Create a short URL                           |
| DELETE | `/api/v1/urls/:id` | Delete a URL owned by the authenticated user |

---

### Public

| Method | Endpoint      | Description                  |
| ------ | ------------- | ---------------------------- |
| GET    | `/:shortCode` | Redirect to the original URL |

---

## 🔐 Authentication

Protected endpoints require a JWT access token.

Example:

```
Authorization: Bearer <your_jwt_token>
```

Protected routes:

* `GET /api/v1/me`
* `GET /api/v1/urls`
* `POST /api/v1/shorten`
* `DELETE /api/v1/urls/:id`

---

## ⚙️ Running the Project

### Start PostgreSQL

```bash
docker compose up -d
```

### Run Migrations

```bash
migrate -path migrations -database "<your_database_url>" up
```

### Start the API

```bash
go run ./cmd/api
```

The server runs at:

```
http://localhost:8080
```

---

## 📦 Architecture

The project follows a layered architecture:

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

This separation keeps business logic independent from HTTP and database implementation.

---

## ✅ Completed

* Project setup
* Configuration management
* PostgreSQL integration
* Dockerized database
* Database migrations
* User registration
* User login
* Password hashing with bcrypt
* JWT authentication
* Authentication middleware
* URL validation
* Random short code generation
* Create short URLs
* Redirect using short code
* Click count tracking
* Get authenticated user's URLs
* Delete authenticated user's URLs
* Get authenticated user profile
* Clean Architecture implementation

---

## 🚀 Future Improvements

* React frontend
* Custom aliases
* URL expiration
* URL analytics dashboard
* Unit and integration tests
* Rate limiting
* Redis caching
* CI/CD pipeline
* Production deployment

---

## License

MIT
