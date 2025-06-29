# --- Stage 1: Go Build ---
FROM golang:1.22 AS builder

# Set working directory to match the Go module root
WORKDIR /app/cherry-blossom-hunters-app

# Copy go.mod and go.sum first for better layer caching
COPY cherry-blossom-hunters-app/go.mod cherry-blossom-hunters-app/go.sum ./
RUN go mod download

# Copy the entire Go module source code
COPY cherry-blossom-hunters-app/ ./

# Build the Go binary from the module root
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# --- Stage 2: Runtime ---
FROM python:3.11-slim

# Install essential system packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        jq \
        tree \
        curl \
        bash \
        ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Create a Python virtual environment and install required packages
ENV VENV_PATH=/opt/venv

RUN python3 -m venv $VENV_PATH && \
    $VENV_PATH/bin/pip install --no-cache-dir beautifulsoup4

# Ensure the virtual environment's Python is used by default
ENV PATH="$VENV_PATH/bin:$PATH"

# Set the working directory for the final container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/server .

# Copy the Python script directory
COPY --from=builder /app/cherry-blossom-hunters-app/script ./script

# Copy any additional files if needed
COPY --from=builder /app/cherry-blossom-hunters-app/swagger.yml ./

# Expose the application port (adjust if needed)
EXPOSE 8080

# Default command: run the Go server binary,
# which may call Python scripts using the virtual environment
CMD ["./server"]
