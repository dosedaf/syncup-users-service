# Summary of Learnings: SyncUp Users Service

This document summarizes the key concepts and software engineering patterns learned and applied during the development of the SyncUp Users Service.

---

### 1. Software Architecture & Design Patterns

* **Layered Architecture:** Built the application using a clean three-layer architecture (Handler, Service, Repository) to ensure separation of concerns. The `Handler` manages HTTP requests, the `Service` holds the business logic, and the `Repository` handles database interactions.
* **Dependency Injection (DI):** Injected all dependencies (like the database connection, logger, and other layers) at application startup. This makes the application modular and testable.
* **Programming to Interfaces:** Decoupled the layers by having them depend on interfaces (`ServiceInstance`, `RepositoryInstance`) rather than concrete structs. This is the key pattern that enables isolated unit testing.
* **Microservice Principles:** Designed the service as a self-contained unit with its own database, understanding the distinction from a monolithic architecture where all components share one database.

### 2. Advanced Error Handling

* **Error Wrapping:** Used `fmt.Errorf` with the `%w` verb in the service layer to add valuable context to errors returned from the repository. This creates descriptive error logs that pinpoint the exact location and context of a failure.
* **Sentinel Errors:** Created and used specific, exported error variables (e.g., `helper.ErrUserNotFound`, `helper.ErrEmailAlreadyExists`) to represent known business rule failures.
* **Error Checking with `errors.Is`:** Used `errors.Is()` in the handler and tests to reliably check for specific sentinel errors, even after they have been wrapped with additional context.
* **Error Translation:** Implemented the pattern of translating low-level, dependency-specific errors (like `pgx.ErrNoRows`) into high-level, application-specific sentinel errors (like `helper.ErrUserNotFound`) at the correct architectural boundary (the repository).

### 3. Testing

* **Unit Testing in Isolation:** Wrote unit tests for the service layer that run completely isolated from the database.
* **Mocking Dependencies:** Created mock repositories that implement the `RepositoryInstance` interface to simulate database behavior during tests.
* **Configurable Mocks:** Designed mocks with function fields, allowing their behavior to be configured on a per-test basis to simulate both success ("happy path") and various failure ("sad path") scenarios.

### 4. Go & Backend Fundamentals

* **`context.Context` Propagation:** Passed `context.Context` through all application layers (handler, service, repository) to gracefully handle request timeouts and cancellations.
* **Structured Logging:** Implemented structured, key-value pair logging using Go's standard `slog` library for better observability.
* **Configuration Management:** Loaded sensitive data and configuration (database URL, JWT secret) from environment variables (`.env` file) at startup.
* **Password Hashing:** Secured user passwords using the `bcrypt` algorithm for one-way hashing and comparison.
* **JWT Authentication:** Generated stateless JSON Web Tokens (JWT) for users upon successful login.
* **Docker Compose for Development:** Used `docker-compose` to create a reproducible local development environment that includes the Go application and its PostgreSQL database.