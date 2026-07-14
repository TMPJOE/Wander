package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type PgTourRepository struct {
	pool *pgxpool.Pool
}

func NewPgTourRepository(pool *pgxpool.Pool) TourRepository {
	return &PgTourRepository{pool: pool}
}

func (r *PgTourRepository) Create(ctx context.Context, guideID int, req models.TourCreateRequest) (*models.Tour, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO tours (guide_id, category_id, title, description, location, latitude, longitude, duration_minutes, price_per_person, max_guests, difficulty, languages, what_included, meeting_point, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, guide_id, category_id, title, description, location, latitude, longitude, duration_minutes, price_per_person, max_guests, difficulty, languages, what_included, meeting_point, is_published, avg_rating, review_count, created_at, updated_at
	`
	t := &models.Tour{}
	err = tx.QueryRow(ctx, query,
		guideID, req.CategoryID, req.Title, req.Description, req.Location, req.Latitude, req.Longitude,
		req.DurationMinutes, req.PricePerPerson, req.MaxGuests, req.Difficulty, req.Languages,
		req.WhatIncluded, req.MeetingPoint, req.IsPublished,
	).Scan(
		&t.ID, &t.GuideID, &t.CategoryID, &t.Title, &t.Description, &t.Location, &t.Latitude, &t.Longitude,
		&t.DurationMinutes, &t.PricePerPerson, &t.MaxGuests, &t.Difficulty, &t.Languages, &t.WhatIncluded,
		&t.MeetingPoint, &t.IsPublished, &t.AvgRating, &t.ReviewCount, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert tour: %w", err)
	}

	// Insert images
	for i, imgURL := range req.Images {
		_, err = tx.Exec(ctx, `
			INSERT INTO tour_images (tour_id, url, position)
			VALUES ($1, $2, $3)
		`, t.ID, imgURL, i+1)
		if err != nil {
			return nil, fmt.Errorf("insert tour image %d: %w", i, err)
		}
	}
	t.Images = req.Images

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create tx: %w", err)
	}

	return t, nil
}

func (r *PgTourRepository) GetByID(ctx context.Context, id int, userID int) (*models.Tour, error) {
	query := `
		SELECT t.id, t.guide_id, t.category_id, t.title, t.description, t.location, t.latitude, t.longitude,
		       t.duration_minutes, t.price_per_person, t.max_guests, t.difficulty, t.languages, t.what_included,
		       t.meeting_point, t.is_published, t.avg_rating, t.review_count, t.created_at, t.updated_at,
		       u.first_name || ' ' || u.last_name as guide_name, u.avatar_url as guide_avatar,
		       c.name as category_name, c.slug as category_slug,
		       EXISTS(SELECT 1 FROM favorites f WHERE f.tour_id = t.id AND f.user_id = $2) as is_favorited
		FROM tours t
		JOIN users u ON t.guide_id = u.id
		JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1
	`
	t := &models.Tour{}
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.GuideID, &t.CategoryID, &t.Title, &t.Description, &t.Location, &t.Latitude, &t.Longitude,
		&t.DurationMinutes, &t.PricePerPerson, &t.MaxGuests, &t.Difficulty, &t.Languages, &t.WhatIncluded,
		&t.MeetingPoint, &t.IsPublished, &t.AvgRating, &t.ReviewCount, &t.CreatedAt, &t.UpdatedAt,
		&t.GuideName, &t.GuideAvatar, &t.CategoryName, &t.CategorySlug, &t.IsFavorited,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}

	// Fetch images
	rows, err := r.pool.Query(ctx, `
		SELECT url FROM tour_images WHERE tour_id = $1 ORDER BY position ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("fetch tour images: %w", err)
	}
	defer rows.Close()

	t.Images = []string{}
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			t.Images = append(t.Images, url)
		}
	}

	return t, nil
}

