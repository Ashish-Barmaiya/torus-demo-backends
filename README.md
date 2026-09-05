# Torus Demo Backends

This repository contains the lightweight Go backend services used to demonstrate [Torus](https://github.com/Ashish-Barmaiya/torus-proxy), a reverse proxy and load balancer. The backends generate realistic HTTP traffic, simulate latency and failure conditions, and expose controllable response sizes so the demo can exercise Torus routing, health checking, and load-balancing behavior without implementing a production application.

## What this repository demonstrates

This repository is only the backend portion of the larger Torus demo. The backends exist to create realistic traffic for a Torus deployment, not to provide a complete public SaaS service.

```text
Browser
   |
   | HTTPS / WebSocket
   v
Demo Controller
   |
   | HTTPS
   v
Torus
   |
   | HTTP
   v
Demo Backends
```

The current Compose topology includes four backend instances that are intentionally simple HTTP simulators:

- `users-a` and `users-b` serve the users service
- `orders-a` and `orders-b` serve the orders service

Each instance listens on its own internal port `3000` and is designed to demonstrate backend selection, latency, failures, health checks, and concurrent traffic patterns.

## Services

The repository exposes two demo services, each with two independent instances.

| Service | Instances | Port | Purpose |
| --- | --- | --- | --- |
| `users` | `users-a`, `users-b` | `3000` | Demo user data and user-related request flows |
| `orders` | `orders-a`, `orders-b` | `3000` | Demo order data and order-related request flows |

The service name is configured by the `SERVICE` environment variable and the instance name by `INSTANCE`. The default port is `3000` when `PORT` is unset.

## API

Each service exposes a small set of HTTP routes. The data is deterministic and in-memory. It is not persisted across requests.

### Users service

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Returns backend health status |
| `GET` | `/api/v1/users` | Returns all users in the static demo dataset |
| `GET` | `/api/v1/users/:id` | Returns one user by ID |
| `POST` | `/api/v1/users` | Simulates user creation |
| `PATCH` | `/api/v1/users/:id` | Simulates user update |
| `DELETE` | `/api/v1/users/:id` | Simulates user deletion |

### Orders service

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Returns backend health status |
| `GET` | `/api/v1/orders` | Returns all orders in the static demo dataset |
| `GET` | `/api/v1/orders/:id` | Returns one order by ID |
| `POST` | `/api/v1/orders` | Simulates order creation |
| `PATCH` | `/api/v1/orders/:id` | Simulates order update |
| `DELETE` | `/api/v1/orders/:id` | Simulates order deletion |

The repository returns static deterministic demo data. `GET` handlers read in-memory demo data and return it as JSON. `POST`, `PATCH`, and `DELETE` simulate operations for demo purposes; they do not persist any changes and should not be described as a production CRUD service.

### Example requests

```bash
curl http://localhost:3000/health
curl http://localhost:3000/api/v1/users/1
curl http://localhost:3000/api/v1/orders/1
```

## Response payload simulation

The backend supports response-size simulation through the internal header `X-Demo-Response-Size`.

This is not a public client-facing API. It is internal demo-control metadata intended for a controller-to-backend path in the larger architecture. The browser should not directly control these internals.

Supported values are:

- `0b`
- `1kb`
- `16kb`
- `64kb`
- `256kb`
- `1mb`
- `4mb`

When the header is absent, the backend returns its normal JSON response. When the header is present, the backend generates a JSON envelope containing the normal `data` and `meta` plus a filler `payload` field sized to the requested total response size. The maximum payload is intentionally limited to `4mb` as a demo constraint.

Example:

```bash
curl -H 'X-Demo-Response-Size: 1mb' \
  http://localhost:3000/api/v1/users
```

This response includes the regular user payload plus an additional filler payload to reach the target size.

## Request simulation

The backend supports request-scoped simulation headers:

- `X-Demo-Simulation`
- `X-Demo-Simulation-Delay`

Supported simulation modes:

- `normal`
- `slow`
- `error`

When the mode is omitted, the request behaves normally. When `slow` is used, the backend waits for a delay before continuing. With `error`, the request fails with `503 Service Unavailable`.

Delay behavior:

- default delay when omitted: `750ms`
- minimum delay: `250ms`
- maximum delay: `2000ms`

Example requests:

```bash
curl -H 'X-Demo-Simulation: normal' http://localhost:3000/api/v1/users
curl -H 'X-Demo-Simulation: slow' -H 'X-Demo-Simulation-Delay: 1000' \
  http://localhost:3000/api/v1/users
curl -H 'X-Demo-Simulation: error' http://localhost:3000/api/v1/users
```

Request simulation is intentionally independent per request and does not mutate backend-instance state. There is no global simulation state shared across all requests.

```text
users-a instance
  ├─ request A: X-Demo-Simulation=slow, delay=1000ms
  ├─ request B: X-Demo-Simulation=normal
  └─ request C: X-Demo-Simulation=error
```

These requests can all hit the same backend instance concurrently, and each request can have a different simulation mode without affecting the others.

## Health checks

The `/health` endpoint is special and does not participate in request simulation.

This is intentional and important: Torus actively calls `/health` to maintain its healthy backend pool. If a slow or error simulation were allowed to delay or fail `/health`, Torus could incorrectly mark otherwise healthy backends as unhealthy. The demo must simulate application request behavior without corrupting the health-check control plane.

Therefore:

- `/health` remains fast and healthy even when regular application requests are slow or erroring
- error simulation does not affect `/health`
- request simulation is limited to regular API traffic, not the readiness/health path used by Torus

This separation ensures that Torus can still evaluate backend health independently of application-level request simulation.

## Local development

### Prerequisites

- Go toolchain installed on the machine
- `SERVICE` and `INSTANCE` must be set when running the backend process
- `PORT` defaults to `3000` if omitted

### Run the test suite

```bash
go test ./...
```

```bash
go test -race ./...
```

### Run a single backend instance

```bash
SERVICE=users INSTANCE=users-a PORT=3000 go run ./cmd/backend
```

This starts one backend process, serving the `users` service as instance `users-a` on port `3000`.

### Example curl

```bash
curl http://localhost:3000/health
```

## Docker

The repository includes a `Dockerfile` that builds a multi-stage Go image:

- builder stage: `golang:1.26`
- runtime stage: `gcr.io/distroless/static-debian12`

The runtime image is intentionally minimal and has no shell. This keeps the container image small and avoids exposing interactive shell access in the runtime environment.

Build the image:

```bash
docker build -t torus-demo-backend:dev .
```

Example local run:

```bash
docker run --rm \
  -e SERVICE=users \
  -e INSTANCE=users-a \
  -e PORT=3000 \
  -p 3000:3000 \
  torus-demo-backend:dev
```

## Docker Compose

The `compose.yml` file defines four backend services:

- `users-a`
- `users-b`
- `orders-a`
- `orders-b`

Each service runs the same image with different `SERVICE`, `INSTANCE`, and `PORT` environment values. No host ports are exposed by the Compose topology, which is intentional because the Torus deployment is expected to sit in front of these internal services.

Start the demo backend topology:

```bash
docker compose up -d
```

Inspect running containers:

```bash
docker compose ps
```

Current topology mapping:

```text
users-a -> SERVICE=users INSTANCE=users-a PORT=3000
users-b -> SERVICE=users INSTANCE=users-b PORT=3000
orders-a -> SERVICE=orders INSTANCE=orders-a PORT=3000
orders-b -> SERVICE=orders INSTANCE=orders-b PORT=3000
```

## Project structure

```text
cmd/
  backend/
    main.go
internal/
  api/
  config/
  data/
  payload/
  simulation/
Dockerfile
compose.yml
go.mod
backend
```

Package responsibilities:

- `cmd/backend`: entrypoint for the backend binary
- `internal/api`: HTTP handlers, validation, and JSON response generation
- `internal/config`: environment-based configuration loading for `SERVICE`, `INSTANCE`, and `PORT`
- `internal/data`: deterministic in-memory demo data and response models
- `internal/payload`: response-size generation and payload envelope logic
- `internal/simulation`: request-scoped latency and failure simulation middleware
- `Dockerfile`: multi-stage build for the demo backend image
- `compose.yml`: Compose topology for the four backend instances

## Architecture principles

These are the main design rules of this repository:

- This is a demo backend, not a production SaaS application.
- Demo data is deterministic and static in memory.
- There is no persistence layer or database.
- Simulation is request-scoped and does not mutate backend instance state.
- `/health` is isolated from simulation so Torus can maintain a correct health model.
- Response-size simulation is done through internal demo metadata headers.
- Multiple backend instances are intentionally independent and run in parallel.
- Torus remains responsible for routing, health checking, and load balancing.

## License

This project is licensed under the [MIT License](./LICENSE).
