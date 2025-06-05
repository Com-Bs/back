# API Tests

This directory contains comprehensive tests for the Learning Go API.

## Quick Start

### 1. Set up MongoDB Authentication

If your MongoDB requires authentication, create a `.env` file:

```bash
# Copy the example file
cp tests/env.example tests/.env

# Edit with your MongoDB credentials
# Option 1: Complete URI (recommended)
TEST_MONGO_URI=mongodb://username:password@localhost:27017/?authSource=admin

# Option 2: Individual settings
# MONGO_USERNAME=your_username
# MONGO_PASSWORD=your_password
```

### 2. Run Tests

```bash
# Run all tests
./tests/run_tests.sh

# Check prerequisites first
./tests/run_tests.sh check
```

## Structure

```
tests/
├── unit/               # Unit tests for individual handlers
│   ├── user_handler_test.go
│   ├── logs_handler_test.go
│   ├── compile_handler_test.go
│   └── problem_handler_test.go
├── integration/        # Integration tests for complete workflows
│   ├── auth_integration_test.go
│   └── problems_integration_test.go
├── test_config.go      # Centralized test configuration
├── env.example         # Environment configuration template
└── README.md
```

## Test Coverage

### Unit Tests

- **user_handler_test.go**: Tests for `/signUp` and `/logIn` endpoints

  - Valid registration and login
  - Invalid JSON payloads
  - Missing required fields
  - Duplicate user registration
  - Invalid credentials
  - JWT token generation

- **logs_handler_test.go**: Tests for `/logs` endpoint

  - Valid GET requests with authentication
  - Invalid HTTP methods
  - Empty logs database scenarios

- **compile_handler_test.go**: Tests for `/compile` endpoint

  - Valid compile requests
  - Invalid HTTP methods
  - Cache functionality testing
  - External service integration

- **problem_handler_test.go**: Tests for `/problems` and `/problems/{id}` endpoints
  - Get all problems
  - Get specific problem by ID
  - Invalid ID formats
  - Non-existent problems
  - Empty database scenarios

### Integration Tests

- **auth_integration_test.go**: Complete authentication flow

  - User signup → login → access protected endpoints
  - Authentication middleware testing
  - JWT token validation
  - Database logging middleware

- **problems_integration_test.go**: Problems endpoints with full middleware stack
  - Authentication required for all problem endpoints
  - Complete CRUD operations
  - Error handling with proper HTTP status codes

## Prerequisites

1. **MongoDB**: Running MongoDB instance (with or without authentication)
2. **Go 1.23+**: Required for the test framework features used

## MongoDB Configuration

### Option 1: MongoDB without Authentication (Development)

```bash
# No additional configuration needed
# Tests will connect to mongodb://localhost:27017
```

### Option 2: MongoDB with Authentication (Production-like)

```bash
# Create tests/.env file
TEST_MONGO_URI=mongodb://username:password@localhost:27017/?authSource=admin

# Or use individual settings
MONGO_USERNAME=your_username
MONGO_PASSWORD=your_password
MONGO_HOST=localhost
MONGO_PORT=27017
```

### Option 3: Docker MongoDB for Testing

```bash
# Run MongoDB in Docker without authentication
docker run -d --name mongo-test -p 27017:27017 mongo:latest

# Or with authentication
docker run -d --name mongo-test \
  -e MONGO_INITDB_ROOT_USERNAME=testuser \
  -e MONGO_INITDB_ROOT_PASSWORD=testpass \
  -p 27017:27017 mongo:latest

# Then set in .env:
TEST_MONGO_URI=mongodb://testuser:testpass@localhost:27017/?authSource=admin
```

## Running Tests

### Using the Test Runner Script (Recommended)

```bash
# Check prerequisites and MongoDB connection
./tests/run_tests.sh check

# Run all tests
./tests/run_tests.sh

# Run specific test types
./tests/run_tests.sh unit
./tests/run_tests.sh integration

# Run with coverage report
./tests/run_tests.sh coverage

# Clean test databases
./tests/run_tests.sh clean

# Show help
./tests/run_tests.sh help
```

### Using Make

```bash
cd tests
make check      # Check prerequisites
make test       # Run all tests
make unit       # Unit tests only
make integration # Integration tests only
make coverage   # Tests with coverage
make clean      # Clean test databases
```

### Direct Go Commands

```bash
cd api

# Run all tests
go test ./tests/... -v

# Run with environment variables
TEST_MONGO_URI='mongodb://user:pass@localhost:27017/?authSource=admin' go test ./tests/... -v

# Run specific test types
go test ./tests/unit/... -v
go test ./tests/integration/... -v

# Run with coverage
go test ./tests/... -v -cover
```

## Test Database

- Tests use separate test databases to avoid interfering with development data
- Database names: `test_learning_go_*`
- Collections are automatically cleaned up before each test
- Tests will skip if MongoDB connection fails

## Authentication Testing

Tests use the same JWT secret key as the application (`"secret-key"`) and generate valid tokens for testing protected endpoints. Each test creates a standard test user:

- **Username**: `testuser`
- **Email**: `test@example.com`
- **Password**: `testpassword`

## Mock External Services

The compile endpoint tests may fail if the external compile service at `http://10.49.12.48:3001/runCompile` is not available. This is expected in testing environments.

## Test Data

Tests create their own test data including:

- Test users with known credentials
- Test problems with defined test cases
- Test logs for middleware verification

## Expected Behavior

### Success Cases

- User registration and login with valid data
- Access to protected endpoints with valid JWT tokens
- Proper JSON responses with expected data structures
- Database operations working correctly

### Failure Cases

- Appropriate HTTP status codes for errors
- Proper error messages in response bodies
- Authentication failures for missing/invalid tokens
- Validation errors for malformed requests

## Troubleshooting

### MongoDB Connection Issues

```bash
# Test your MongoDB connection manually
mongosh "mongodb://username:password@localhost:27017/?authSource=admin"

# Check if MongoDB is running
sudo systemctl status mongod  # Linux
brew services list | grep mongodb  # macOS
```

### Common Errors

1. **"Command insert requires authentication"**

   - MongoDB requires authentication but credentials not provided
   - Create `.env` file with correct MongoDB credentials

2. **"MongoDB not available for testing"**

   - MongoDB service is not running
   - Wrong connection credentials
   - Network connectivity issues

3. **"Tests are skipped"**
   - MongoDB connection failed
   - Tests automatically skip when database is unavailable

### Environment Variables

```bash
# Debug environment loading
echo $TEST_MONGO_URI
echo $MONGO_USERNAME

# Run tests with debug output
TEST_MONGO_URI='mongodb://localhost:27017' ./tests/run_tests.sh check
```

## Debugging Tests

If tests fail:

1. Check MongoDB connection with `./tests/run_tests.sh check`
2. Verify test database permissions
3. Check for port conflicts (27017)
4. Review test output for specific error messages
5. Ensure all dependencies are installed
6. Check `.env` file configuration

## Contributing

When adding new endpoints:

1. Add unit tests for the handler function
2. Add integration tests for the complete flow
3. Include both success and failure scenarios
4. Test authentication if the endpoint requires it
5. Use the centralized test configuration
6. Update this README if needed
