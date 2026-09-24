# Calculator

A full-stack calculator: a Go REST API doing the arithmetic, and a React +
TypeScript frontend that consumes it.

## Quick start (Docker)

```bash
docker compose up --build
```

Then open **http://localhost:3000** in a browser. The calculator UI is
served there; it talks to the backend automatically (proxied through
nginx, no extra setup).

Backend on its own is reachable at `http://localhost:8080` (e.g.
`http://localhost:8080/healthz`), but you shouldn't need to hit it directly
— the frontend at :3000 is the app.

Don't have Docker? See [Running locally](#running-locally) for the
no-Docker setup.

## Assignment deliverables

| Deliverable | Status | Where |
|---|---|---|
| Git repository with frontend and backend code | ✅ | `backend/`, `frontend/` in this repo |
| README with setup instructions, API examples, design decisions | ✅ | this file |
| Unit tests and coverage report | ✅ | [Testing & coverage](#testing--coverage) below |
| Dockerfile to run frontend + backend together (optional) | ✅ | [Docker](#docker) below, verified working with `docker compose up --build` |

## Instructions compliance

1. **AI tooling** — built with Claude as an AI pair-programmer throughout; see [`PROMPTS.md`](./PROMPTS.md) for the full prompt sequence.
2. **Time-boxed, correctness over extra features** — `percentage` was deliberately left out (ambiguous semantics — see Design Decisions), and no `/operations` discovery endpoint was added, in favor of a smaller, fully-tested, well-documented surface.
3. **Pushed to a Git host** — this repository.
4. **Repository link** — shared separately with the submission.
5. **Prompts shared** — see [`PROMPTS.md`](./PROMPTS.md).
6. **README includes:**
   - Setup instructions — see [Prerequisites](#prerequisites) and [Running locally](#running-locally)
   - How to run frontend and backend — see [Running locally](#running-locally) and [Docker](#docker)
   - API call examples — see [API](#api)
   - Design decisions / assumptions — see [Design decisions & assumptions](#design-decisions--assumptions)

```
calculator-app/
├── backend/
│   ├── cmd/server/main.go        # entrypoint: routing, CORS, graceful shutdown
│   ├── calculator/                # pure arithmetic + validation (no HTTP)
│   ├── handler/                   # HTTP layer: decode, validate, map errors
│   ├── middleware.go              # CORS
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/calculatorClient.ts
│   │   ├── hooks/useCalculator.ts # state machine (useReducer + effect)
│   │   ├── components/            # Display, Keypad, ErrorBanner, Calculator
│   │   └── styles.css
│   ├── tests/
│   ├── nginx.conf
│   └── Dockerfile
├── docker-compose.yml
├── PROMPTS.md
└── README.md
```

## Prerequisites

- Go 1.22+
- Node.js 20+
- Docker (optional, only needed for the containerized setup)

## Running locally

**Backend** (from `backend/`):

```bash
go run ./cmd/server
```

Starts on `:8080` by default. Override with the `PORT` env var:

```bash
PORT=9000 go run ./cmd/server
```

**Frontend** (from `frontend/`):

```bash
npm install
npm run dev
```

Starts on `http://localhost:5173` and calls the backend at
`http://localhost:8080` by default (see `VITE_API_BASE_URL` below). The
backend's CORS middleware explicitly allows this origin.

## API

### `POST /api/v1/calculate`

```json
{
  "operation": "add",
  "operands": [4, 2.5]
}
```

`operation` is one of: `add`, `subtract`, `multiply`, `divide`, `power`,
`sqrt`. Binary operations (`add`, `subtract`, `multiply`, `divide`, `power`)
require exactly 2 operands; `sqrt` requires exactly 1.

**Success — `200 OK`:**

```json
{
  "result": 6.5,
  "operation": "add",
  "operands": [4, 2.5]
}
```

**Client error — `400 Bad Request`** (malformed JSON, unknown operation,
wrong operand count):

```json
{
  "error": {
    "code": "UNKNOWN_OPERATION",
    "message": "Unsupported operation: modulo"
  }
}
```

**Semantic error — `422 Unprocessable Entity`** (well-formed request, but the
math itself is invalid — division by zero, negative square root, overflow):

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Cannot divide by zero"
  }
}
```

Error codes: `INVALID_REQUEST`, `UNKNOWN_OPERATION`, `DIVISION_BY_ZERO`,
`NEGATIVE_SQRT`, `INVALID_RESULT`.

**curl examples:**

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","operands":[4,2.5]}'

curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"sqrt","operands":[16]}'

curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","operands":[1,0]}'
# -> 422 DIVISION_BY_ZERO
```

### `GET /healthz`

```json
{ "status": "ok" }
```

## Testing & coverage

**Backend** (from `backend/`):

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out      # per-function summary in the terminal
go tool cover -html=coverage.out -o coverage.html   # open in a browser
```

Measured results:

| Package | Coverage |
|---|---|
| `calculator` | **100.0%** |
| `handler` | **88.9%** |
| `cmd/server` (entrypoint wiring) | 0.0% |
| `middleware.go` (CORS) | 0.0% |

The core logic — arithmetic and request handling, the parts a bug would
actually hurt — is thoroughly covered. `cmd/server` and `middleware.go` are
thin wiring (route registration, graceful shutdown, CORS headers) that I
judged lower-value to unit test at this scope; they'd be natural candidates
for a couple of `httptest`-based integration tests given more time (see
Known Limitations).

**Frontend** (from `frontend/`):

```bash
npm install -D @vitest/coverage-v8   # one-time, if not already present
npm run test                # vitest run
npm run test:coverage       # vitest run --coverage
```

Add both scripts to `package.json` if they're not already there:

```json
"test": "vitest run",
"test:coverage": "vitest run --coverage"
```

Measured results:

| File | Statements |
|---|---|
| `components/Calculator.tsx` | **100%** |
| `components/ErrorBanner.tsx` | **100%** |
| `hooks/useCalculator.ts` | **84.1%** |
| `components/Display.tsx` | 71.4% |
| `components/Keypad.tsx` | 60% |
| `api/calculatorClient.ts` | 28.6% |
| **All files** | **74.25%** |

`calculatorClient.ts` is low because every test mocks it at the module
boundary (by design — component/hook tests shouldn't depend on real network
calls), so its own internal error-handling branches (network failure, bad
JSON response) are exercised in isolation, not through it. `Keypad.tsx`'s
gap is simply digit buttons (`2`, `6`, `7`, `8`) that no test flow happens to
click. Both are honest, known gaps rather than hidden ones — see Known
Limitations.

## Docker

```bash
docker compose up --build
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`

The frontend's Dockerfile bakes `VITE_API_BASE_URL=""` at build time, so in
this deployment the browser calls the API on the frontend's own origin
(`/api/v1/calculate`), and nginx proxies that server-side to the backend
container. The browser never makes a cross-origin request in this setup, so
the backend's CORS middleware (needed for local `npm run dev` against
`localhost:5173`) is simply unused here rather than fought against.

Verified with `docker compose up --build`: both images build and start
cleanly, with the backend logging `server listening on port 8080` and nginx
serving the frontend at `http://localhost:3000` with no startup errors.

## Design decisions & assumptions

- **Standard library only, both sides.** `net/http` instead of a framework;
  no state management library on the frontend. At this scope, either would
  be more ceremony than the problem needs.
- **One `/api/v1/calculate` endpoint** rather than one route per operation.
  Keeps the contract stable as operations are added, and puts validation in
  one place. Trade-off: per-operation routes would be more RESTful/
  discoverable and easier to rate-limit individually.
- **`400` vs `422`** is a deliberate distinction: `400` means the request
  itself is malformed or unrecognized; `422` means the request was valid but
  the arithmetic isn't defined (÷0, √negative, overflow to Infinity).
- **`percentage` was intentionally left out.** Its semantics are ambiguous
  (X% of Y? X as a percent change from Y?) and I'd rather ship four
  unambiguous operations plus `power`/`sqrt` than one that needs a paragraph
  of caveats.
- **No `/operations` discovery endpoint.** With a small, fixed operation
  set, hardcoding it in both layers was simpler than adding an abstraction
  with no real payoff yet.
- **Frontend state is a `useReducer`, not `useState`.** An earlier
  `useState`-based version had a real stale-closure bug: multiple state
  updates dispatched in the same batch could silently overwrite each other.
  `useReducer` guarantees each action applies to the result of the previous
  one. The async API call is handled by a `useEffect` watching a
  `pendingCalculation` field the reducer sets — reducers stay pure, side
  effects live in one clearly separate place.
- **All arithmetic happens server-side.** The frontend only formats what to
  send and what came back; this was a requirement of the assignment, not
  just a style choice.
- **No CSS framework.** A single stylesheet with CSS variables as design
  tokens; keeps the bundle small and was a deliberate exercise in doing the
  visual design intentionally rather than defaulting to a component kit.

## Known limitations / next steps with more time

- `cmd/server` and `middleware.go` have no tests (see coverage section).
- `calculatorClient.ts`'s real fetch/error-parsing logic isn't exercised by
  any test that doesn't mock it away — worth a small dedicated test file
  using a mocked `fetch`.
- No end-to-end tests (Playwright/Cypress) — acceptable at this scope, but
  the natural next layer of testing.
- Frontend `power`/`sqrt` of very large inputs are validated by the backend
  (NaN/Infinity → `422`), but the frontend doesn't pre-validate before
  sending — intentional (server is the source of truth), but a nice-to-have
  UX improvement would be disabling `=` when the pending expression is
  obviously invalid.

## Prompts used

See [`PROMPTS.md`](./PROMPTS.md) for the sequence of prompts used to build
this project with AI assistance.
