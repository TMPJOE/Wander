package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
}

type PgCategoryRepository struct {
	pool *pgxpool.Pool
}

func NewPgCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &PgCategoryRepository{pool: pool}
}

func (r *PgCategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, slug, icon, description, sort_order, created_at FROM categories ORDER BY sort_order ASC")
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Icon, &c.Description, &c.SortOrder, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *PgCategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	c := &models.Category{}
	err := r.pool.QueryRow(ctx, "SELECT id, name, slug, icon, description, sort_order, created_at FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.Slug, &c.Icon, &c.Description, &c.SortOrder, &c.CreatedAt)
	if err != nil {
		return nil, models.ErrNotFound
	}
	return c, nil
}
