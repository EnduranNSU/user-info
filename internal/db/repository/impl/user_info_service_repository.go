package impl

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/EnduranNSU/end-user-info/internal/db/entity"
	"github.com/EnduranNSU/end-user-info/internal/db/repository"
)

//go:embed query/*
var queryFS embed.FS

type UserInfoRepositoryImpl struct {
	db      *sql.DB
	queries map[string]string
}

func NewUserRepository(db *sql.DB) (repository.UserInfoRepository, error) {
	repo := &UserInfoRepositoryImpl{
		db:      db,
		queries: make(map[string]string),
	}
	queries := []string{
		"create_user_info",
		"get_latest_user_info_by_user_id",
		"get_all_user_info_by_user_id",
	}

	for _, queryName := range queries {
		query, err := queryFS.ReadFile(fmt.Sprintf("query/%s.sql", queryName))
		if err != nil {
			return nil, fmt.Errorf("failed to load query %s: %w", queryName, err)
		}
		repo.queries[queryName] = string(query)
	}

	return repo, nil
}

func (r *UserInfoRepositoryImpl) CreateUserInfo(ctx context.Context, info *entity.UserInfo) error {
	query := r.queries["create_user_info"]

	_, err := r.db.ExecContext(ctx, query,
		info.Weight, info.Height, time.Now(),
		info.Age, info.UserID,
	)

	log.Debug().Msgf("executing insert into user_info for: %v", info)

	return err
}

func (r *UserInfoRepositoryImpl) GetLatestUserInfoByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error) {
	query := r.queries["get_latest_user_info_by_user_id"]

	var info entity.UserInfo
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&info.ID, &info.Weight, &info.Height,
		&info.Date, &info.Age, &info.UserID,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user_info not found for user_id: %s", userID)
	}

	if err != nil {
		log.Error().Err(err).Msgf("failed to get latest user_info for user_id: %s", userID)
		return nil, err
	}

	log.Debug().Msgf("retrieved latest user_info for user_id: %s, date: %s", userID, info.Date.Format(time.RFC3339))
	return &info, nil
}

func (r *UserInfoRepositoryImpl) GetAllUserInfoByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserInfo, error) {
	query := r.queries["get_all_user_info_by_user_id"]

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Error().Err(err).Msgf("failed to query all user_info for user_id: %s", userID)
		return nil, err
	}
	defer rows.Close()

	var infos []*entity.UserInfo
	for rows.Next() {
		var info entity.UserInfo
		if err := rows.Scan(
			&info.ID, &info.Weight, &info.Height,
			&info.Date, &info.Age, &info.UserID,
		); err != nil {
			log.Error().Err(err).Msgf("failed to scan user_info row for user_id: %s", userID)
			return nil, err
		}
		infos = append(infos, &info)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msgf("error iterating user_info rows for user_id: %s", userID)
		return nil, err
	}

	log.Debug().Msgf("retrieved %d user_info records for user_id: %s", len(infos), userID)

	if len(infos) == 0 {
		return nil, fmt.Errorf("no user_info records found for user_id: %s", userID)
	}

	return infos, nil
}
