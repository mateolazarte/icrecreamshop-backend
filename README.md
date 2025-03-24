# 🚀 Ice Cream Shop API with Go

---
## 🍦 About the Project
REST API built in Go for ice cream shop order management. Powered by the Gin framework, it provides secure endpoints for both customers and administrators, featuring JWT authentication and Swagger API documentation.

The application uses Docker and PostgreSQL for database management, and integrates with Mercado Pago for payment.

---
## 📦 Dependencies

- **Go 1.23.4+**: Required to build and run the project.
- **Docker**: Used for running PostgreSQL.
---
## 🛠 Quick Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mateolazarte/icecreamshop-backend.git
   cd icecreamshop-backend
   ```

2. **Configure environment**

    ```bash
    cp .env.example .env
    # Edit .env with your credentials
   ```

3. **Start services with Docker**

   ```bash
   docker-compose up -d --build
    ```

## 🔌 Environment Variables

- 
    ```.env
    # Database
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=postgres
    
    # Testing Database
    TEST_DB_HOST=localhost
    TEST_DB_PORT=5433
    TEST_DB_USER=test_user
    TEST_DB_PASSWORD=test_password
    TEST_DB_NAME=app_test
    
    # API
    API_PORT=8080
    API_ENV=testing
    
    # Secrets
    JWT_SECRET=your-secret
    
    # Tests (mock or integration)
    TEST_MODE=mock 
    
    # MercadoPago
    MP_ACCESS_TOKEN=your_mp_token
    ```
---
## 💻 Run local

1. **Install dependencies**
    ```bash
    go mod tidy
    ```

2. **Start the server**
    ```bash
    go run ./cmd
    # Set env variable API_ENV=development to turn Gin logs on and use the Database
    # Set env variable API_ENV=testing to turn Gin logs off and use local memory
    ```
---
## ✅ Running Tests
-   ```bash
    cd internal/tests
    go test
    # Set env variable TEST_MODE=integration for integration tests.
    # Set env variable API_ENV=testing to turn Gin logs off.
    ```
---
## 📚 API Documentation

The API is fully documented using Swagger (OpenAPI). You can find the documentation in the `swagger` folder at `icecreamshop-backend/internal/api/swagger/swagger.yaml`.

---
Built by Mateo Lazarte and Esteban Mena