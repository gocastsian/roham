# Roham

Roham is a microservices-based **map server** platform built with **Go**, providing geospatial **vector layer** processing.

It is composed of multiple services behind **Traefik** (reverse proxy + API gateway) and uses **Temporal** to orchestrate long-running geospatial workflows (imports, conversions, validation, etc.).

---

## 🏗️ Architecture

Roham is a microservices architecture consisting of three main services:

- **User Service**: user management, authentication, authorization
- **Filer Service**: file upload/download (resumable uploads)
- **Vector Layer Service**: geospatial vector layer processing and management (Temporal workflows)

---

## 🚀 Features

### User Service
- User registration and authentication
- JWT-based authentication with access/refresh tokens
- Role-based access control (RBAC) with OPA policies
- User profile management

### Filer Service
- Resumable file uploads using **TUS** protocol
- File downloads with direct and pre-signed URL support
- Multiple storage backends (Filesystem, S3/MinIO)
- File metadata management

### Vector Layer Service
- Vector layer import and processing
- Temporal-based workflow orchestration
- Layer management and styling
- Integration with Filer Service for shapefile processing

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Docker** and **Docker Compose** (v2.0+)
- **Go** 1.24 or higher (for local development; Docker uses Go 1.25)
- **Make** (optional, for convenience commands)

## 🚦 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd roham
```

### 2. Configure Hosts File

Add the following line to your `/etc/hosts` file (Linux/Mac) or `C:\Windows\System32\drivers\etc\hosts` (Windows):

```
127.0.0.1 roham.local
```

This allows you to access services via the domain name used in Traefik routing.

### 3. Set Up Environment Variables

Copy the example environment file and configure it:

```bash
cp deploy/.env.example deploy/.env
```

Edit `deploy/.env` to customize:
- Database credentials (usernames, passwords)
- Service ports
- Domain names
- Go build configuration

**Important**: All sensitive configuration (database credentials, passwords) should be set in the `.env` file. See [Configuration](#configuration) for details.

### 4. Start All Services

```bash
make start-all
```

Or manually:

```bash
./deploy/docker-compose-dev.bash --profile traefik --profile user --profile filer --profile vectorlayer up --build
```

### 5. Access Services

Once all services are running, you can access:

- **Traefik Dashboard**: http://localhost:8080
- **User Service**: http://roham.local/users
- **Filer Service**: http://roham.local/filer
- **Vector Layer Service**: http://roham.local/vectorlayer

## 📁 Project Structure

```
roham/
├── cmd/                    # Application entry points
│   ├── user/              # User service entry point
│   ├── filer/             # Filer service entry point
│   └── vectorlayer/       # Vector layer service entry point
├── deploy/                # Deployment configurations
│   ├── .env               # Environment variables (create from .env.example)
│   ├── .env.example       # Environment variables template
│   ├── docker-compose-dev.bash  # Docker Compose orchestration script
│   ├── user/              # User service deployment configs
│   ├── filer/             # Filer service deployment configs
│   ├── vectorlayer/       # Vector layer service deployment configs
│   └── roham/             # Traefik configuration
├── userapp/               # User service implementation
├── filer/                 # Filer service implementation
├── vectorlayerapp/        # Vector layer service implementation
├── pkg/                   # Shared packages and utilities
│   ├── cfg_loader/        # Configuration loader (YAML + env vars)
│   ├── postgresql/        # PostgreSQL database utilities
│   ├── redis/             # Redis cache utilities
│   ├── logger/            # Logging utilities
│   └── ...
├── adapter/               # External service adapters
│   ├── redis/             # Redis adapter
│   └── temporal/          # Temporal workflow adapter
└── Makefile               # Convenience commands
```

## ⚙️ Configuration

### Environment Variables

All configuration is centralized in `deploy/.env`. The file contains:

1. **Docker Compose Variables**: Used in docker-compose.yaml files
   - `USER_DB_NAME`, `USER_DB_USER`, `USER_DB_PASSWORD` - User service database
   - `FILER_DB_NAME`, `FILER_DB_USER`, `FILER_DB_PASSWORD` - Filer service database
   - `VECTORLAYER_DB_NAME`, `VECTORLAYER_DB_USER`, `VECTORLAYER_DB_PASSWORD` - Vector layer database
   - `SERVICE_NAME`, `SERVICE_DOMAIN` - Service identification

2. **Application Environment Variables**: Used by Go services
   - `USER__POSTGRES_DB__HOST`, `USER__POSTGRES_DB__USER`, `USER__POSTGRES_DB__PASSWORD` - User service config
   - `FILER__POSTGRES_DB__HOST`, `FILER__POSTGRES_DB__USER`, `FILER__POSTGRES_DB__PASSWORD` - Filer service config
   - `VECTORLAYER__POSTGRES_DB__HOST`, etc. - Vector layer service config

### Configuration Loading

Services use the `cfg_loader` package which:
1. Loads YAML configuration files (defaults)
2. Overrides with environment variables (environment variables take precedence)

**Important**: Sensitive values (passwords, database credentials) should only be set in the `.env` file, not in YAML config files.

See `DOCKER_SETUP.md` for detailed configuration documentation.

## 🐳 Running with Docker Compose

### Start All Services

```bash
make start-all
```

### Start Individual Services

```bash
# Start only user service
./deploy/docker-compose-dev.bash --profile traefik --profile user up

