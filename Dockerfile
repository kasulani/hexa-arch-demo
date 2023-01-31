FROM golang:1.17 AS dev-image

# Define current working directory
WORKDIR /shortener

RUN apt-get update
RUN apt-get install -y ca-certificates

# Download modules to local cache so we can skip re-downloading
# on consecutive docker-compose builds commands
COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add sources
COPY . .

# Build dev dependencies
RUN make deps-dev

# Build migration tool
RUN make build-migration-tool
