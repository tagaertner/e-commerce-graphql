# E-Commerce Microservices Platform

A demonstration of microservices architecture using **Go** and **Node.js**, and GraphQL APIs with Docker support.

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
                  │     (Node.js+Express)       │
                  │     Unified GraphQL API     │
                  └─────────────────────────────┘
```

## Features

- **Microservices Architecture**: Three independent Go services
- **API Gateway**: Node.js proxy layer to unify GraphQL access
- **Cross-Service Queries**: Fetch data across services
- **Health Monitoring**: Each service exposes a `/health` endpoint
- **Dockerized**: One-line boot-up for all services

---

## 🚀 Quick Start (Docker)

### Prerequisites

- [Docker](https://www.docker.com/) installed

### Run Everything

```bash
git clone https://github.com/tagaertner/e-commerce-graphql.git
cd e-commerce-graphql
docker compose up --build

🧠 **You won’t see terminal logs until you run queries.**

Then, open the GraphQL Playground:

👉 [http://localhost:4000/graphql](http://localhost:4000/graphql)

You can now run live GraphQL queries against the composed services.


### Example Queries
```

**Get products:**

```graphql
query {
  es
  products {
    id
    name
    price
    inventory
  }
}
```

**Get user orders:**

```graphql
query {
  user(id: "1") {
    name
    email
  }
  ordersByUser(userId: "1") {
    id
    totalPrice
    status
  }
}
```

````
**Cross-service data:**

```graphql
query {
  products {
    id
    name
    price
  }
  users {
    id
    name
    email
  }
  orders {
    id
    userId
    productId
    totalPrice
  }
}
````

## Service Endpoints

| Service  | Port | GraphQL Playground            | Health Check                 |
| -------- | ---- | ----------------------------- | ---------------------------- |
| Products | 4001 | http://localhost:4001/        | http://localhost:4001/health |
| Users    | 4002 | http://localhost:4002/        | http://localhost:4002/health |
| Orders   | 4003 | http://localhost:4003/        | http://localhost:4003/health |
| Gateway  | 4000 | http://localhost:4000/graphql | http://localhost:4000/health |

## Project Structure

```
e-commerce-graphql/
├── docker-compose.yml         # Docker orchestration
├── gateway/                   # Node.js API Gateway
│   ├── gateway.js            # Service composition logic
│   ├── package.json
│   └── Dockerfile            # Gateway container config
├── services/                 # Go microservices
│   ├── products/
│   │   ├── main.go          # HTTP server setup
│   │   ├── schema.graphql   # GraphQL schema definition
│   │   ├── gqlgen.yml       # Code generation config
│   │   ├── generated/       # Auto-generated GraphQL code
│   │   ├── models/          # Data models
│   │   ├── resolvers/       # Business logic
│   │   └── Dockerfile       # Products service container
│   ├── users/               # Same structure
│   └── orders/              # Same structure
└── README.md
```

## Docker Configuration

The project includes Docker support with:

- **Multi-stage builds** for optimized Go service images
- **Service networking** for inter-service communication
- **Health checks** for container monitoring
- **Volume mounting** for development (optional)

### Docker Services

- **gateway**: Node.js API Gateway (port 4000)
- **products**: Go Products service (port 4001)
- **users**: Go Users service (port 4002)
- **orders**: Go Orders service (port 4003)

## Future Development

**1. Testing**

**2. PostgreSQL Database Integration**

**3. Advanced Queries**

**4. Authentication & Authorization**

**5. AWS/GCP Deployment**

## Technical Details

**Backend Services:** Go with gqlgen for GraphQL server generation  
**API Gateway:** Node.js with Express and node-fetch for service composition  
**Communication:** HTTP/GraphQL between gateway and services  
**Containerization:** Docker with multi-stage builds for production-ready images

## Sample Data

The services include sample data for demonstration:

- 3 products (laptop, smartphone, headphones)
- 3 users (customers and admin)
- 3 orders with different statuses

---

**MIT License Copyright (c) 2025 Tami Gaertner**