# Start only filer service
./deploy/docker-compose-dev.bash --profile traefik --profile filer up

# Start only vectorlayer service
./deploy/docker-compose-dev.bash --profile traefik --profile vectorlayer up
```

### View Logs

```bash
# All services
./deploy/docker-compose-dev.bash logs -f

# Specific service
./deploy/docker-compose-dev.bash logs -f user-service
./deploy/docker-compose-dev.bash logs -f filer-service
./deploy/docker-compose-dev.bash logs -f vectorlayer-service
```

### Stop Services

```bash
./deploy/docker-compose-dev.bash down
```

### Rebuild Services

```bash
./deploy/docker-compose-dev.bash --profile traefik --profile user --profile filer --profile vectorlayer up --build
```

## 🔧 Development

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific package
make test-user
make test-filer
make test-vectorlayer

# Run tests with coverage
make test-coverage

# Run E2E tests (requires services to be running)
make test-e2e
```

See `TESTING.md` for detailed testing documentation.

### Local Development (Without Docker)

For local development, you can run services directly with Go:

1. **Set up databases** (PostgreSQL and Redis) - either locally or using Docker
2. **Configure environment variables** in `deploy/.env`
3. **Run the service**:

```bash
# User service
go run cmd/user/main.go

# Filer service
go run cmd/filer/main.go serve

# Vector layer service
go run cmd/vectorlayer/main.go
```

## 🌐 API Gateway (Traefik)

Traefik acts as the API gateway and reverse proxy:

- **Routing**: Routes requests based on path prefixes (`/users`, `/filer`, `/vectorlayer`)
- **Authentication**: Validates JWT tokens via forwardAuth middleware
- **Authorization**: Checks policies using OPA via forwardAuth middleware
- **Path Stripping**: Removes service prefixes before forwarding to services

### Traefik Dashboard

Access the Traefik dashboard at http://localhost:8080 to:
- View active routes and services
- Monitor request metrics
- Debug routing issues

### Authentication Flow

1. **Public endpoints** (login, health checks) bypass authentication
2. **Protected endpoints** go through:
   - `roham_auth` middleware: Validates JWT token
   - `roham_authz` middleware: Checks authorization policies
   - Request is forwarded to the target service

## 📚 Documentation

- **[DOCKER_SETUP.md](DOCKER_SETUP.md)**: Detailed Docker Compose and Traefik setup guide
- **[TESTING.md](TESTING.md)**: Testing guide and best practices
- **[filer/Readme.md](filer/Readme.md)**: Filer service specific documentation

## 🛠️ Available Make Commands

```bash
# Start services
make start-all              # Start all services with Docker Compose
make start-filer-dev        # Start only filer service

# Testing
make test                   # Run all tests
make test-coverage          # Run tests with coverage
make test-coverage-html     # Generate HTML coverage report
make test-user              # Run user service tests
make test-filer             # Run filer service tests
make test-vectorlayer       # Run vector layer service tests
make test-e2e               # Run end-to-end tests
```

## 🔐 Security Notes

- **Never commit** the `.env` file to version control
- Database credentials are stored in `.env` file only
- JWT tokens are used for authentication
- OPA policies control authorization
- Traefik dashboard is insecure by default (development only)

For production deployments:
- Use secure Traefik dashboard authentication
- Enable HTTPS/TLS
- Use Docker secrets for sensitive data
- Secure database connections

## 🐛 Troubleshooting

### Services Not Reachable

1. Check if services are running: `docker ps`
2. Verify Traefik dashboard: http://localhost:8080
3. Check service logs: `./deploy/docker-compose-dev.bash logs -f <service-name>`
4. Verify `/etc/hosts` has `roham.local` entry

### Port Conflicts

If ports are already in use:
1. Check what's using the port: `lsof -i :80` or `netstat -tulpn | grep :80`
2. Stop conflicting services or change ports in `deploy/.env`

### Database Connection Issues

1. Verify database credentials in `deploy/.env`
2. Check database containers are running: `docker ps | grep postgres`
3. Verify network connectivity: services must be on the same Docker network

See `DOCKER_SETUP.md` for more troubleshooting tips.

## 📝 License

See [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

For issues and questions, please open an issue on the repository.
