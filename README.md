# E-Commerce GraphQL Microservices Platform

A production-style microservices backend built with Go, GraphQL Federation, and PostgreSQL, featuring a lightweight Python Gradio UI for end-to-end API validation and exploration. Deployed on Render with Neon PostgreSQL.
This project focuses on backend architecture, service composition, and real-world deployment behavior across distributed services.

## Demo & Availability

This project is hosted on Render’s free tier. Because the services can go idle, the first request may fail or take time while the services restart.

## Live Demo

- **Gradio UI:** https://gradio-ui-render-e-commerce-graphql.onrender.com
- **GraphQL Gateway:** https://gateway-render-e-commerce-graphql.onrender.com

### Warm Up Render Services

If the live demo does not respond right away, warm up the Render services first.

```bash
curl -X POST https://gateway-render-e-commerce-graphql.onrender.com/ \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __typename }"}'
curl -X POST https://order-render-e-commerce-graphql.onrender.com/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __typename }"}'
curl -X POST https://products-render-ecommercegraphql.onrender.com/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __typename }"}'
curl -X POST https://user-render-e-commercegraphql.onrender.com/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __typename }"}'
```

Once the services are awake, refresh the Gradio UI or retry the GraphQL Gateway.

⸻

Video Walkthroughs (with audio)

E-Commerce UI Demo

https://github.com/user-attachments/assets/8176fa10-a4d4-4653-84bb-032f2cb8c933

GraphQL Demo

https://github.com/user-attachments/assets/a79566af-e485-46d7-9b76-c01a3b7c6ba2

⸻

Quick Start

Local Development

Prerequisites: Docker Desktop or Docker Engine

cp .env.example .env
docker compose up --build

Access

- GraphQL Playground: http://localhost:4000
- Gradio UI: http://localhost:4004

⸻

Deployment

Cloud Platform: Render (Docker Web Services)

Each microservice runs in its own container with a shared PostgreSQL database. Services communicate through environment-configured URLs, making the system portable between local Docker Compose and cloud deployments.

Key Architecture Difference

- Local: Services communicate through Docker DNS (http://gateway:4000)
- Cloud: Services communicate through deployed URLs (https://gateway-render-e-commerce-graphql.onrender.com)

Configuration is managed entirely through environment variables for seamless local-to-cloud deployment.

⸻

Why This Project

I built this system to better understand how distributed backend systems behave outside of local development environments.

The project focuses on:

- Independent service deployment
- Federated GraphQL architecture
- Explicit service ownership
- Environment-based configuration
- Real-world deployment and debugging experience across distributed services

The Gradio UI validates backend behavior without requiring a full frontend stack.

⸻

Architecture

Services

- Products Service (Go + GraphQL) - Product catalog with cursor-based pagination
- Users Service (Go + GraphQL) - User management
- Orders Service (Go + GraphQL) - Order lifecycle management
- Apollo Federation Gateway (Node.js) - Schema composition and query routing
- PostgreSQL - Shared database with service-scoped models
- Gradio UI (Python) - Lightweight API validation interface

Service Communication

- Federation handled through Apollo Gateway
- Each service exposes its own GraphQL schema with @key directives
- Gateway resolves entities across services

⸻

Key Engineering Highlights

GraphQL Federation

- Type extensions across service boundaries
- Entity reference resolution between services
- Schema composition through Apollo Gateway

Pagination

- Cursor-based pagination for products
- Configurable page sizes
- Forward-only navigation pattern

Infrastructure

- Docker health checks for startup ordering
- Environment-based service discovery
- Independently deployable services
- Shared PostgreSQL database with clear model ownership

Cloud Deployment

- Dockerized microservices deployed on Render
- Dynamic port binding using environment variables
- Service-to-service communication through deployed URLs
- Real-world debugging and deployment constraints

⸻

Tech Stack

Backend

- Go (gqlgen, GORM)
- GraphQL Federation (Apollo Gateway)
- PostgreSQL

Infrastructure

- Docker & Docker Compose
- Render
- Neon PostgreSQL

Validation

- Python (Gradio)

⸻

Current UI Functionality

Working

- User creation through GraphQL mutations
- Product listing with cursor-based pagination
- Configurable page sizes
- Product detail view with federated data
- Cloud deployment on Render

In Progress

- Order creation and management
- Authentication and authorization
- Expanded automated test coverage

The Gradio UI is intentionally lightweight. It exists to exercise the GraphQL API and validate federated queries rather than serve as a production frontend.

⸻

Project Status

This is an active learning project focused on backend architecture, GraphQL Federation, and distributed system deployment patterns.

The system is functional but still evolving as I expand authentication, authorization, and order management workflows.

Additional documentation: See README.dev.md for environment variables, troubleshooting, sample GraphQL queries, and local development details.

⸻

Troubleshooting

Port Conflicts

docker ps -a

Reset Everything

docker compose down -v
docker compose up --build

Database Connection Delays

Docker health checks can take ~30 seconds during startup. Wait for all services to report healthy before testing the system.
