package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ziyadovea/task_manager/users/internal/app/entity"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) userRepository {
	return userRepository{db: db}
}

func (r userRepository) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	if err := r.db.QueryRowxContext(ctx, query, user.Name, user.Email, user.Password).Scan(&user.ID); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r userRepository) GetUserByID(ctx context.Context, id int64) (entity.User, error) {
	const query = `
		SELECT 
			id "id",
			name "name",
			email "email",
			password "password"
		FROM
			users
		WHERE
			id = $1
	`
	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r userRepository) GetUserByName(ctx context.Context, name string) (entity.User, error) {
	const query = `
		SELECT 
			id "id",
			name "name",
			email "email",
			password "password"
		FROM
			users
		WHERE
			name = $1
	`
	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, name); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	const query = `
		SELECT 
			id "id",
			name "name",
			email "email",
			password "password"
		FROM
			users
		WHERE
			email = $1
	`
	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r userRepository) ListUsers(ctx context.Context) ([]entity.User, error) {
	const query = `
		SELECT 
		    id "id",
			name "name",
			email "email",
			password "password"
		FROM
			users
	`
	var users []entity.User
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}
	return users, nil
}

func (r userRepository) UpdateUserByID(ctx context.Context, user entity.User) (int64, error) {
	ub := sq.Update("users").Where(sq.Eq{"id": user.ID}).PlaceholderFormat(sq.Dollar)
	if user.Name != "" {
		ub = ub.Set("name", user.Name)
	}
	if user.Email != "" {
		ub = ub.Set("email", user.Email)
	}
	if user.Password != "" {
		ub = ub.Set("password", user.Password)
	}

	query, args, err := ub.ToSql()
	if err != nil {
		return 0, fmt.Errorf("unable to build sql query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("unable to exec sql query: %w", err)
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to get affected rows: %w", err)
	}

	return rowsUpdated, nil
}

func (r userRepository) RemoveUserByID(ctx context.Context, id int64) (int64, error) {
	const query = `DELETE FROM users where id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}
