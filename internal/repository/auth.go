package repository

import (
	"context"
	"database/sql"
	"errors"
	"medods-test-task/internal/db/postgres"
	"medods-test-task/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Auth struct {
	db postgres.DB
}

func NewAuthRepo(db postgres.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func (r *Auth) CreateSession(ctx context.Context, session *models.RefreshSession) error {
	_, err := sq.
		Insert("refreshSessions").
		Columns("userId", "ip", "refreshToken", "expiresAt", "createdAt").
		Values(session.UserID, session.IP, session.Token, session.ExpiresAt, session.CreatedAt).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Auth) DeleteSessionByUserID(ctx context.Context, userID uuid.UUID) error {
	res, err := sq.
		Delete("refreshSessions").
		Where(sq.Eq{"userId": userID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrSessionNotFound
	}

	return nil
}

func (r *Auth) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*models.RefreshSession, error) {
	row := sq.
		Select("id", "userId", "ip", "refreshToken", "expiresAt", "createdAt").
		From("refreshSessions").
		Where(sq.Eq{"userId": userID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow()

	var session models.RefreshSession

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.IP,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrSessionNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (r *Auth) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	row := sq.
		Select("id").
		From("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow()

	var user models.User

	err := row.Scan(
		&user.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Auth) CreateUser(ctx context.Context, user *models.User) error {
	res, err := sq.
		Insert("users").
		Columns("id").
		Values(user.ID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrSessionNotFound
	}

	return nil
}
