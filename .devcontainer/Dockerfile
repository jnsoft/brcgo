
# Use the official Go image as a parent image
FROM docker.io/golang:latest

# Install necessary tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    unzip \
    graphviz \
    && rm -rf /var/lib/apt/lists/*

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install Go tools
RUN go install -v golang.org/x/tools/cmd/goimports@latest

# Install gopls
RUN go install -v golang.org/x/tools/gopls@latest

# Create a non-root user to use the container
RUN useradd -ms /bin/bash vscode

RUN mkdir -p /go && chown -R vscode:vscode /go

# Set the user to vscode
USER root

# Set the working directory
WORKDIR /workspace

# Expose the port for Delve debugger
EXPOSE 2345
