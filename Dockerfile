# Use a robust base image like Ubuntu
FROM ubuntu:22.04

# Set environment variables for non-interactive installs
ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1

# --- 1. Install Base Tools and Dependencies ---
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    git \
    vim \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# --- 2. Install Go (Using Go's official package) ---
ENV GOLANG_VERSION 1.21.1
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
# Use NodeSource to get a modern Node.js version
ENV NODE_VERSION 20.x
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - && \
    apt-get install -y nodejs
RUN npm install -g typescript ts-node

# --- 5. Install Julia (via its official download) ---
ENV JULIA_VERSION 1.9.4
RUN curl -LO https://julialang-s3.s3.amazonaws.com/bin/linux/x64/${JULIA_VERSION%.*}/julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    tar -C /opt -xzf julia-${JULIA_VERSION}-linux-x64.tar.gz && \
    ln -s /opt/julia-${JULIA_VERSION}/bin/julia /usr/local/bin/julia && \
    rm julia-${JULIA_VERSION}-linux-x64.tar.gz

# --- 6. Final Setup: Copy Code and Install Python Dependencies ---
WORKDIR /app
COPY . /app

# Install Python requirements (assuming you have a requirements.txt)
RUN if [ -f requirements.txt ]; then pip install -r requirements.txt; fi

# Set the default command to open a shell/terminal for testing/running scripts
CMD ["/bin/bash"]
