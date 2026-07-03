# Wander — Local Guide Hub: Full-Stack Implementation

Bring the Stitch "Local Guide Hub" project to life as a working full-stack application. Go backend with PostgreSQL (pgx driver), Vue 3 + Vite frontend. Media handled client-side (no S3/MinIO). Logging via `slog`.

## User Review Required

> [!IMPORTANT]
> The Stitch UI has 8 mobile screens. Per your direction, we are **not** strictly imitating the UI — the screens serve as an overview. We're building the full data model and the screens needed for guides to create/manage tours, plus all the interfaces implied by the UI but not shown.

> [!IMPORTANT]  
> **Auth strategy for demo**: We'll use JWT (already in config). Simple email+password registration/login. The `User` model gets a `role` field (`traveler` | `guide` | `admin`) to gate guide-specific flows.

> [!WARNING]
> **PostgreSQL must be running locally** before the backend can start. The existing `.env` has `wander_db` / `wander_user` / `wander_pass`. You'll need to create the DB + user manually, or I can add a setup script.

## Open Questions

> [!IMPORTANT]
> 1. **Seeding**: Do you want a seed script that populates the DB with demo categories, guides, and tours so the app has data on first run?
> 2. **Language**: The Stitch screens are in Spanish (Explorar, Reservas, etc.). Should the app stay in Spanish or switch to English?
> 3. **Messaging**: The "Mensajes" screen implies real-time chat. For a demo, should I do simple REST-based messaging (poll), or skip messaging entirely for now?

---

## Proposed Changes

### Database Layer

PostgreSQL schema with migrations using raw SQL files (no ORM). Uses `pgx/v5` as the driver with `pgxpool` for connection pooling.

#### [NEW] [000001_init_schema.up.sql](file:///c:/Users/josed/Documents/Projects/Wander/backend/migrations/000001_init_schema.up.sql)
#### [NEW] [000001_init_schema.down.sql](file:///c:/Users/josed/Documents/Projects/Wander/backend/migrations/000001_init_schema.down.sql)

Tables:
| Table | Purpose |
|---|---|
| `users` | Auth + profile. `role` enum: `traveler`, `guide`, `admin`. Avatar URL stored as string (frontend handles upload/display). |
| `categories` | Tour categories (Hiking, Cultural, Food, Adventure, etc.). Name, slug, icon name, description. |
| `tours` | Core entity. Belongs to a guide (user). Title, description, category_id, location, lat/lng, duration_minutes, price_per_person, max_guests, difficulty, language, what_included (JSONB), meeting_point, images (JSONB array of URLs — frontend-managed). |
| `tour_schedules` | Available dates/times for a tour. tour_id, start_time, end_time, available_spots. Guides create these to make tours bookable. |
| `bookings` | Links traveler → tour_schedule. status enum: `pending`, `confirmed`, `cancelled`, `completed`. guest_count, total_price, notes. |
| `reviews` | traveler → tour. rating (1-5), comment, created_at. |
| `messages` | Simple conversation model: sender_id, receiver_id, booking_id (optional), content, read_at. |
| `favorites` | user_id → tour_id. |

---

### Backend — Core Infrastructure

