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

# --- 1. Install Base Tools ---
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    git \
    ca-certificates \
    software-properties-common \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# --- 2. Install Go ---
ARG GOLANG_VERSION=1.21.5
RUN curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# --- 3. Install Node.js + TypeScript ---
ARG NODE_VERSION=20.x
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - && \
    apt-get install -y nodejs && \
    npm install -g typescript ts-node

# --- 4. Install Julia ---
ARG JULIA_VERSION=1.10.0
RUN JULIA_MAJOR_MINOR=$(echo $JULIA_VERSION | cut -d. -f1,2) && \
    curl -fLO https://julialang-s3.s3.amazonaws.com/bin/linux/x64/${JULIA_MAJOR_MINOR}/julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    tar -C /opt -xzf julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    ln -s /opt/julia-${JULIA_VERSION}/bin/julia /usr/local/bin/julia && \
    rm julia-${JULIA_VERSION}-linux-x64.tar.gz

# --- 5. Final Setup ---
WORKDIR /app

# Copy all project files
COPY . /app

# Install Python dependencies if requirements.txt exists
RUN if [ -f requirements.txt ]; then pip install -r requirements.txt; fi

# Copy Go binary from builder stage
COPY --from=builder /app/app /usr/local/bin/app

# Default command
CMD ["/bin/bash"]
