package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type FavoriteRepository interface {
	Add(ctx context.Context, userID int, tourID int) error
	Remove(ctx context.Context, userID int, tourID int) error
	List(ctx context.Context, userID int) ([]models.Tour, error)
}

type PgFavoriteRepository struct {
	pool *pgxpool.Pool
}

func NewPgFavoriteRepository(pool *pgxpool.Pool) FavoriteRepository {
	return &PgFavoriteRepository{pool: pool}
}

func (r *PgFavoriteRepository) Add(ctx context.Context, userID int, tourID int) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO favorites (user_id, tour_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, tourID)
	if err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}

func (r *PgFavoriteRepository) Remove(ctx context.Context, userID int, tourID int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM favorites WHERE user_id = $1 AND tour_id = $2", userID, tourID)
	if err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}

func (r *PgFavoriteRepository) List(ctx context.Context, userID int) ([]models.Tour, error) {
	// Note: the tours.images JSONB column was migrated into the relational
	// tour_images table by migration 000005, so we select tour core columns
	// only and fetch images separately (same pattern as tour_repo.List) to
	// avoid an N+1 query and to not reference a dropped column.
	query := `
		SELECT t.id, t.guide_id, t.category_id, t.title, t.description, t.location, t.latitude, t.longitude,
		       t.duration_minutes, t.price_per_person, t.max_guests, t.difficulty, t.languages, t.what_included,
		       t.meeting_point, t.is_published, t.avg_rating, t.review_count, t.created_at, t.updated_at,
		       u.first_name || ' ' || u.last_name as guide_name, u.avatar_url as guide_avatar,
		       c.name as category_name, c.slug as category_slug,
		       (SELECT MIN(start_time) FROM tour_schedules WHERE tour_id = t.id AND start_time > NOW() AND is_active = true) as next_schedule_start,
		       (SELECT available_spots FROM tour_schedules WHERE tour_id = t.id AND start_time > NOW() AND is_active = true ORDER BY start_time ASC LIMIT 1) as available_spots,
		       true as is_favorited
		FROM favorites f
		JOIN tours t ON f.tour_id = t.id
		JOIN users u ON t.guide_id = u.id
		JOIN categories c ON t.category_id = c.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list favorites: %w", err)
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
			&t.GuideName, &t.GuideAvatar, &t.CategoryName, &t.CategorySlug,
			&t.NextScheduleStart, &t.AvailableSpots, &t.IsFavorited,
		)
		if err != nil {
			return nil, err
		}
		t.Images = []string{}
		tours = append(tours, t)
		tourIDs = append(tourIDs, t.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list favorites rows: %w", err)
	}

	if len(tourIDs) > 0 {
		// Bulk fetch images for all favorited tours in one query, mirroring
		// tour_repo.List so favorites cards can render image galleries.
		imgQuery := `SELECT tour_id, url FROM tour_images WHERE tour_id = ANY($1) ORDER BY position ASC`
		imgRows, err := r.pool.Query(ctx, imgQuery, tourIDs)
		if err != nil {
			return nil, fmt.Errorf("fetch favorite tour images: %w", err)
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
