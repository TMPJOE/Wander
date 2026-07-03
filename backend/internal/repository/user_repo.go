package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"wander/backend/internal/models"
)

type PgUserRepository struct {
	pool *pgxpool.Pool
}

func NewPgUserRepository(pool *pgxpool.Pool) UserRepository {
	return &PgUserRepository{pool: pool}
}

func (r *PgUserRepository) Create(ctx context.Context, req models.UserCreateRequest, hashedPassword string) (*models.User, error) {
	role := req.Role
	if role == "" {
		role = "traveler"
	}

	query := `
		INSERT INTO users (email, username, password_hash, first_name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, username, first_name, last_name, role, bio, phone, avatar_url, languages, created_at, updated_at
	`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, query,
		req.Email, req.Username, hashedPassword, req.FirstName, req.LastName, role,
	).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Role, &user.Bio, &user.Phone, &user.AvatarURL, &user.Languages,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user db: %w", err)
	}

	return user, nil
}

func (r *PgUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, role, bio, phone, avatar_url, languages, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.Bio, &user.Phone, &user.AvatarURL, &user.Languages,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}
	return user, nil
}

func (r *PgUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, role, bio, phone, avatar_url, languages, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.Bio, &user.Phone, &user.AvatarURL, &user.Languages,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, models.ErrNotFound
	}
	return user, nil
}

func (r *PgUserRepository) Update(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error) {
	query := "UPDATE users SET "
	args := []interface{}{}
	argID := 1

	if req.Email != "" {
		query += fmt.Sprintf("email = $%d, ", argID)
		args = append(args, req.Email)
		argID++
	}
	if req.Username != "" {
		query += fmt.Sprintf("username = $%d, ", argID)
		args = append(args, req.Username)
		argID++
	}
	if req.FirstName != "" {
		query += fmt.Sprintf("first_name = $%d, ", argID)
		args = append(args, req.FirstName)
		argID++
	}
	if req.LastName != "" {
		query += fmt.Sprintf("last_name = $%d, ", argID)
		args = append(args, req.LastName)
		argID++
	}
	if req.Bio != "" {
		query += fmt.Sprintf("bio = $%d, ", argID)
		args = append(args, req.Bio)
		argID++
	}
	if req.Phone != "" {
		query += fmt.Sprintf("phone = $%d, ", argID)
		args = append(args, req.Phone)
		argID++
	}
	if req.AvatarURL != "" {
		query += fmt.Sprintf("avatar_url = $%d, ", argID)
		args = append(args, req.AvatarURL)
		argID++
	}
	if req.Languages != nil {
		query += fmt.Sprintf("languages = $%d, ", argID)
		args = append(args, req.Languages)
		argID++
	}

	if len(args) == 0 {
		return r.GetByID(ctx, id)
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, email, username, first_name, last_name, role, bio, phone, avatar_url, languages, created_at, updated_at", argID)
	args = append(args, id)

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Role, &user.Bio, &user.Phone, &user.AvatarURL, &user.Languages,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update user db: %w", err)
	}

	return user, nil
}

func (r *PgUserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete user db: %w", err)
	}
	return nil
}

func (r *PgUserRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, email, username, first_name, last_name, role, bio, phone, avatar_url, languages, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users db: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
			&u.Role, &u.Bio, &u.Phone, &u.AvatarURL, &u.Languages,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
