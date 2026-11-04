set dotenv-load

image := env("IMAGE", "depot")
version := `git describe --tags --always --dirty`
commit := `git rev-parse HEAD`
short := `git rev-parse --short HEAD`

export DATABASE_URL := "postgres://depot:" + env("POSTGRES_PASSWORD", "") + "@localhost:5432/depot?sslmode=disable"

# List all recipes
default:
    @just --list

# Build and start all services
up:
    docker compose up -d --build

# Stop all services, data survives
down:
    docker compose down

# Start only the database
db:
    docker compose up -d --wait db

# Check the code
lint:
    go vet ./...

# Run the tests, needs a running database
test:
    go test ./...

# Build the image, stamped with version, commit and build time
build:
    docker build \
        --build-arg VERSION={{version}} \
        --build-arg COMMIT={{commit}} \
        --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
        -t {{image}}:{{short}} -t {{image}}:latest .

# Push the image built by build
push:
    docker push {{image}}:{{short}}
    docker push {{image}}:latest
