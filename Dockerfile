
 # --- STAGE 1: Builder Stage (for Go) ---
FROM golang:1.21-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod file first
COPY go.mod .

# Download dependencies and generate/update go.sum.
# This step works even if go.sum is missing or incomplete in the repo.
# We also run 'go mod tidy' to ensure all dependencies are correct.
RUN go mod download && go mod tidy

# Copy the rest of the source code
COPY . .

# Build the application
# -o app: names the output binary 'app'
# CGO_ENABLED=0: creates a statically linked binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/app ./main.go

# Use a robust and widely supported base image for a multi-language environment
FROM ubuntu:22.04

# Set environment variables for non-interactive installs and Python
ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1

# --- 1. Install Base Tools and Dependencies ---
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    git \
    ca-certificates \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# --- 2. Install Go ---
ENV GOLANG_VERSION 1.21.5
RUN curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# --- 3. Install Python and Pip ---
RUN apt-get update && \
    apt-get install -y python3 python3-pip && \
    ln -sf /usr/bin/python3 /usr/bin/python && \
    pip install --upgrade pip

# --- 4. Install Node.js (for TypeScript/TSX) ---
# Add NodeSource repository for a modern Node.js version
ENV NODE_VERSION 20.x
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - && \
    apt-get install -y nodejs
# Install TypeScript globally for running TSX/TS files
RUN npm install -g typescript ts-node

# --- 5. Install Julia (FIXED VERSION) ---
# NOTE: Using a recent, stable version (1.10.0) which is known to be available.
ENV JULIA_VERSION 1.10.0
RUN curl -LO https://julialang-s3.s3.amazonaws.com/bin/linux/x64/${JULIA_VERSION%.*}/julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    tar -C /opt -xzf julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    ln -s /opt/julia-${JULIA_VERSION}/bin/julia /usr/local/bin/julia && \
    rm julia-${JULIA_VERSION}-linux-x64.tar.gz

# --- 6. Final Setup: Copy Code and Install Dependencies ---
WORKDIR /app
# Copy everything from your local repo into the container
COPY . /app

# Install Python requirements (if file exists)
RUN if [ -f requirements.txt ]; then pip install -r requirements.txt; fi

# Set the default command to open an interactive shell when the container starts
CMD ["/bin/bash"]
