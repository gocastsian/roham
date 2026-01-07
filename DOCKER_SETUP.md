# Docker Compose Setup Guide

This guide explains how to run the Roham project locally with Docker Compose and Traefik.

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Go 1.24+ (for local development; Docker uses Go 1.25 for air compatibility)
- `roham.local` added to `/etc/hosts` (see below)

### Step 1: Configure `/etc/hosts`

Add this line to your `/etc/hosts` file:
```
127.0.0.1 roham.local
```

### Step 2: Configure Environment File

The project uses a centralized `deploy/.env` file for all configuration. 

**First-time setup:**
```bash
# Copy the example file (if it doesn't exist)
cp deploy/.env.example deploy/.env
```

**Edit `deploy/.env` to customize:**
- Database credentials (usernames, passwords)
- Service ports
- Domain names
- Go build configuration

The `.env` file contains:
- **Docker Compose variables**: Used in docker-compose.yaml files (e.g., `${USER_DB_NAME}`, `${FILER_DB_PASSWORD}`)
- **Application variables**: Used by Go services with prefixes (`USER__POSTGRES_DB__USER`, `FILER__POSTGRES_DB__PASSWORD`, etc.)

All database credentials, ports, and service configurations are centralized in this single file. See the [Environment Variables](#environment-variables) section for details.

### Step 3: Run All Services

```bash
# Using Make (recommended)
make start-all

# Or directly using the script
./deploy/docker-compose-dev.bash --profile traefik --profile user --profile filer --profile vectorlayer up --build
```

### Step 4: Access Services

- **Traefik Dashboard**: http://localhost:8080
- **User Service**: http://roham.local/users
- **Filer Service**: http://roham.local/filer
- **Vector Layer Service**: http://roham.local/vectorlayer

---

## How Traefik Works in This Project

### Architecture Overview

Traefik acts as a reverse proxy and API gateway that routes requests to your microservices:

```
Client Request → Traefik (Port 80) → Service (user/filer/vectorlayer)
```

### Configuration Files

1. **`deploy/roham/development/traefik/traefik.yml`** - Main Traefik configuration
   - Entry point: Port 80 (web)
   - Dashboard: Port 8080 (insecure mode)
   - Docker provider enabled (auto-discovers services via labels)
   - File provider for dynamic configuration

2. **`deploy/roham/development/traefik/dynamic/dynamic.yml`** - Dynamic middleware configuration
   - `roham_auth`: Authentication middleware (validates JWT via `/v1/auth` endpoint)
   - `roham_authz`: Authorization middleware (checks policies via `/v1/authz` endpoint)

### How Services Register with Traefik

Each service uses Docker labels to configure Traefik routing:

```yaml
labels:
  - "traefik.enable=true"  # Enable Traefik for this service
  - "traefik.http.routers.service_name.rule=Host(`roham.local`) && PathPrefix(`/service`)"
  - "traefik.http.routers.service_name.entrypoints=web"
  - "traefik.http.services.service_name.loadbalancer.server.port=5001"
  - "traefik.http.middlewares.service_name_strip.stripprefix.prefixes=/service"
```

### Authentication Flow

1. **Unauthenticated Routes** (e.g., `/users/v1/health-check`, `/users/v1/users/login`)
   - Traefik routes directly to the service without auth middleware

2. **Authenticated Routes** (e.g., `/users/*` except login/health)
   - Request → `roham_auth` middleware → Validates JWT at `user-service:5001/v1/auth`
   - If valid, adds `X-User-Info` header
   - Request → `roham_authz` middleware → Checks policies at `user-service:5001/v1/authz`
   - If authorized, request proceeds to target service

### Path Prefix Stripping

Traefik strips path prefixes before forwarding to services:
- Request: `http://roham.local/users/v1/users` 
- Traefik strips `/users` → Forwards to service as `/v1/users`

---

## Issues and Anti-Patterns Found

### 🔴 Critical Issues

#### 1. **Network Configuration Missing**
**Location**: `deploy/roham/development/traefik/traefik.yml:20`
**Issue**: Traefik references `network: traefik_default` but no network is defined in docker-compose files.
**Impact**: Services may not be able to communicate with Traefik or each other.
**Fix**: Define a shared network in all docker-compose files and reference it explicitly.

#### 2. **Incorrect Dockerfile Path**
**Location**: `deploy/vectorlayer/development/docker-compose.yaml:26`
**Issue**: References `dockerfile: ./Dockerfile` but the actual file is at `./deploy/vectorlayer/development/Dockerfile`
**Impact**: Build will fail.
**Fix**: Change to `dockerfile: ./deploy/vectorlayer/development/Dockerfile`

#### 3. **Port Conflicts**
**Issue**: Both `user-db` and `filer-db` use port `5433`
**Impact**: Only one database can run at a time.
**Fix**: Use different ports (already done: user-db=5433, filer-db=5433, but vectorlayer-db=5434)

### 🟡 Medium Priority Issues

#### 4. **Inconsistent Environment File Paths**
**Issue**: 
- `user-service` and `vectorlayer-service` use `.env` (project root)
- `filer-service` uses `./deploy/.env`
- Script uses `--env-file ./deploy/.env`
**Impact**: Confusion and potential configuration errors.
**Recommendation**: Standardize on `./deploy/.env` for all services.

#### 5. **Insecure Traefik Dashboard**
**Location**: `deploy/roham/development/traefik/traefik.yml:11`
**Issue**: `api.dashboard.insecure: true` exposes dashboard without authentication
**Impact**: Security risk in production (acceptable for dev).
**Recommendation**: Add basic auth or disable in production.

#### 6. **Unnecessary Port Exposures**
**Location**: `deploy/filer/development/docker-compose.yaml:21-23`
**Issue**: Filer service exposes ports 5005 and 5006 directly
**Impact**: Bypasses Traefik routing, inconsistent with other services.
**Recommendation**: Remove direct port exposure if using Traefik exclusively.

#### 7. **Missing Network Definition**
**Issue**: No explicit network defined for services to communicate
**Impact**: Docker Compose creates default network, but explicit definition is better practice.
**Recommendation**: Define a shared network explicitly.

---

## Recommendations

### 1. Create a Shared Network

Add to all docker-compose files:
```yaml
networks:
  roham_network:
    name: roham_network
    driver: bridge

services:
  service-name:
    networks:
      - roham_network
```

Update `traefik.yml`:
```yaml
providers:
  docker:
    network: roham_network
```

### 2. Fix Dockerfile Path

In `deploy/vectorlayer/development/docker-compose.yaml`:
```yaml
build:
  context: .
  dockerfile: ./deploy/vectorlayer/development/Dockerfile
```

### 3. Standardize Environment Files

Update `deploy/user/development/docker-compose.yaml` and `deploy/vectorlayer/development/docker-compose.yaml`:
```yaml
env_file:
  - ./deploy/.env
```

### 4. Add Health Checks

Add health checks to services:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:5001/v1/health-check"]
  interval: 30s
  timeout: 10s
  retries: 3
```

### 5. Use Docker Compose Profiles Correctly

The current setup uses profiles effectively for service isolation. Consider:
- Adding a `traefik` profile explicitly
- Documenting which profiles are needed for which scenarios

### 6. Environment-Specific Configuration

Consider:
- Separate `docker-compose.dev.yaml` and `docker-compose.prod.yaml`
- Use environment variables for sensitive values
- Add `.env.example` file for documentation

### 7. Add Traefik Network Labels

Ensure all services that need Traefik routing are on the same network:
```yaml
services:
  traefik:
    networks:
      - roham_network
    labels:
      - "traefik.docker.network=roham_network"
```

### 8. Security Improvements for Production

- Add authentication to Traefik dashboard
- Use HTTPS/TLS certificates
- Remove insecure dashboard configuration
- Use Docker secrets for passwords

---

## Troubleshooting

### Services Not Reachable Through Traefik

1. Check if services are on the same Docker network:
   ```bash
   docker network ls
   docker network inspect roham_network
   ```

2. Verify Traefik can see services:
   - Check Traefik dashboard: http://localhost:8080
   - Look for services in "HTTP" → "Routers" section

3. Check service labels are correct

4. Verify `/etc/hosts` has `roham.local` entry

### Port Already in Use

If you get port conflicts:
- Check what's using the port: `lsof -i :80` or `netstat -tulpn | grep :80`
- Stop conflicting services or change ports in docker-compose files

### Build Failures

1. Check Dockerfile paths are correct
2. Verify base images exist: `docker pull golang:1.24.0-alpine3.21`
3. Check build context paths

---

## Useful Commands

```bash
# Start all services
make start-all

# Start specific service
./deploy/docker-compose-dev.bash --profile traefik --profile user up

# View logs
./deploy/docker-compose-dev.bash logs -f user-service

# Stop all services
./deploy/docker-compose-dev.bash down

# Rebuild specific service
./deploy/docker-compose-dev.bash --profile user up --build user-service

# Check Traefik configuration
docker exec <traefik-container> traefik version
docker exec <traefik-container> cat /etc/traefik/traefik.yml
```

---

## Environment Variables

All configuration is centralized in `deploy/.env`. This file contains:

### Docker Compose Variables

These variables are used in docker-compose.yaml files:

```bash
# Service Configuration
SERVICE_NAME=roham
SERVICE_DOMAIN=roham.local

# User Service Database
USER_DB_NAME=user_db
USER_DB_USER=user_admin
USER_DB_PASSWORD=password123
USER_DB_PORT=5433

# Filer Service Database
FILER_DB_NAME=filer_db
FILER_DB_USER=roham
FILER_DB_PASSWORD=secret
FILER_DB_PORT=5433

# Vectorlayer Service Database
VECTORLAYER_DB_NAME=vectorlayer_db
VECTORLAYER_DB_USER=vectorlayer_admin
VECTORLAYER_DB_PASSWORD=password123
VECTORLAYER_DB_PORT=5434
```

### Application Environment Variables

These variables are automatically loaded by the Go services using the `cfg_loader` package:

- **User Service**: Variables prefixed with `USER__` (e.g., `USER__POSTGRES_DB__USER`)
- **Filer Service**: Variables prefixed with `FILER__` (e.g., `FILER__POSTGRES_DB__PASSWORD`)
- **Vectorlayer Service**: Variables prefixed with `VECTORLAYER__` (e.g., `VECTORLAYER__POSTGRES_DB__HOST`)

The double underscore (`__`) is converted to a dot (`.`) in the configuration structure. For example:
- `USER__POSTGRES_DB__USER` → `postgres_db.user`
- `FILER__POSTGRES_DB__PASSWORD` → `postgres_db.password`

**How it works:**
1. Docker Compose loads variables from `deploy/.env` using `--env-file`
2. Services use `env_file: - ./deploy/.env` to pass variables to containers
3. Applications use `cfg_loader` which reads YAML config files AND environment variables
4. Environment variables override YAML values, allowing you to change config without modifying YAML files

### Creating Your .env File

```bash
# Use the example template
cp deploy/.env.example deploy/.env

# Edit with your values
nano deploy/.env  # or your preferred editor
```

**Security Note:** Never commit `.env` to version control. The `.env.example` file is a template without sensitive values.

---

## Additional Notes

- The project uses **Docker Compose profiles** to isolate services
- Traefik uses **Docker provider** to auto-discover services via labels
- Services use **path-based routing** with prefix stripping
- Authentication is handled via **forwardAuth** middleware pattern
- All configuration is centralized in `deploy/.env` for easy management

