# URL Shortener

A high-performance URL shortening service built with Go and Fiber. This service allows users to convert long URLs into short, manageable links that redirect to the original destination.

## Key Features

- URL shortening using SHA256 hashing + timestamp
- Automatic redirection to original URLs
- Flexible configuration via:
    - Environment variables
    - Command-line flags (-a for SERVER_ADDRESS; -b for BASE_URL)
    - Default values (SERVER_ADDRESS = :8080 ; BASE_URL = localhost)
- High-performance Fiber web framework
- Dependency injection with Uber FX
- Comprehensive test coverage

---

## Instalation
```bash
git clone https://github.com/VladimirAzanza/url-shortener.git
cd url-shortener
```

```bash
go build -o cmd/shortener/shortener cmd/shortener/main.go
```


## Usage

### Basic Execution
```bash
./shortener
```

### Configuration Options
Using command-line flags:

```bash
./shortener -a :8081 -b http://my-domain.com
```
Example2:
```bash
cmd/shortener/shortener -f ./files/records.json
```

Using environment variables:
```bash
SERVER_ADDRESS=:8082 BASE_URL=http://other-domain.com ./shortener
```

### Flags 

./shortener -h (to see what means each flag)

- (-a) : port
- (-b) : host
- (-f) : file storage path
- (-dt): database type (sqlite|postgres)

### 1. Shorten a URL

To shorten a URL, make a POST request to the / endpoint with the original URL in the request body.

Example with curl:
```bash
curl -X POST -d "http://myurl.com" http://localhost:8080/
```
Response example:
```bash
http://localhost:8080/12310
```

### 2. Access the Original URL
To access the original URL, make a GET request to the shortened URL.

Example with curl:
```bash
curl -v http://localhost:8080/12310
```
Response:

The server will respond with a 307 Temporary Redirect status code and a Location header pointing to the original URL.

```
< HTTP/1.1 307 Temporary Redirect
< Location: http://myurl.com
< Date: Mon, 17 Mar 2025 19:00:47 GMT
< Content-Length: 0
```

### Testing

Run all tests:
```bash
make test
```

Generate test coverage report:
```bash
make cover
```

## Dependencies

- Fiber - Fast web framework
- Uber FX - Dependency injection
- Testify - Testing toolkit

## Requirements

- Go 1.24.0
- Make (optional)
- SQLite : https://www.sqlite.org/download.html
- PostgreSQL

## Create Database with PostgreSQL

```bash
sudo apt install postgresql postgresql-contrib
sudo -u postgres psql
```

Then set the password-> postgres=# \password postgres
Quit -> \q

```bash
sudo -i -u postgres
psql -U postgres
```

Then create database -> postgres=# create database urlshortener
Check for existence of the database
```bash
sudo -u postgres psql -U postgres -c "\l"
```

Make migrations (PostgreSQL):
```bash
sudo -u postgres psql -d urlshortener

CREATE TABLE short_urls (
    uuid UUID PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE
);

```

Make migrations (SQLite):
```bash
sqlite3 urlshortener.db

CREATE TABLE short_urls (
    uuid UUID PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE
);
```

