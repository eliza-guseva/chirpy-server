# Chirpy - A Twitter-like Social Media Server

A learning project for building HTTP servers in Go. Chirpy is a social media backend that lets users post short messages ("chirps"), manage accounts, and interact with a Twitter-like platform.

## What This Project Does

- **User Management**: Create accounts, login/logout, password hashing with bcrypt
- **Authentication**: JWT tokens with refresh token support
- **Post Chirps**: Create and share short messages (140 characters max)
- **Social Features**: View all chirps, filter by author, sort by date
- **Content Moderation**: Automatic profanity filtering
- **Premium Features**: Upgrade users to "Chirpy Red" via webhook integration
- **Database**: PostgreSQL with SQLC for type-safe SQL queries

## Why Someone Should Care

This project demonstrates modern Go web development patterns:

- **Clean Architecture**: Layered structure with handlers, database layer, and internal packages
- **Type-Safe Database Access**: Uses SQLC to generate Go code from SQL queries
- **Standard Library Focus**: Built with `net/http` instead of frameworks like Gin or Echo
- **Real Authentication**: Proper JWT implementation with secure password hashing
- **Production Patterns**: Environment configuration, structured logging, middleware

Perfect for developers learning Go web development or studying backend architecture patterns.

## How to Install and Run

### Prerequisites

- Go 1.22+
- PostgreSQL
- SQLC (for code generation)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/eliza-guseva/chirpy-server.git
   cd chirpy-server
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables:**
   Create a `.env` file:
   ```env
   DB_URL=postgres://username:password@localhost/chirpy_db?sslmode=disable
   JWT_SECRET=your-secret-key-here
   PLATFORM=dev
   ```

4. **Set up the database:**
   - Create a PostgreSQL database called `chirpy_db`
   - Run your database migrations (schema files in `sql/schema/`)

5. **Generate database code:**
   ```bash
   sqlc generate
   ```

6. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

### API Endpoints

- `POST /api/users` - Create user account
- `POST /api/login` - User login
- `GET /api/chirps` - Get all chirps (supports `?author_id=` and `?sort=` query params)
- `POST /api/chirps` - Create new chirp (requires authentication)
- `DELETE /api/chirps/{id}` - Delete chirp (requires authentication)

### Development Commands

```bash
go test ./...          # Run tests
go vet ./...          # Static analysis
go fmt ./...          # Format code
sqlc generate         # Regenerate database code after SQL changes
```

---

*This is a learning project created as part of studying Go web development and HTTP server patterns.*