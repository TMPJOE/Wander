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