func (r *PgTourRepository) Update(ctx context.Context, id int, req models.TourUpdateRequest) (*models.Tour, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "UPDATE tours SET "
	args := []interface{}{}
	argID := 1

	if req.CategoryID != nil {
		query += fmt.Sprintf("category_id = $%d, ", argID)
		args = append(args, *req.CategoryID)
		argID++
	}
	if req.Title != "" {
		query += fmt.Sprintf("title = $%d, ", argID)
		args = append(args, req.Title)
		argID++
	}
	if req.Description != "" {
		query += fmt.Sprintf("description = $%d, ", argID)
		args = append(args, req.Description)
		argID++
	}
	if req.Location != "" {
		query += fmt.Sprintf("location = $%d, ", argID)
		args = append(args, req.Location)
		argID++
	}
	if req.Latitude != nil {
		query += fmt.Sprintf("latitude = $%d, ", argID)
		args = append(args, req.Latitude)
		argID++
	}
	if req.Longitude != nil {
		query += fmt.Sprintf("longitude = $%d, ", argID)
		args = append(args, req.Longitude)
		argID++
	}
	if req.DurationMinutes != nil {
		query += fmt.Sprintf("duration_minutes = $%d, ", argID)
		args = append(args, *req.DurationMinutes)
		argID++
	}
	if req.PricePerPerson != nil {
		query += fmt.Sprintf("price_per_person = $%d, ", argID)
		args = append(args, *req.PricePerPerson)
		argID++
	}
	if req.MaxGuests != nil {
		query += fmt.Sprintf("max_guests = $%d, ", argID)
		args = append(args, *req.MaxGuests)
		argID++
	}
	if req.Difficulty != "" {
		query += fmt.Sprintf("difficulty = $%d, ", argID)
		args = append(args, req.Difficulty)
		argID++
	}
	if req.Languages != nil {
		query += fmt.Sprintf("languages = $%d, ", argID)
		args = append(args, req.Languages)
		argID++
	}
	if req.WhatIncluded != nil {
		query += fmt.Sprintf("what_included = $%d, ", argID)
		args = append(args, req.WhatIncluded)
		argID++
	}
	if req.MeetingPoint != "" {
		query += fmt.Sprintf("meeting_point = $%d, ", argID)
		args = append(args, req.MeetingPoint)
		argID++
	}
	if req.IsPublished != nil {
		query += fmt.Sprintf("is_published = $%d, ", argID)
		args = append(args, *req.IsPublished)
		argID++
	}

	if len(args) > 0 {
		query = query[:len(query)-2]
		query += fmt.Sprintf(" WHERE id = $%d", argID)
		args = append(args, id)

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("update tour core: %w", err)
		}
	}

	// If images are provided in the update, replace them
	if req.Images != nil {
		_, err = tx.Exec(ctx, "DELETE FROM tour_images WHERE tour_id = $1", id)
		if err != nil {
			return nil, fmt.Errorf("clear tour images: %w", err)
		}
		for i, imgURL := range req.Images {
			_, err = tx.Exec(ctx, `
				INSERT INTO tour_images (tour_id, url, position)
				VALUES ($1, $2, $3)
			`, id, imgURL, i+1)
			if err != nil {
				return nil, fmt.Errorf("insert updated tour image %d: %w", i, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit update tx: %w", err)
	}

	return r.GetByID(ctx, id, 0)
}

func (r *PgTourRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Bookings have ON DELETE RESTRICT to tour_id, so we must delete them explicitly
	_, err = tx.Exec(ctx, "DELETE FROM bookings WHERE tour_id = $1", id)
	if err != nil {
		return fmt.Errorf("delete related bookings: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM tours WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete tour: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *PgTourRepository) List(ctx context.Context, filter models.TourFilter) ([]models.Tour, error) {
	query := `
		SELECT t.id, t.guide_id, t.category_id, t.title, t.description, t.location, t.latitude, t.longitude,
		       t.duration_minutes, t.price_per_person, t.max_guests, t.difficulty, t.languages, t.what_included,
		       t.meeting_point, t.is_published, t.avg_rating, t.review_count, t.created_at, t.updated_at,
		       u.first_name || ' ' || u.last_name as guide_name, u.avatar_url as guide_avatar,
		       c.name as category_name, c.slug as category_slug,
		       EXISTS(SELECT 1 FROM favorites f WHERE f.tour_id = t.id AND f.user_id = $1) as is_favorited
		FROM tours t
		JOIN users u ON t.guide_id = u.id
		JOIN categories c ON t.category_id = c.id
		WHERE 1=1
	`
	args := []interface{}{filter.UserID}
	argID := 2

	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND t.category_id = $%d", argID)
		args = append(args, filter.CategoryID)
		argID++
	}
	if filter.CategorySlug != "" {
		query += fmt.Sprintf(" AND c.slug = $%d", argID)
		args = append(args, filter.CategorySlug)
		argID++
	}
	if filter.Difficulty != "" {
		query += fmt.Sprintf(" AND t.difficulty = $%d", argID)
		args = append(args, filter.Difficulty)
		argID++
	}
	if filter.GuideID > 0 {
		query += fmt.Sprintf(" AND t.guide_id = $%d", argID)
		args = append(args, filter.GuideID)
		argID++
	} else {
		query += " AND t.is_published = true"
	}
	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND t.price_per_person >= $%d", argID)
		args = append(args, *filter.MinPrice)
		argID++
	}
	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND t.price_per_person <= $%d", argID)
		args = append(args, *filter.MaxPrice)
		argID++
	}
	if filter.Location != "" {
		query += fmt.Sprintf(" AND t.location ILIKE $%d", argID)
		args = append(args, "%"+filter.Location+"%")
		argID++
	}
	if filter.Query != "" {
		query += fmt.Sprintf(" AND to_tsvector('spanish', coalesce(t.title, '') || ' ' || coalesce(t.description, '') || ' ' || coalesce(t.location, '')) @@ plainto_tsquery('spanish', $%d)", argID)
		args = append(args, filter.Query)
		argID++
	}

	limit := 50
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	query += fmt.Sprintf(" ORDER BY t.avg_rating DESC, t.id DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tours query: %w", err)
	}
	defer rows.Close()

	var tours []models.Tour
	tourIDs := []int{}
	for rows.Next() {
		var t models.Tour
		err := rows.Scan(
			&t.ID, &t.GuideID, &t.CategoryID, &t.Title, &t.Description, &t.Location, &t.Latitude, &t.Longitude,
			&t.DurationMinutes, &t.PricePerPerson, &t.MaxGuests, &t.Difficulty, &t.Languages, &t.WhatIncluded,
			&t.MeetingPoint, &t.IsPublished, &t.AvgRating, &t.ReviewCount, &t.CreatedAt, &t.UpdatedAt,
			&t.GuideName, &t.GuideAvatar, &t.CategoryName, &t.CategorySlug, &t.IsFavorited,
		)
		if err != nil {
			return nil, err
		}
		t.Images = []string{}
		tours = append(tours, t)
		tourIDs = append(tourIDs, t.ID)
	}

	if len(tourIDs) > 0 {
		// Bulk fetch images to avoid N+1 query problem
		imgQuery := `SELECT tour_id, url FROM tour_images WHERE tour_id = ANY($1) ORDER BY position ASC`
		imgRows, err := r.pool.Query(ctx, imgQuery, tourIDs)
		if err != nil {
			return nil, fmt.Errorf("fetch list images: %w", err)
		}
		defer imgRows.Close()

		imgMap := make(map[int][]string)
		for imgRows.Next() {
			var tID int
			var url string
			if err := imgRows.Scan(&tID, &url); err == nil {
				imgMap[tID] = append(imgMap[tID], url)
			}
		}

		for i := range tours {
			if imgs, ok := imgMap[tours[i].ID]; ok {
				tours[i].Images = imgs
			}
		}
	}

	return tours, nil
}

