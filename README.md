# Wander — Local Guide Hub

Wander is a full-stack marketplace platform that connects **travelers** with **local guides** for guided experiences, tours, and adventures. Guides publish tours with detailed itineraries and exact meeting points, travelers browse and book, and the platform handles scheduling, payments, reviews, favorites, and direct messaging between both sides.

The UI is in Spanish (`es`), tuned for the LATAM tourism market.

## Features

### For travelers
- Explore tours with search, category filters, and difficulty/price/location refinement
- Tour detail pages with itinerary, schedules, reviews, and an interactive meeting-point map with **"Abrir en Google Maps"** directions
- Booking flow powered by Stripe (test mode), with notes sent to the guide
- Booking detail with map, status tracking, and cancellation
- Favorites and one-rate-per-tour reviews
- Direct messaging with guides, streamed via SSE

### For guides
- Dashboard with tour and earnings stats
- Tour creation form with ordered itinerary editor and Google Maps **Place Picker** for exact meeting-point coordinates
- Schedule manager for availability windows and guest capacity
- Incoming booking queue showing traveler notes and a one-tap chat action
- Confirm / reject / complete booking lifecycle controls
- Profile with bio, languages, and avatar

### Platform
- Role-based auth (`traveler`, `guide`, `admin`) with JWT
- REST API with layered Go backend (handler → service → repository)
- PostgreSQL schema with versioned migrations and trigger-maintained aggregates (tour rating/count, `updated_at`)
- Full-text search index (Spanish) over tour title, description, and location
- Image uploads served from the backend, with the single binary also serving the built SPA

## Tech Stack

| Layer        | Technology                                                                 |
|--------------|----------------------------------------------------------------------------|
| Frontend     | Vue 3 (`<script setup>` + TS), Vue Router, Pinia, Vite                     |
| UI           | Lucide icons, design-token CSS, Google Maps JS API + Places Autocomplete   |
| Payments     | Stripe (server-side intents + client confirmation via `@stripe/stripe-js`) |
| Backend      | Go 1.25, `chi/v5` router, `pgx/v5`, `golang-jwt/v5`, `godotenv`            |
| Database     | PostgreSQL with versioned `.up.sql/.down.sql` migrations                   |
| Tests        | Vitest (frontend), Go testing (backend)                                   |

## Repository Structure

```
Wander/
├── backend/
│   ├── api/            # chi route setup + static/SPA serving
│   ├── cmd/            # server entrypoint
│   ├── internal/       # config, handler, service, repository, models, middleware
│   ├── migrations/     # versioned PostgreSQL migrations
│   ├── pkg/            # reusable helpers
│   ├── tests/
│   └── uploads/        # runtime image storage (gitignored)
├── frontend/
│   ├── src/
│   │   ├── components/ # TourForm, LocationMap, CategoryPill, FilterDrawer, ...
│   │   ├── views/      # Explore, TourDetail, Booking, Checkout, Guide*, ...
│   │   ├── composables/# useApi, useGoogleMaps, ...
│   │   └── assets/     # Wander logo, design-tokens.css
│   └── package.json
├── implementation_plan.md
├── setup_db.bat
└── .env.example
```

## Prerequisites

- **Go** 1.25+
- **Node.js** ^22.18.0 or >=24.12.0  (use **pnpm**)
- **PostgreSQL** (psql on PATH)
- A **Stripe** account (test keys) and a **Google Maps** API key with the Maps JavaScript API and Places API enabled

## Getting Started

### 1. Configure environment

Copy the example env files and fill in your secrets:

```sh
cp .env.example .env              # backend-facing combined config
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Required variables:

| Variable                      | Where              | Purpose                                  |
|-------------------------------|--------------------|-------------------------------------------|
| `DB_*`                        | `backend/.env`     | PostgreSQL connection                      |
| `JWT_SECRET`                  | `backend/.env`     | Token signing secret (change in prod!)    |
| `STRIPE_SECRET_KEY`           | `backend/.env`     | Stripe test secret key                     |
| `VITE_API_BASE_URL`           | `frontend/.env`    | API origin, e.g. `http://localhost:8080/api/v1` |
| `VITE_GOOGLE_MAPS_API_KEY`     | `frontend/.env`    | Maps JS API + Places key                   |
| `VITE_STRIPE_PUBLISHABLE_KEY`  | `frontend/.env`    | Stripe test publishable key                |

### 2. Create the database (Windows)

```cmd
setup_db.bat
```

This provisions the `wander_db` database and `wander_user` role, syncs the password into `db_setup.sql` and `backend/.env`, and runs the setup SQL via `psql`. On other platforms, run the equivalent statements from `backend/migrations/db_setup.sql` manually.

### 3. Apply migrations and run the backend

```sh
cd backend
go run ./cmd/...        # applies migrations on boot and serves :8080
```

The backend serves:
- API at `/api/v1/*`
- Uploaded images at `/uploads/*`
- The frontend production build at `/*` (SPA fallback) when `frontend/dist` exists

### 4. Run the frontend

```sh
cd frontend
pnpm install
pnpm dev        # Vite dev server (http://localhost:5173)
```

For a production build that the Go backend can serve:

```sh
cd frontend
pnpm build      # type-check + vite build -> frontend/dist
```

## Scripts

### Frontend (`frontend/package.json`)
```sh
pnpm dev          # dev server with HMR
pnpm build        # type-check + production build
pnpm type-check   # vue-tsc
pnpm test:unit    # vitest
pnpm lint         # oxlint + eslint --fix
pnpm format       # prettier
```

### Backend
```sh
go build ./...
go test ./...
```

## API Overview

All endpoints are versioned under `/api/v1`. Selected routes:

- `POST /auth/register`, `POST /auth/login`
- `GET /categories`, `GET /tours`, `GET /tours/{id}`
- `GET /tours/{id}/schedules`, `GET /tours/{id}/reviews`
- `POST /bookings`, `GET /bookings`, `GET /bookings/{id}`, `PATCH /bookings/{id}/cancel`
- `POST /payments/bookings/{id}/intent`, `POST /payments/bookings/{id}/confirm`
- `POST /tours/{tourId}/reviews`, `GET /reviews/me`
- `GET /favorites`, `POST /favorites/{tourId}`, `DELETE /favorites/{tourId}`
- `GET /messages/conversations`, `GET /messages/{userId}`, `POST /messages/{userId}`, `GET /messages/stream` (SSE)
- **Guide-only:** `POST/PUT/DELETE /tours`, `GET /guide/tours`, `GET /guide/bookings`, `GET /guide/earnings`, `PATCH /bookings/{id}/{confirm|reject|complete}`, `POST/PUT/DELETE /schedules`

See [`implementation_plan.md`](./implementation_plan.md) for the active roadmap, including the guided-tour itinerary and maps work.

## License

This project is currently private/not licensed for redistribution. Contact the maintainer for usage details.
