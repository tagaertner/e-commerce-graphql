# E-Commerce GraphQL Microservices Platform

A production-style **microservices backend** built with **Go**, **GraphQL Federation**, and **PostgreSQL**, featuring a lightweight Python **Gradio UI** for end-to-end API validation and exploration. Deployed on Render.

This project demonstrates **backend architecture**, **service composition**, and **stateful system design** with real-world deployment considerations.

## Demo & Availability

This project has a hosted demo on Render. It can be slow or unavailable at times due to cold starts or rate limiting.

For a smoother experience, run the project locally with Docker Compose. See **Run Locally** below.

## Try it yourself

**Live apps**

- **Gradio UI:** https://gradio-ui-render-e-commerce-graphql.onrender.com
- **GraphQL Gateway:** https://gateway-render-e-commerce-graphql.onrender.com

---

## Video walkthroughs (with audio)

### E-Commerce UI Demo

https://github.com/user-attachments/assets/8176fa10-a4d4-4653-84bb-032f2cb8c933

### GraphQL Demo

https://github.com/user-attachments/assets/a79566af-e485-46d7-9b76-c01a3b7c6ba2

## Quick Start

### Local Development

**Prerequisites:** Docker Desktop or Docker Engine

```bash
cp .env.example .env
docker compose up --build
```

**Access:**

- GraphQL Playground: http://localhost:4000
- Gradio UI: http://localhost:4004

---

## Deployment

**Cloud Platform:** Render (Docker Web Services)

Each microservice runs in its own container with a shared PostgreSQL database. Services communicate via environment-configured URLs, making the system portable between local Docker Compose and cloud deployments.

**Key Architecture Difference:**

- **Local:** Services use Docker DNS (`http://gateway:4000`)
- **Cloud:** Services use full URLs (`https://gateway-render-e-commerce-graphql.onrender.com`)

Configuration managed entirely through environment variables for seamless local-to-production deployment.

---

## Why This Project

I built this system to demonstrate real backend engineering patterns:

- **Independent service deployment** - Each microservice can be updated without touching others
- **Federated GraphQL** - Type system spans services while maintaining domain boundaries
- **Explicit data ownership** - Each service owns its database models
- **Production patterns** - Cursor-based pagination, health checks, environment-based config
- **Schema-driven contracts** - GraphQL schemas define service APIs

The Gradio UI validates backend behavior without requiring a full frontend stack.

---

## Architecture

**Services:**

- **Products Service** (Go + GraphQL) - Product catalog, cursor-based pagination
- **Users Service** (Go + GraphQL) - User management
- **Orders Service** (Go + GraphQL) - Order lifecycle
- **Apollo Federation Gateway** (Node.js) - Schema composition, query routing
- **PostgreSQL** - Shared database (service-scoped models)
- **Gradio UI** (Python) - API validation interface

**Service Communication:**

- Federation via Apollo Gateway
- Each service exposes GraphQL schema with `@key` directives
- Gateway handles entity resolution across services

---

## Key Engineering Highlights

**GraphQL Federation:**

- Type extensions across service boundaries
- Entity reference resolution (e.g., Orders referencing Users)
- Schema composition via Apollo Gateway

**Pagination:**

- Cursor-based pagination for products (not offset-based)
- Configurable page sizes
- Forward-only navigation (production pattern)

**Infrastructure:**

- Docker health checks ensure startup ordering
- Environment-based service discovery (Docker Compose → Cloud)
- Each service independently scalable
- Shared database with clear model ownership

**Cloud Deployment:**

- Dockerized microservices on Render
- Dynamic port binding (`PORT` env var)
- Service-to-service communication via environment URLs
- Production debugging experience included

---

## Tech Stack

**Backend:**

- Go (gqlgen, GORM)
- GraphQL Federation (Apollo Gateway)
- PostgreSQL

**Infrastructure:**

- Docker & Docker Compose
- Render (cloud platform)

**Validation:**

- Python (Gradio)

---

## Current UI Functionality

**Working:**

- User creation via GraphQL mutation
- Product listing with cursor-based pagination
- Configurable page size
- Product detail view with federated data
- Cloud deployment on Render

**In Progress:**

- Order creation and management
- Authentication/authorization layer
- Comprehensive test coverage

The Gradio UI is intentionally lightweight - it exists to exercise the GraphQL API and validate federated queries, not to serve as a production frontend.

---

## Project Status

This is an **active learning project** focused on backend architecture and deployment patterns. The system is functional but evolving as I add authentication, authorization, and expand the order management flow.

**Full documentation:** See `README.dev.md` for detailed setup, environment variables, sample queries, and troubleshooting.

---

## Troubleshooting

**Port conflicts:**

```bash
docker ps -a  # Check running containers
```

**Reset everything:**

```bash
docker compose down -v
docker compose up --build
```

**Database connection errors:**
Health checks take ~30s. Wait for all services to show as healthy.
