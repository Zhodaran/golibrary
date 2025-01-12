package repository

import (
	"context"
	"database/sql"
)

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

// Create добавляет нового пользователя в базу данных
func (r *PostgresBookRepository) Create(ctx context.Context, user User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
	return err
}

// GetByID получает пользователя по ID
func (r *PostgresBookRepository) GetByID(ctx context.Context, id string) (User, error) {
	var user User
	query := "SELECT id, name, email, deleted_at FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Update обновляет данные пользователя
func (r *PostgresBookRepository) Update(ctx context.Context, user User) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

// Delete помечает пользователя как удаленного
func (r *PostgresBookRepository) Delete(ctx context.Context, id string) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List возвращает список пользователей с пагинацией
func (r *PostgresBookRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
	query := "SELECT id, name, email, deleted_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
