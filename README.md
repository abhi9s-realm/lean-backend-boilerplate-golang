🧠 Lean Backend Boilerplate - Go/Gin
A minimal, production-ready backend API boilerplate built with Go and Gin. Get your project started in minutes, not hours.

🎯 Philosophy

Lean by Design: Only the essentials - no feature bloat
5-Minute Setup: Clone, configure, run
Extend Don't Remove: Add what you need, when you need it
Production Ready: Proper error handling, logging, and structure from day one


📁 Project Structure
```
lean-backend-boilerplate-golang/
├── api/                      # HTTP layer
│   ├── handlers/
│   │   ├── health.go        # Health check endpoint
│   │   └── user.go          # Example user handler
│   ├── middleware/
│   │   ├── cors.go          # CORS middleware
│   │   └── logger.go        # Request logging
│   └── routes/
│       └── routes.go        # Route definitions
├── cmd/
│   └── api/
│       └── main.go          # Application entrypoint
├── config/
│   └── config.go            # Environment configuration
├── internal/
│   ├── domain/
│   │   ├── models/
│   │   │   └── user.go      # Example user model
│   │   └── services/        # Business logic layer
│   └── infrastructure/
│       ├── database/
│       │   └── postgres.go  # Database connection
│       └── logger/
│           └── logger.go    # Structured logging
├── pkg/
│   └── utils/
│       └── response.go      # API response helpers
├── tests/
│   ├── api/
│   │   ├── health_test.go   # API tests
│   │   └── user_test.go     # User endpoint tests
│   └── testutils/
│       └── setup.go         # Test helpers and utilities
├── .env                     # Base configuration
├── docker-compose.yml       # Container orchestration
├── Dockerfile              # Container definition
├── go.mod                  # Go modules file
└── README.md              # Project documentation
```


🚀 Quick Start
Prerequisites

Go 1.21+
PostgreSQL (or use Docker)

Setup (< 5 minutes)
# 1. Clone
git clone https://github.com/your-username/lean-backend-boilerplate-golang.git
cd lean-backend-boilerplate-golang

# 2. Install Dependencies
go mod tidy

# 3. Configure Environment
# Base .env file is included with documentation
# Create environment-specific files as needed:
cp .env .env.dev    # Development settings
cp .env .env.test   # Test settings
cp .env .env.prod   # Production settings

# 4. Run
make run           # Development mode (default)
# or
make run-dev      # Explicit development mode
make run-prod     # Production mode

API ready at http://localhost:8080 ✅

⚙️ Environment Configuration
We support multiple environments through separate config files:

# Base configuration (.env)
ENVIRONMENT=development    # development, test, or production
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=myapp
LOG_LEVEL=info

Environment-specific files (.env.dev, .env.test, .env.prod) inherit from base config.
Only .env is versioned - other files are gitignored for security.


📡 Default Endpoints



Method
Endpoint
Description



GET
/api/health
Health check


GET
/api/users
List users (example)


POST
/api/users
Create user (example)


Response Format:
{
  "success": true,
  "data": {...},
  "message": "Success"
}


🛠️ Essential Commands
# Running
make run          # Start server (development mode)
make run-dev      # Start server in development mode
make run-prod     # Start server in production mode

# Testing
make test         # Run tests (uses test environment)
make test-dev     # Run tests in development mode
make test-prod    # Run tests in production mode

# Building
make build        # Build binary
make docker-build # Build Docker image


🧰 What's Included
Core Features:

🚀 Gin HTTP server with middleware (Gin v1.9.1)
🗄️ PostgreSQL with GORM (GORM v1.25.4)
📝 Structured logging (Zap v1.26.0)
⚙️ Multi-environment configuration (Viper v1.18.2)
🔒 CORS handling
📦 Standardized API responses
🧪 Comprehensive testing setup (Testify v1.8.4)
🐳 Docker & Docker Compose ready

What's NOT Included (Add When Needed):

Authentication/JWT
Rate limiting
Redis/caching
Database migrations
Complex validation
Metrics/monitoring
Kubernetes configs


🏗️ Architecture
HTTP Request → Middleware → Handler → Service → Repository → Database


Handlers: HTTP request/response handling
Services: Business logic
Models: Data structures
Infrastructure: External dependencies (DB, logger, etc.)


🐳 Docker
We provide both Dockerfile and docker-compose.yml for containerization:

# Using Docker directly:
docker build -t lean-backend-boilerplate-golang .
docker run -p 8080:8080 \
    -e ENVIRONMENT=production \
    -e DB_HOST=postgres \
    lean-backend-boilerplate-golang

# Using Docker Compose (recommended):
docker-compose up     # Starts both API and PostgreSQL
docker-compose down   # Stop all services

The docker-compose setup includes:
- API service with hot-reload
- PostgreSQL database
- Persistent database volume
- Proper service networking


🧪 Testing
# Run all tests
make test

# Test specific package
go test ./api/handlers -v

Basic test structure included - expand as needed.

🚀 Extension Points
Need Authentication? Add JWT middleware in api/middleware/
Need Database Migrations? Add golang-migrate and migration files
Need Caching? Add Redis client in infrastructure/
Need API Docs? Add Swagger/OpenAPI annotations
Need More Validation? Extend with go-playground/validator

🤝 Contributing

Fork the repo
Add your feature
Write tests
Submit PR

Keep it lean - new features should be opt-in, not default.

📄 License
MIT License - Build something awesome! 🚀

Ready to code in 5 minutes. Scale when you need to.
Updated as of June 7, 2025
