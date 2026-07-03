package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type PgTourScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewPgTourScheduleRepository(pool *pgxpool.Pool) TourScheduleRepository {
	return &PgTourScheduleRepository{pool: pool}
}

func (r *PgTourScheduleRepository) Create(ctx context.Context, s models.TourSchedule) (*models.TourSchedule, error) {
	query := `
		INSERT INTO tour_schedules (tour_id, start_time, end_time, available_spots, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, tour_id, start_time, end_time, available_spots, is_active, created_at
	`
	res := &models.TourSchedule{}
	err := r.pool.QueryRow(ctx, query, s.TourID, s.StartTime, s.EndTime, s.AvailableSpots, s.IsActive).Scan(
		&res.ID, &res.TourID, &res.StartTime, &res.EndTime, &res.AvailableSpots, &res.IsActive, &res.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}
	return res, nil
}

func (r *PgTourScheduleRepository) GetByID(ctx context.Context, id int) (*models.TourSchedule, error) {
	s := &models.TourSchedule{}
	err := r.pool.QueryRow(ctx, "SELECT id, tour_id, start_time, end_time, available_spots, is_active, created_at FROM tour_schedules WHERE id = $1", id).Scan(
		&s.ID, &s.TourID, &s.StartTime, &s.EndTime, &s.AvailableSpots, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}
	return s, nil
}

func (r *PgTourScheduleRepository) ListByTourID(ctx context.Context, tourID int, onlyActive bool) ([]models.TourSchedule, error) {
	query := "SELECT id, tour_id, start_time, end_time, available_spots, is_active, created_at FROM tour_schedules WHERE tour_id = $1"
	if onlyActive {
		query += " AND is_active = true AND start_time > NOW()"
	}
	query += " ORDER BY start_time ASC"

	rows, err := r.pool.Query(ctx, query, tourID)
	if err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.TourSchedule
	for rows.Next() {
		var s models.TourSchedule
		err := rows.Scan(&s.ID, &s.TourID, &s.StartTime, &s.EndTime, &s.AvailableSpots, &s.IsActive, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *PgTourScheduleRepository) Update(ctx context.Context, id int, req models.ScheduleUpdateRequest) (*models.TourSchedule, error) {
	query := "UPDATE tour_schedules SET "
	args := []interface{}{}
	argID := 1

	if req.StartTime != nil {
		query += fmt.Sprintf("start_time = $%d, ", argID)
		args = append(args, *req.StartTime)
		argID++
	}
	if req.EndTime != nil {
		query += fmt.Sprintf("end_time = $%d, ", argID)
		args = append(args, *req.EndTime)
		argID++
	}
	if req.AvailableSpots != nil {
		query += fmt.Sprintf("available_spots = $%d, ", argID)
		args = append(args, *req.AvailableSpots)
		argID++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", argID)
		args = append(args, *req.IsActive)
		argID++
	}

	if len(args) == 0 {
		return r.GetByID(ctx, id)
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, tour_id, start_time, end_time, available_spots, is_active, created_at", argID)
	args = append(args, id)

	s := &models.TourSchedule{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&s.ID, &s.TourID, &s.StartTime, &s.EndTime, &s.AvailableSpots, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update schedule: %w", err)
	}
	return s, nil
}

func (r *PgTourScheduleRepository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tour_schedules WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}
	return nil
}

func (r *PgTourScheduleRepository) AdjustSpots(ctx context.Context, id int, delta int) error {
	_, err := r.pool.Exec(ctx, "UPDATE tour_schedules SET available_spots = available_spots + $1 WHERE id = $2", delta, id)
	if err != nil {
		return fmt.Errorf("adjust spots: %w", err)
	}
	return nil
}
