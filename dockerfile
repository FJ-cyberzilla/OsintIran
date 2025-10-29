# --- STAGE 1: Builder Stage (Go) ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ca-certificates (needed for some Go modules)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod AND go.sum (both are required)
COPY go.mod go.sum ./

# Configure Go proxy with better error handling
RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go env -w GOSUMDB=sum.golang.org

# Download dependencies with retry and better error output
RUN go mod download -x 2>&1 || \
    (echo "=== GO MOD DOWNLOAD FAILED ===" && \
     echo "Trying with GOSUMDB=off..." && \
     go env -w GOSUMDB=off && \
     go mod download -x)

# Copy source code and build statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags '-static'" -o /app/app ./main.go

# --- STAGE 2: Runtime Environment ---
FROM ubuntu:22.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1

# Set working directory
WORKDIR /app

# --- 1. Install Base Tools with better error handling ---
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    git \
    ca-certificates \
    software-properties-common \
    python3 \
    python3-pip \
    python3-venv \
    wget \
    gnupg \
    postgresql-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates

# --- 2. Install Go ---
ARG GOLANG_VERSION=1.21.5
RUN curl -fsSL -o /tmp/go.tar.gz https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# --- 3. Install Node.js + TypeScript ---
ARG NODE_VERSION=20.x
RUN mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_${NODE_VERSION} nodistro main" > /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && apt-get install -y nodejs && \
    npm install -g typescript ts-node @types/node

# --- 4. Install Julia ---
ARG JULIA_VERSION=1.10.0
RUN JULIA_MAJOR_MINOR=$(echo $JULIA_VERSION | cut -d. -f1,2) && \
    curl -fsSL -o /tmp/julia.tar.gz https://julialang-s3.julialang.org/bin/linux/x64/${JULIA_MAJOR_MINOR}/julia-${JULIA_VERSION}-linux-x86_64.tar.gz && \
    tar -C /opt -xzf /tmp/julia.tar.gz && \
    ln -s /opt/julia-${JULIA_VERSION}/bin/julia /usr/local/bin/julia && \
    rm /tmp/julia.tar.gz

# --- 5. Copy project files ---
COPY . .

# --- 6. Install Python dependencies if requirements.txt exists ---
RUN if [ -f requirements.txt ]; then \
        pip3 install --no-cache-dir --break-system-packages -r requirements.txt; \
    fi

# --- 7. Install Node.js dependencies if package.json exists ---
RUN if [ -f package.json ]; then \
        npm ci --only=production; \
    fi

# --- 8. Copy Go binary from builder stage ---
COPY --from=builder /app/app /usr/local/bin/app

# --- 9. Create non-root user for security ---
RUN groupadd -r appuser && useradd -r -g appuser appuser && \
    chown -R appuser:appuser /app
USER appuser

# --- 10. Health check ---
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# --- 11. Expose port ---
EXPOSE 8080

# --- 12. Default command ---
CMD ["/usr/local/bin/app"]
