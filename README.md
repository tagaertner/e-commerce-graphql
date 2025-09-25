# E-Commerce Microservices Platform

A demonstration of a microservices architecture using **Go** and **Node.js**, built with **GraphQL Federation**, **PostgreSQL**, and **Docker**.

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Products       │     │  Users          │     │  Orders         │
│  Service        │     │  Service        │     │  Service        │
│  Port: 4001     │     │  Port: 4002     │     │  Port: 4003     │
│  (Go+GraphQL)   │     │  (Go+GraphQL)   │     │  (Go+GraphQL)   │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                  ┌──────────────┴──────────────┐
                  │     API Gateway             │
                  │     Port: 4000              │
                  │     (Node.js+Apollo Gateway)│
                  │     Unified GraphQL API     │
                  └─────────────┬───────────────┘
                                │
                  ┌─────────────┴───────────────┐
                  │     PostgreSQL Database     │
                  │     Port: 5432              │
                  │     (Shared across services)│
                  └─────────────────────────────┘
```

## Features

- **Microservices architecture** with separate services for products, users, and orders
- **Apollo Federation** GraphQL Gateway for unified API access
- **PostgreSQL database** with GORM for automatic migrations
- **Automated seeding** with test data
- **Docker** setup with health checks and proper service dependencies
- **Cross-service query** capability through federation

---

## 🚀 Quick Start

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose installed

### To Run Everything

1. **Clone the repository:**

   ```bash
   git clone https://github.com/tagaertner/e-commerce-graphql.git
   cd e-commerce-graphql
   ```

2. **Start all services:**

   ```bash
   docker-compose up --build
   ```

3. **Wait for services to initialize:**

   - Database will start and become healthy
   - Go services will connect and create tables via GORM
   - Seed data will be automatically inserted
   - Gateway will start and compose the federated schema

4. **Open GraphQL Playground:**
   👉 [http://localhost:4000](http://localhost:4000)

You can now run live GraphQL queries against the federated services.

---

## Sample Queries

### 🛍️ Get All Products

```graphql
query {
  products {
    id
    name
    description
    price
    inventory
  }
}
```

### 👤 Get User and Their Orders

```graphql need to fix to retreve orders
query {
  user(id: "1") {
    id
    name
    email
    orders {
      id
      productId
      quantity
      totalPrice
    }
  }
}
```

### 📦 Get Order by ID

```graphql
query {
  order(id: "3") {
    id
    userId
    productId
    quantity
    totalPrice
    status
  }
}
```

### 👥 Get Orders by User

```graphql
query {
  ordersByUser(userId: "4") {
    id
    quantity
    totalPrice
    status
  }
}
```

---

## Sample Mutations

### ➕ Create Order

```graphql
mutation {
  createOrder(
    input: {
      userId: "1"
      productId: "1"
      quantity: 2
      totalPrice: 3999.98
      status: "PENDING"
      createdAt: "2025-09-08T12:00:00Z"
    }
  ) {
    id
    userId
    quantity
    status
  }
}
```

### 🔁 Update Order

```graphql
mutation {
  updateOrder(input: { orderId: "3", quantity: 5, totalPrice: 1249.95, status: "CONFIRMED" }) {
    id
    quantity
    totalPrice
    status
  }
}
```

### ❌ Delete Order

```graphql
mutation {
  deleteOrder(input: { orderId: "4", userId: "4" })
}
```

---

## Service Endpoints

| Service  | Port | GraphQL Endpoint                                           | Container Health |
| -------- | ---- | ---------------------------------------------------------- | ---------------- |
| Products | 4001 | [http://localhost:4001/query](http://localhost:4001/query) | ✓ Health checks  |
| Users    | 4002 | [http://localhost:4002/query](http://localhost:4002/query) | ✓ Health checks  |
| Orders   | 4003 | [http://localhost:4003/query](http://localhost:4003/query) | ✓ Health checks  |
| Gateway  | 4000 | [http://localhost:4000](http://localhost:4000)             | Federated API    |
| Database | 5432 | PostgreSQL                                                 | ✓ Health checks  |

---

## Project Structure

```
e-commerce-graphql/
├── docker-compose.yml             # Service orchestration
├── .env                          # Environment variables (not tracked in git)
├── database/
│   └── init/
│       └── 01-seed-data.sql      # Test data for all services
├── gateway/                      # Node.js Apollo Federation Gateway
│   ├── gateway.js               # Gateway with federation composition
│   ├── package.json
│   └── dockerfile
├── services/                     # Go microservices
│   ├── orders/                  # Orders service
│   │   ├── main.go
│   │   ├── database/
│   │   ├── models/
│   │   ├── resolvers/
│   │   ├── services/
│   │   ├── schema.graphql
│   │   └── dockerfile
│   ├── products/                # Products service (same structure)
│   └── users/                   # Users service (same structure)
└── README.md
```

---

## Database Features

- **Automatic migrations:** GORM creates tables on service startup
- **Seed data:** Test users, products, and orders inserted automatically
- **Shared database:** All services connect to the same PostgreSQL instance
- **Health monitoring:** Database health checks ensure services start in correct order

---

## Sample Data Overview

- 👥 **Users:** 10 users (admins + customers, with active/inactive statuses)
- 📦 **Products:** 15 Apple ecosystem products
- 🧾 **Orders:** 20 orders with realistic statuses (pending, shipped, completed, cancelled)

---

## Development

### Environment Variables

Create a `.env` file:

```bash
# Database
POSTGRES_USER=ecom_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=ecom_db
DB_HOST=db
DB_PORT=5432

# Service Ports
PORT_PRODUCTS=4001
PORT_USERS=4002
PORT_ORDERS=4003
PORT_GATEWAY=4000
```

### Docker Compose Boot Order

1. PostgreSQL database starts first
2. Go services wait until DB is healthy and run GORM migrations
3. Seed data is inserted
4. Gateway starts and composes the federated schema

---

## Technical Stack

- **Backend Services:** Go + gqlgen (GraphQL API & resolvers)
- **Gateway:** Node.js + Apollo Federation
- **Database:** PostgreSQL (shared across services)
- **Containerization:** Docker & Docker Compose
- **Dev Experience:** Auto-migrations, seed data, health checks

---

## Future Enhancements

- **Automated Testing** — Add unit and integration tests for core services and GraphQL resolvers
- **Security Improvements** — Implement JWT-based authentication, input validation, and basic rate limiting

---

**MIT License Copyright (c) 2025 Tami Gaertner**
