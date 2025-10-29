# --- STAGE 1: Builder Stage (Go) ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ca-certificates (needed for some Go modules)
RUN apk add --no-cache git ca-certificates

# Copy go.mod AND go.sum (both are required)
COPY go.mod go.sum ./

# Configure Go proxy and download dependencies with verbose output
RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go env -w GOSUMDB=sum.golang.org && \
    go mod download -x

# Copy source code and build statically linked binary
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/app ./main.go

# --- STAGE 2: Runtime Environment ---
FROM ubuntu:22.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1

# Set working directory
WORKDIR /app

# --- 1. Install Base Tools ---
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
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates

# --- 2. Install Go ---
ARG GOLANG_VERSION=1.21.5
RUN curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# --- 3. Install Node.js + TypeScript ---
ARG NODE_VERSION=20.x
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - && \
    apt-get install -y nodejs && \
    npm install -g typescript ts-node @types/node

# --- 4. Install Julia ---
ARG JULIA_VERSION=1.10.0
RUN JULIA_MAJOR_MINOR=$(echo $JULIA_VERSION | cut -d. -f1,2) && \
    curl -fLO https://julialang-s3.s3.amazonaws.com/bin/linux/x64/${JULIA_MAJOR_MINOR}/julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    tar -C /opt -xzf julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    ln -s /opt/julia-${JULIA_VERSION}/bin/julia /usr/local/bin/julia && \
    rm julia-${JULIA_VERSION}-linux-x64.tar.gz

# --- 5. Copy project files FIRST ---
COPY . .

# --- 6. Install Python dependencies if requirements.txt exists ---
RUN if [ -f requirements.txt ]; then \
        pip3 install --no-cache-dir -r requirements.txt; \
    fi

# --- 7. Install Node.js dependencies if package.json exists ---
RUN if [ -f package.json ]; then \
        npm install; \
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

# --- 11. Expose port (adjust as needed) ---
EXPOSE 8080

# --- 12. Default command ---
CMD ["/usr/local/bin/app"]
