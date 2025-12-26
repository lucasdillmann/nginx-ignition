FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Install nginx from official repository to get version 1.25.1+ with http2 directive support
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    curl \
    gnupg2 \
    ca-certificates \
    lsb-release \
    ubuntu-keyring \
    && curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor \
    | tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null \
    && echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] \
    http://nginx.org/packages/ubuntu `lsb_release -cs` nginx" \
    | tee /etc/apt/sources.list.d/nginx.list \
    && apt-get update \
    && apt-get -y install --no-install-recommends \
    build-essential \
    wget \
    git \
    postgresql-client \
    nginx \
    && mkdir -p /usr/lib/nginx/modules \
    && mkdir -p /var/cache/nginx/client_temp \
    && mkdir -p /var/cache/nginx/proxy_temp \
    && mkdir -p /var/cache/nginx/fastcgi_temp \
    && mkdir -p /var/cache/nginx/uwsgi_temp \
    && mkdir -p /var/cache/nginx/scgi_temp \
    && chown -R vscode:vscode /var/cache/nginx \
    && chown -R vscode:vscode /var/log/nginx \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install Go 1.23
ARG GO_VERSION=1.23.5
RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-arm64.tar.gz" -O /tmp/go.tar.gz \
    && tar -xzf /tmp/go.tar.gz -C /usr/local \
    && rm /tmp/go.tar.gz \
    && ln -s /usr/local/go/bin/go /usr/local/bin/go \
    && ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt

# Install Node.js 22
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@11 \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Set up Go environment for vscode user
ENV GOPATH=/home/vscode/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Install Go development tools as vscode user
USER vscode
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install mvdan.cc/gofumpt@latest \
    && go install github.com/air-verse/air@latest

# Set up workspace
WORKDIR /workspaces/nginx-ignition

# Expose ports for development
# 8090 - Backend API
# 5173 - Vite dev server
# 5432 - PostgreSQL (if accessing from host)
EXPOSE 8090 5173 5432
