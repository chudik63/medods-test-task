package repository

import (
	"medods-test-task/internal/db/postgres"
	"medods-test-task/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Auth struct {
	db postgres.DB
}

func NewAuthRepo(db postgres.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func (r *Auth) Create(session *models.RefreshSession) error {
	_, err := sq.
		Insert("refreshSessions").
		Columns("userId", "refreshToken", "ip", "expiresAt", "createdAt").
		Values(session.UserID, session.Token, session.IP, session.ExpiresAt.Unix(), session.CreatedAt).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Auth) GetByToken(token string) (*models.RefreshSession, error) {
	row := sq.
		Select("id", "userId", "refreshToken", "ip", "expiresAt", "createdAt").
		From("refreshSessions").
		Where(sq.Eq{"refreshToken": token}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow()

	var session models.RefreshSession

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.IP,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Auth) DeleteExpired() error {
	_, err := sq.
		Delete("refreshSessions").
		Where("expiresAt < ?", time.Now().Unix()).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Auth) GetByUserID(userID string) (*models.RefreshSession, error) {
	row := sq.
		Select("id", "userId", "refreshToken", "ip", "expiresAt", "createdAt").
		From("refreshSessions").
		Where(sq.Eq{"userId": userID}).
		Where("expiresIn >= ?", time.Now().Unix()).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRow()

	var session models.RefreshSession

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.IP,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