#### [MODIFY] [go.mod](file:///c:/Users/josed/Documents/Projects/Wander/backend/go.mod)
- Add `github.com/jackc/pgx/v5` + `github.com/golang-jwt/jwt/v5` + `golang.org/x/crypto` (bcrypt)
- Remove `github.com/go-chi/chi/v5` and `github.com/go-chi/cors` (we're using stdlib `net/http` + our own CORS middleware already)

#### [MODIFY] [config.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/config/config.go)
- Add `DatabaseURL()` method returning pgx-compatible connection string

#### [NEW] [database.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/config/database.go)
- `pgxpool.New()` initialization, connection health check, auto-run migrations from `migrations/` folder

#### [MODIFY] [main.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/cmd/server/main.go)
- Replace in-memory repo with pgx pool
- Wire all new services/handlers
- Use `slog` for structured startup logging
- Graceful shutdown with `os.Signal`

---

### Backend — Logging (slog)

#### [MODIFY] [logger.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/middleware/logger.go)
- Replace `log.Printf` with `slog.Info` structured logging (method, path, status, duration, remote_addr)

#### [MODIFY] [recovery.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/middleware/recovery.go)
- Replace `fmt.Printf` with `slog.Error` for panic recovery

---

### Backend — Models

#### [MODIFY] [user.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/user.go)
- Add `Role`, `Bio`, `Phone`, `AvatarURL`, `Languages` fields

#### [NEW] [category.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/category.go)
#### [NEW] [tour.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/tour.go)
- Tour + TourCreateRequest + TourUpdateRequest + TourFilter (for browsing with category, price range, difficulty, location)
#### [NEW] [tour_schedule.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/tour_schedule.go)
#### [NEW] [booking.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/booking.go)
#### [NEW] [review.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/review.go)
#### [NEW] [message.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/message.go)
#### [NEW] [favorite.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/models/favorite.go)

---

### Backend — Repositories (pgx)

All repos get a **PostgreSQL implementation** using `pgxpool.Pool`. The existing `UserRepository` interface stays, but the in-memory impl is replaced.

#### [MODIFY] [user_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/user_repo.go)
- Keep interface, replace `InMemoryUserRepository` with `PgUserRepository`

#### [NEW] [category_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/category_repo.go)
#### [NEW] [tour_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/tour_repo.go)
- List with filters (category, price range, difficulty, search query, location), pagination
- GetByGuideID for guide dashboard
#### [NEW] [tour_schedule_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/tour_schedule_repo.go)
#### [NEW] [booking_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/booking_repo.go)
#### [NEW] [review_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/review_repo.go)
#### [NEW] [message_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/message_repo.go)
#### [NEW] [favorite_repo.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/repository/favorite_repo.go)

---

### Backend — Services

#### [MODIFY] [user_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/user_service.go)
- Use bcrypt for password hashing
- Add Login method (returns JWT)

#### [NEW] [auth_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/auth_service.go)
- JWT generation/validation, token claims with user ID + role

#### [NEW] [category_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/category_service.go)
#### [NEW] [tour_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/tour_service.go)
- CRUD + filtering/search + ownership validation (only guide's own tours)
#### [NEW] [booking_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/booking_service.go)
- Create booking (check availability), cancel, confirm, list by user/guide
#### [NEW] [review_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/review_service.go)
#### [NEW] [message_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/message_service.go)
#### [NEW] [favorite_service.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/service/favorite_service.go)

---

### Backend — Handlers

#### [MODIFY] [handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/handler.go)
- Wire all new sub-handlers

#### [MODIFY] [user_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/user_handler.go)
- Add profile update, get current user (from JWT context)

#### [NEW] [auth_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/auth_handler.go)
- `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/auth/me`

#### [NEW] [category_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/category_handler.go)
- `GET /api/v1/categories`

#### [NEW] [tour_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/tour_handler.go)
- Public: `GET /api/v1/tours` (filter/search), `GET /api/v1/tours/{id}`
- Guide: `POST /api/v1/tours`, `PUT /api/v1/tours/{id}`, `DELETE /api/v1/tours/{id}`
- Guide: `GET /api/v1/guide/tours` (my tours)

#### [NEW] [schedule_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/schedule_handler.go)
- `GET /api/v1/tours/{tourId}/schedules`
- Guide: `POST /api/v1/tours/{tourId}/schedules`, `PUT .../{id}`, `DELETE .../{id}`

#### [NEW] [booking_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/booking_handler.go)
- `POST /api/v1/bookings`, `GET /api/v1/bookings` (my bookings), `PATCH /api/v1/bookings/{id}/cancel`
- Guide: `GET /api/v1/guide/bookings`, `PATCH /api/v1/guide/bookings/{id}/confirm`

#### [NEW] [review_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/review_handler.go)
- `POST /api/v1/tours/{tourId}/reviews`, `GET /api/v1/tours/{tourId}/reviews`

#### [NEW] [message_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/message_handler.go)
- `GET /api/v1/messages/conversations`, `GET /api/v1/messages/{userId}`, `POST /api/v1/messages/{userId}`

#### [NEW] [favorite_handler.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/handler/favorite_handler.go)
- `POST /api/v1/favorites/{tourId}`, `DELETE /api/v1/favorites/{tourId}`, `GET /api/v1/favorites`

---

### Backend — Auth Middleware

#### [MODIFY] [auth.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/internal/middleware/auth.go)
- Actual JWT validation, extract user ID + role into `context.Context`
- Add `RequireRole("guide")` middleware for guide-only routes

---

### Backend — Routes

#### [MODIFY] [routes.go](file:///c:/Users/josed/Documents/Projects/Wander/backend/api/routes.go)
- Reorganize into public routes, authenticated routes, and guide-only routes
- All routes under `/api/v1/` prefix

---

### Frontend — Vue 3 App

The frontend is a Vue 3 + Vite + Pinia app. We build the screens needed to interact with the backend, inspired by but not slavishly copying the Stitch designs.

#### [MODIFY] [App.vue](file:///c:/Users/josed/Documents/Projects/Wander/frontend/src/App.vue)
- Replace boilerplate with app shell: sticky bottom nav (Explorar, Reservas, Mensajes, Perfil), top header, RouterView

#### [MODIFY] [index.ts](file:///c:/Users/josed/Documents/Projects/Wander/frontend/src/router/index.ts)
- Full route table with auth guards

**New views (pages):**

| Route | View | Auth? | Description |
|---|---|---|---|
| `/` | `ExploreView` | No | Browse tours by category, search, filter |
| `/tours/:id` | `TourDetailView` | No | Tour details, reviews, guide info, book CTA |
| `/tours/:id/book` | `BookingView` | Yes | Select schedule, guests, confirm booking |
| `/booking-success/:id` | `BookingSuccessView` | Yes | Confirmation with details |
| `/bookings` | `MyBookingsView` | Yes | List user's bookings |
| `/messages` | `MessagesView` | Yes | Conversations list |
| `/messages/:userId` | `ChatView` | Yes | Individual conversation |
| `/profile` | `ProfileView` | Yes | User profile, edit, favorites |
| `/login` | `LoginView` | No | Login form |
| `/register` | `RegisterView` | No | Registration form |
| `/guide/dashboard` | `GuideDashboardView` | Guide | Stats, recent bookings |
| `/guide/tours` | `GuideToursView` | Guide | List/manage my tours |
| `/guide/tours/new` | `TourFormView` | Guide | Create a new tour |
| `/guide/tours/:id/edit` | `TourFormView` | Guide | Edit existing tour |
| `/guide/tours/:id/schedules` | `ScheduleManagerView` | Guide | Manage available dates/times |
| `/guide/bookings` | `GuideBookingsView` | Guide | View/confirm incoming bookings |
| `/category/:slug` | `CategoryView` | No | Filtered tour listing by category |

**New Pinia stores:**

#### [NEW] `stores/auth.ts` — JWT token management, current user state, login/logout/register actions
#### [NEW] `stores/tours.ts` — Tour listing, filters, current tour detail, CRUD for guides
#### [NEW] `stores/bookings.ts` — User bookings, guide bookings, create/cancel
#### [NEW] `stores/categories.ts` — Category list
#### [NEW] `stores/messages.ts` — Conversations, messages
#### [NEW] `stores/favorites.ts` — Favorited tours

**New composables:**

#### [NEW] `composables/useApi.ts` — Axios/fetch wrapper with JWT header injection, base URL from env
#### [NEW] `composables/useAuth.ts` — Auth guard logic, role check

**New components (reusable):**

| Component | Description |
|---|---|
| `BottomNav.vue` | Sticky bottom navigation with 4 tabs |
| `TourCard.vue` | Card for tour listings (image, title, price, rating, location) |
| `CategoryPill.vue` | Horizontal scrollable category filter pill |
| `ReviewCard.vue` | Individual review display |
| `BookingCard.vue` | Booking status card |
| `MessageBubble.vue` | Chat message bubble |
| `ConversationItem.vue` | Conversation list item |
| `StarRating.vue` | Interactive star rating |
| `ImageGallery.vue` | Tour image gallery (client-side images via URL) |
| `TourForm.vue` | Multi-step form for guides to create/edit tours |
| `ScheduleCalendar.vue` | Calendar for guides to add available dates |
| `FilterDrawer.vue` | Filter panel for tour search |
| `GuideCard.vue` | Guide mini-profile card |
| `EmptyState.vue` | Empty state illustration + message |

**Styles:**

#### [NEW] `assets/design-tokens.css` — CSS custom properties from the Stitch design system (terracotta/emerald palette, Inter typography, spacing, elevation)
#### [MODIFY] `assets/main.css` — Reset + global styles using design tokens

---

## Interface Mapping: UI Options → Backend Endpoints

These are features implied by the UI but not explicitly shown as screens:

| UI Element | Implied Backend Interface |
|---|---|
| Category pills on Explore | `GET /api/v1/categories` + filter param on tours |
| Search bar | `GET /api/v1/tours?q=...` full-text search |
| Heart/favorite icon on tour cards | `POST/DELETE /api/v1/favorites/{tourId}` |
| Star rating display | `GET /api/v1/tours/{id}/reviews` with avg calculation |
| "Book Now" button | `POST /api/v1/bookings` with schedule + guest count |
| Guide profile link on tour detail | `GET /api/v1/users/{id}` (public profile) |
| Guide dashboard stats | `GET /api/v1/guide/stats` (total bookings, revenue, avg rating) |
| Booking status changes | `PATCH /api/v1/guide/bookings/{id}/confirm` |
| Tour images in gallery | Stored as JSONB URLs in `tours.images` — frontend manages upload to any URL (e.g., data URLs, external links) |
| Price per person on cards | Part of `Tour` model, returned in list/detail endpoints |
| Tour difficulty badge | Enum in Tour model: `easy`, `moderate`, `challenging`, `extreme` |
| Language of tour | Field on Tour model |
| "What's included" list | JSONB array on Tour model |
| Meeting point | Text + optional lat/lng on Tour model |

---

## Verification Plan

### Automated Tests
```bash
# Backend — build check
cd backend && go build ./...

# Backend — run tests
cd backend && go test ./...

# Frontend — type check
cd frontend && pnpm type-check

# Frontend — lint
cd frontend && pnpm lint
```

### Manual Verification
1. Start PostgreSQL, create database with migration
2. `cd backend && go run ./cmd/server` — verify health check at `http://localhost:8080/health`
3. Register a user, login, get JWT
4. Create categories, create a tour as guide, add schedule
5. `cd frontend && pnpm dev` — browse the UI, verify tour listing, detail, booking flow
6. Test guide flow: create tour, add schedules, view incoming bookings
