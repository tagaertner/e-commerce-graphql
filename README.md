# E-Commerce GraphQL Microservices Platform

A production-style **microservices backend** built with **Go**, **GraphQL Federation**, and **PostgreSQL**, featuring a lightweight Python **Gradio UI** used for end-to-end API validation and exploration.

This project was designed to demonstrate **backend architecture**, **service composition**, and **stateful system design**, not frontend polish.

---

## 🚀 Quick Start

> Note: The Gradio UI is still evolving and is used to validate backend behavior rather than serve as a full client. It currently supports user creation, product browsing, and product detail exploration.

- GraphQL Playground: `http://localhost:4000`
- Gradio UI: `http://localhost:4004`

To run locally:

```bash
docker-compose up --build
```

## Why This Project

I built this system to go beyond CRUD demos and model how real backend teams design:

- independently deployable services
- federated GraphQL schemas
- explicit data ownership boundaries
- cursor-based pagination
- service-to-service contracts

The Gradio UI exists solely to validate backend behavior without needing a full frontend stack.

---

## Architecture Overview

- **Products Service** (Go + GraphQL)
- **Users Service** (Go + GraphQL)
- **Orders Service** (Go + GraphQL)
- **Apollo Federation Gateway** (Node.js)
- **PostgreSQL** (shared infrastructure, service-scoped models)
- **Gradio UI** (Python) for live querying & state inspection

---

## Key Engineering Highlights

- Cursor-based pagination for products
- Federated entity resolution across services
- Clear separation of domain logic and GraphQL resolvers
- Docker-first local development with health checks
- Schema-driven API design (GraphQL as contract)
- No frontend abstractions (API exercised directly)

---

## Tech Stack

- **Go** (gqlgen, GORM)
- **GraphQL Federation** (Apollo Gateway)
- **PostgreSQL**
- **Docker & Docker Compose**
- **Python (Gradio)**

---

## Full Documentation

For setup details, environment variables, and sample queries:  
➡️ `README.dev.md`

## Status

Actively evolving — authentication, authorization, and testing planned next.

The Gradio UI is intentionally lightweight and still evolving.

**Completed:**

- User creation
- Product listing with cursor-based pagination
- Configurable page size
- Product detail view

**In Progress:**

- Cart-style order creation
- Order history and order detail views

The UI exists to validate backend behavior and federated GraphQL flows, not to serve as a production frontend.
