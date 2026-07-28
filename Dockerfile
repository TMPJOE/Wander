# syntax=docker/dockerfile:1.7

# ============================================================
# Stage 1 — Build the frontend (Vue/Vite SPA)
# ============================================================
FROM node:22-alpine AS frontend-build

WORKDIR /build/frontend

# Cache deps separately so source changes don't re-install.
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && corepack prepare pnpm@latest --activate \
    && pnpm install --no-frozen-lockfile

# Build the SPA.
COPY frontend/ ./

# Build args
ARG VITE_API_BASE_URL
ARG VITE_APP_NAME
ARG VITE_APP_ENV
ARG VITE_STRIPE_PUBLISHABLE_KEY
ARG VITE_GOOGLE_MAPS_API_KEY

RUN pnpm build

# ============================================================
# Stage 2 — Build the Go backend
# ============================================================
FROM golang:1.25-alpine AS backend-build

WORKDIR /build/backend

# Cache module downloads.
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend source. Migrations/ is embedded via
# //go:embed *.sql, so the binary bakes them in at compile time — no
# need to carry migrations/ into the runtime image.
COPY backend/ ./

# Backend build args (for embedding into the binary).

ARG APP_NAME
ARG APP_ENV
ARG APP_PORT
ARG APP_HOST
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG DB_SSLMODE
ARG DB_SSLROOTCERT
ARG STORAGE_DRIVER
ARG STORAGE_LOCAL_DIR
ARG STORAGE_PUBLIC_BASE_URL
ARG S3_BUCKET
ARG S3_REGION
ARG S3_ENDPOINT
ARG S3_ACCESS_KEY
ARG S3_SECRET_KEY
ARG S3_FORCE_PATH_STYLE
ARG S3_PUBLIC_BASE_URL
ARG JWT_SECRET
ARG JWT_EXPIRATION_HOURS
ARG STRIPE_SECRET_KEY
ARG STRIPE_PUBLISHABLE_KEY
ARG ALLOWED_ORIGINS

# Build the server binary as a static-ish binary suitable for Alpine.
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/wander-server ./cmd/server

# Optionally build the seed binary (kept for ops, not needed at runtime).
RUN go build -ldflags="-s -w" -o /out/wander-seed ./cmd/seed || true

# ============================================================
# Stage 3 — Runtime image: just the binary + SPA + certs
# ============================================================
FROM alpine:3.20 AS runtime

# ca-certificates: required for TLS to Supabase/Stripe/S3.
# tzdata: sane time zone handling for booking date math.
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Frontend static assets — served by the Go server's SPA fallback.
COPY --from=frontend-build /build/frontend/dist ./frontend/dist

# The Go server binary.
COPY --from=backend-build /out/wander-server ./server

# Optional seed binary (not used by the start command; present for manual runs).
COPY --from=backend-build /out/wander-seed ./seed

# Volume mount point for local-mode uploads. In S3 storage mode this
# directory stays empty and is ignored.
RUN mkdir -p ./uploads
VOLUME ["/app/uploads"]

EXPOSE 8080

# Render (and most hosts) set PORT. Fall back to 8080 locally.
ENV APP_HOST=0.0.0.0
ENV APP_PORT=8080

CMD ["./server"]
