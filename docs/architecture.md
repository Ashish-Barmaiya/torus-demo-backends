# Architecture

## Purpose

This repository exists to provide controlled backend behavior for demonstrating Torus. The four Go services generate deterministic HTTP traffic, support request-scoped failure and delay simulation, and expose response-size control so a Torus deployment can demonstrate service routing, load balancing, backend selection, health status, and traffic behavior. The code intentionally stays narrow: it simulates a backend environment, not an application platform.

## Runtime topology

```text
Torus
  ├── users
  │   ├── users-a
  │   └── users-b
  └── orders
      ├── orders-a
      └── orders-b
```

Each backend instance is an independent process or container. They are configured separately with `SERVICE`, `INSTANCE`, and `PORT`, and they do not share mutable application state. Torus sits in front of them and is expected to route traffic across the four instances.

## Request lifecycle

```text
HTTP request
  -> simulation middleware
  -> API handler
  -> data
  -> response-size generation
  -> HTTP response
```

The actual behavior is implemented as follows:

- `internal/simulation/middleware.go` reads `X-Demo-Simulation` and `X-Demo-Simulation-Delay` and decides whether to delay, fail, or pass the request through.
- `internal/api/server.go` routes requests to the correct service and HTTP method handlers for `/api/v1/users`, `/api/v1/orders`, and the per-resource endpoints.
- `internal/data/*.go` provides the deterministic in-memory data for users and orders.
- `internal/api/server.go` calls `payload.GenerateJSONWithData` when `X-Demo-Response-Size` is present.
- `internal/payload/*.go` creates a valid JSON envelope with a filler payload sized to the requested target size.

This is a request-by-request pipeline: simulation behavior occurs at the HTTP boundary, then the request is handled by the service, and then the response may be expanded to the target payload size.

## Request-scoped simulation

Simulation state is intentionally not global. The backend does not maintain a mutable simulation mode in shared memory or instance-level state. Instead, each request carries its own simulated behavior through request headers.

This design avoids hidden coupling between requests. Multiple concurrent requests may behave differently even when they hit the same backend instance:

```text
users-a receives three concurrent requests:
  - request 1: X-Demo-Simulation=normal
  - request 2: X-Demo-Simulation=slow, delay=1000ms
  - request 3: X-Demo-Simulation=error
```

Each request is evaluated independently. One request's simulation mode does not mutate another request's state or change the backend instance's long-lived behavior. This is a deliberate design decision that makes concurrent traffic patterns predictable and physically realistic for the demo.

## Health-check isolation

The `/health` endpoint is excluded from simulation in `internal/simulation/middleware.go`.

This isolation matters because Torus uses `/health` to maintain its healthy backend pool. If a slow or error simulation were allowed to affect `/health`, Torus could incorrectly decide that an otherwise healthy instance was unhealthy. The demo must simulate application traffic without corrupting the health-check control plane.

The contract is therefore straightforward:

- regular application requests may be slow or fail based on simulation headers
- `/health` remains fast and healthy regardless of request-scoped simulation
- error simulation does not affect `/health`

This keeps the demo backend semantics aligned with Torus's health model: the control plane for backend health remains independent from the traffic simulation plane.

## Payload generation

The payload package is generic rather than coupled to the `users` or `orders` models. It creates valid JSON envelopes with a `data` field, a `meta` field, and a filler `payload` string so that the total response size can approximate a target byte count.

Supported target sizes are:

- `0b`
- `1kb`
- `16kb`
- `64kb`
- `256kb`
- `1mb`
- `4mb`

The generator uses a repeated filler string and iteratively refines the payload length so the final JSON is as close as possible to the requested total size. The implementation is designed for demo traffic, and the `4mb` cap is an explicit boundary rather than a production payload limit.

The feature is intentionally generic: it can wrap any JSON object that implements the response data/meta interfaces, without being tied to a specific service model.

## Data model

The static dataset is generated deterministically in `internal/data/data.go`.

Current dataset sizes:

- 20 users
- 20 orders

User IDs are generated as `usr_000001` through `usr_000020`.
Order IDs are generated as `ord_000001` through `ord_000020`.

These values are generated in memory each time the data functions run. `POST`, `PATCH`, and `DELETE` handlers do not persist changes to a backing database or any in-memory store beyond the request life cycle. The API simulates mutations for demo purposes only; it does not create durable application state.

## Internal control metadata

Headers such as:

- `X-Demo-Simulation`
- `X-Demo-Simulation-Delay`
- `X-Demo-Response-Size`

exist as internal control metadata for the demo and are not intended to represent a public API contract. They are part of the Demo Controller -> backend path, where a controller validates and injects the appropriate headers before they reach Torus and the backend.

This distinction matters because the browser should not directly control backend internals. The public architecture is designed for safe scenario selection in the controller, followed by validated backend requests routed through Torus.

## Deliberate non-goals

This repository explicitly does not provide the following:

- authentication or authorization
- persistence or database storage
- queues or distributed transactions
- service discovery or dynamic registry behavior
- real business workflows
- production-grade user or order management
- public-facing API security controls
- WebSockets
- the Demo Controller
- the frontend
- Torus configuration management

These are intentionally outside the repository's purpose and are not implemented here.

## Design constraints

The implementation follows these important invariants:

1. Simulation is request-scoped.
2. `/health` is not simulated.
3. Backend instances remain independent.
4. Demo data is deterministic.
5. Mutations are simulated, not persisted.
6. Public users do not directly control backend internals.
7. Torus remains responsible for routing and health checking.

These constraints are visible in the code and help keep the backend aligned with its role as a demo traffic generator rather than a production service.