func (r *PgTourRepository) GetStats(ctx context.Context, guideID int) (*models.GuideStats, error) {
	query := `
		SELECT 
			COALESCE(COUNT(DISTINCT t.id), 0) as total_tours,
			COALESCE(COUNT(DISTINCT t.id) FILTER (WHERE t.is_published = true), 0) as published_tours,
			COALESCE(COUNT(DISTINCT b.id), 0) as total_bookings,
			COALESCE(COUNT(DISTINCT b.id) FILTER (WHERE b.status = 'pending'), 0) as pending_bookings,
			COALESCE(SUM(b.total_price) FILTER (WHERE b.status = 'confirmed' OR b.status = 'completed'), 0.00) as total_revenue,
			COALESCE(AVG(t.avg_rating), 0.00) as avg_rating,
			COALESCE(SUM(t.review_count), 0) as total_reviews
		FROM tours t
		LEFT JOIN bookings b ON t.id = b.tour_id
		WHERE t.guide_id = $1
	`
	stats := &models.GuideStats{}
	err := r.pool.QueryRow(ctx, query, guideID).Scan(
		&stats.TotalTours, &stats.PublishedTours, &stats.TotalBookings, &stats.PendingBookings,
		&stats.TotalRevenue, &stats.AvgRating, &stats.TotalReviews,
	)
	if err != nil {
		return nil, fmt.Errorf("guide stats: %w", err)
	}
	return stats, nil
}
