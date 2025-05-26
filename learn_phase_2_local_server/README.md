# Quiz API Project

This project is a Go (Golang) API server using Gin and PostgreSQL, featuring JWT authentication and a sample quiz API.

## Prerequisites
- Go 1.18 or newer
- PostgreSQL (tested on v14+)
- [Homebrew](https://brew.sh/) (for macOS users)

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
- **Export your database:**
```sh
pg_dump -U postgres -h localhost -F c -b -v -f postgres.backup test_db
```

**Import your database:**
```sh
pg_restore -U postgres -h localhost -d test_db -c -v postgres.backup
```
