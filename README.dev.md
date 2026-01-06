# E-Commerce Microservices Platform

A production-style **e-commerce backend** built with **Go microservices**, **Apollo GraphQL Federation**, **PostgreSQL**, and **Docker**, with a lightweight **Gradio UI** used to validate and explore the API end-to-end.

This project focuses on **backend architecture**, **service boundaries**, and **real-world GraphQL patterns**, rather than frontend polish.

---

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Products       │     │  Users          │     │  Orders         │
│  Service        │     │  Service        │     │  Service        │
│  Port: 4001     │     │  Port: 4002     │     │  Port: 4003     │
│  (Go + GraphQL) │     │  (Go + GraphQL) │     │  (Go + GraphQL) │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                  ┌──────────────┴──────────────┐
                  │     API Gateway              │
                  │     Port: 4000               │
                  │     (Apollo Federation)      │
                  │     Unified GraphQL API      │
                  └─────────────┬───────────────┘
                                │
                  ┌─────────────┴───────────────┐
                  │     PostgreSQL Database      │
                  │     Port: 5432               │
                  │     Shared Across Services   │
                  └─────────────────────────────┘
```

---

## Features

- Microservices architecture with separate **Products**, **Users**, and **Orders** services
- **Apollo Federation Gateway** composing a unified GraphQL schema
- **Cursor-based pagination** for product listings
- **Cross-service queries** via GraphQL federation
- **PostgreSQL** with GORM auto-migrations
- **Automated seed data**
- **Docker Compose** orchestration with health checks
- **Gradio UI** for API validation and exploration

---

## Tech Stack

- **Backend Services:** Go, gqlgen, GORM
- **GraphQL:** Apollo Federation, Apollo Gateway
- **Database:** PostgreSQL
- **UI:** Python, Gradio
- **Infrastructure:** Docker, Docker Compose

---

## 🚀 Quick Start

### Prerequisites

- Docker
- Docker Compose

### Run the Platform

```bash
git clone https://github.com/tagaertner/e-commerce-graphql.git
cd e-commerce-graphql
docker-compose up --build
```

### Available Services

| Service   | Port | URL                         |
| --------- | ---- | --------------------------- |
| Gateway   | 4000 | http://localhost:4000       |
| Products  | 4001 | http://localhost:4001/query |
| Users     | 4002 | http://localhost:4002/query |
| Orders    | 4003 | http://localhost:4003/query |
| Gradio UI | 4004 | http://localhost:4004       |
| Database  | 5432 | PostgreSQL                  |

---

## Sample Queries

### Get Products (Paginated)

```graphql
query {
  productsCursor(first: 10) {
    edges {
      node {
        id
        name
        price
        inventory
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Get Product by ID

```graphql
query {
  product(id: "1") {
    id
    name
    description
    price
    inventory
    available
  }
}
```

### Get Orders for a User

```graphql
query {
  ordersByUser(userId: "1") {
    id
    quantity
    totalPrice
    status
    products {
      id
      name
    }
  }
}
```

---

## Sample Mutations

### Create User

```graphql
mutation {
  createUser(
    input: { name: "Jane Doe", email: "jane@example.com", password: "password123", role: CUSTOMER, active: true }
  ) {
    id
    name
    email
  }
}
```

### Create Order

```graphql
mutation {
  createOrder(
    input: {
      userId: "1"
      productIds: ["1", "2"]
      quantity: 2
      totalPrice: 3999.98
      status: PENDING
      createdAt: "2025-01-01T12:00:00Z"
    }
  ) {
    id
    status
  }
}
```

---

## Project Structure

```
e-commerce-graphql/
├── docker-compose.yml
├── .env
├── database/
│   └── init/
│       └── 01-seed-data.sql
├── gateway/
│   ├── gateway.js
│   ├── package.json
│   └── dockerfile
├── services/
│   ├── products/
│   ├── users/
│   └── orders/
├── gradio_ui/
│   ├── app.py
│   ├── interface.py
│   ├── handlers.py
│   ├── graphql_client.py
│   └── dockerfile
├── README.md
└── README.dev.md
```

---

## Database Behavior

- GORM auto-migrations on startup
- Shared PostgreSQL database across services
- Seed data automatically inserted
- Health checks ensure correct startup order

---

## Sample Data

- **Users:** customers and admins
- **Products:** Apple ecosystem catalog
- **Orders:** realistic order states (PENDING, SHIPPED, COMPLETED)

---

## Development Notes

### Environment Variables

```bash
POSTGRES_USER=ecom_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=ecom_db
DB_HOST=db
DB_PORT=5432

PORT_PRODUCTS=4001
PORT_USERS=4002
PORT_ORDERS=4003
PORT_GATEWAY=4000
PORT_GRADIO=4004
```

---

## Why This Project Matters

This project demonstrates:

- Backend ownership across multiple services
- Real GraphQL Federation (not mocked schemas)
- Cursor pagination and data flow
- Dockerized local development workflows
- Debugging distributed systems and API contracts

---

## Future Enhancements

- JWT authentication and role-based access control
- Integration and unit testing
- Observability and structured logging
- Rate limiting and validation
- Frontend client consuming the same API

---

## License

MIT License  
© 2025 Tami Gaertner
