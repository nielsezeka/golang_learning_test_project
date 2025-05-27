# Quiz API Project

This project is a Go (Golang) API server using Gin and PostgreSQL.  
It features JWT authentication and a sample quiz API.

Built during my quest to master Go (Golang)—and to finally understand  
what all those curly braces are for. :v

If you spot any bugs, have suggestions, or strong opinions about  
tabs vs spaces, let me know!

Special thanks to GitHub Copilot, my tireless AI sidekick,  
for saving me from countless typos and existential crises.

## Prerequisites
- Go 1.18 or newer
- PostgreSQL (tested on v14+)
- [Homebrew](https://brew.sh/) (for macOS users)

## Project Features

This project demonstrates modern Go web development practices with focus on:

### Core Architecture
- **Modular Project Structure** — Clean separation of concerns with organized packages (handler, router, db, utils, docs)
- **Gin Web Framework** — RESTful API server with routing and middleware support
- **PostgreSQL Integration** — Database connection, migration support, and SQL operations

### Authentication & Security
- **JWT Authentication System** — Complete auth flow with login, registration, token refresh, and secure password handling
- **Security Middleware** — JWT token validation, CORS handling, and request authentication
- **Password Security** — Bcrypt hashing for secure password storage

### API Development
- **CRUD Operations** — Full Create, Read, Update, Delete functionality for quiz management
- **Route Handlers** — Organized endpoint handlers for authentication and quiz operations
- **Response Utilities** — Standardized JSON response formatting and error handling
- **Swagger Documentation** — Auto-generated API docs with interactive UI interface

### Development Tools
- **Comprehensive Testing** — Unit tests for router and authentication functionality
- **Database Backup/Restore** — PostgreSQL backup file for easy project setup
- **Environment Configuration** — Flexible database connection and server configuration


## 1. Install PostgreSQL

**On macOS (using Homebrew):**
```sh
brew install postgresql
brew services start postgresql
```

**On Ubuntu/Debian:**
```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start
```

## 2. Create Database and User

Open a terminal and run:
```sh
psql -U postgres -h localhost
```
Then in the psql prompt:
```sql
CREATE DATABASE test_db;
-- (Optional) Create a user if needed:
-- CREATE USER myuser WITH PASSWORD 'mypassword';
-- GRANT ALL PRIVILEGES ON DATABASE test_db TO myuser;
```

## 3. Import the Database

If you have a backup file (e.g., `postgres.backup`), restore it with:
```sh
createdb -U postgres -h localhost test_db  # Only if test_db does not exist
pg_restore -U postgres -h localhost -d test_db -c -v postgres.backup
```

## 4. Configure the Project

- The database connection string is in `db/db.go`. Adjust the user, password, and dbname if needed.

## 5. Install Go Dependencies

```sh
go mod tidy
```

## 6. Run the Project

```sh
go run main.go
```

The server will start on `http://localhost:8080`.

## 7. API Endpoints
- `POST /api/login` — Login and get JWT tokens
- `POST /api/refresh` — Refresh access token
- `GET /api/quiz` — Get a quiz (requires DB setup)
- `POST /api/quiz` — Create a quiz (sample)

## 8. Testing

Run tests with:
```sh
go test ./...
```

---

**Note:**
- For production, use secure password storage and environment variables for secrets.

**Import test database:**
```sh
pg_restore -U postgres -h localhost -d test_db -c -v postgres.backup
```
## 9. API Documentation (Swagger)

The API is documented using Swagger.  
After starting the server, you can access the Swagger UI at:  
[http://localhost:8080/swagger_ui/](http://localhost:8080/swagger_ui/) (if running locally).
