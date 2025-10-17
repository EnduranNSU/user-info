package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/EnduranNSU/end-user-info/internal/domain"
)

const (
	CREATE_USER_INFO                = "create_user_info"
	GET_LATEST_USER_INFO_BY_USER_ID = "get_latest_user_info_by_user_id"
	GET_ALL_USER_INFO_BY_USER_ID    = "get_all_user_info_by_user_id"
)

//go:embed query/*
var queryFS embed.FS

type UserInfoRepositoryImpl struct {
	db      *sql.DB
	queries map[string]string
}

func NewUserInfoRepository(db *sql.DB) (domain.UserInfoRepository, error) {
	repo := &UserInfoRepositoryImpl{
		db:      db,
		queries: make(map[string]string),
	}
	queries := []string{
		CREATE_USER_INFO,
		GET_LATEST_USER_INFO_BY_USER_ID,
		GET_ALL_USER_INFO_BY_USER_ID,
	}

	for _, queryName := range queries {
		query, err := queryFS.ReadFile(fmt.Sprintf("query/%s.sql", queryName))
		if err != nil {
			log.Error().
				Err(err).
				Str("query", queryName).
				Msg("failed to load")
			return nil, err
		}
		repo.queries[queryName] = string(query)
	}

	return repo, nil
}

func (r *UserInfoRepositoryImpl) CreateUserInfo(ctx context.Context, info *domain.UserInfo) error {
	query := r.queries[CREATE_USER_INFO]

	_, err := r.db.ExecContext(ctx, query,
		info.Weight, info.Height, info.Date,
		info.Age, info.UserID,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("operation", CREATE_USER_INFO).
			Float64("weight", info.Weight).
			Int64("height", info.Height).
			Int64("age", info.Age).
			Str("user_id", info.UserID.String()).
			Msg("failed to create user info")
		return err
	}

	log.Debug().
		Str("operation", CREATE_USER_INFO).
		Float64("weight", info.Weight).
		Int64("height", info.Height).
		Int64("age", info.Age).
		Str("user_id", info.UserID.String()).
		Msg("successfully created user info")

	return nil
}

func (r *UserInfoRepositoryImpl) GetLatestUserInfoByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserInfo, error) {
	query := r.queries[GET_LATEST_USER_INFO_BY_USER_ID]

	var info domain.UserInfo
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&info.ID, &info.Weight, &info.Height,
		&info.Date, &info.Age, &info.UserID,
	)

	if err == sql.ErrNoRows {
		log.Warn().
			Str("operation", GET_LATEST_USER_INFO_BY_USER_ID).
			Str("user_id", userID.String()).
			Msg("user info not found")
		return nil, fmt.Errorf("user_info not found for user_id: %s", userID)
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("operation", GET_LATEST_USER_INFO_BY_USER_ID).
			Str("user_id", userID.String()).
			Msg("failed to get latest user info")
		return nil, err
	}

	log.Debug().
		Str("operation", GET_LATEST_USER_INFO_BY_USER_ID).
		Str("user_id", userID.String()).
		Int64("info_id", info.ID).
		Time("date", info.Date).
		Msg("successfully retrieved latest user info")

	return &info, nil
}

func (r *UserInfoRepositoryImpl) GetAllUserInfoByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserInfo, error) {
	query := r.queries[GET_ALL_USER_INFO_BY_USER_ID]

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("operation", GET_ALL_USER_INFO_BY_USER_ID).
			Str("user_id", userID.String()).
			Msg("failed to query user info")
		return nil, err
	}
	defer rows.Close()

	var infos []*domain.UserInfo
	for rows.Next() {
		var info domain.UserInfo
		if err := rows.Scan(
			&info.ID, &info.Weight, &info.Height,
			&info.Date, &info.Age, &info.UserID,
		); err != nil {
			log.Error().
				Err(err).
				Str("operation", GET_ALL_USER_INFO_BY_USER_ID).
				Str("user_id", userID.String()).
				Msg("failed to scan user info row")
			return nil, err
		}
		infos = append(infos, &info)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Str("operation", GET_ALL_USER_INFO_BY_USER_ID).
			Str("user_id", userID.String()).
			Msg("error iterating user info rows")
		return nil, err
	}

	log.Debug().
		Str("operation", GET_ALL_USER_INFO_BY_USER_ID).
		Str("user_id", userID.String()).
		Int("records_count", len(infos)).
		Msg("successfully retrieved user info records")

	if len(infos) == 0 {
		log.Warn().
			Str("operation", GET_ALL_USER_INFO_BY_USER_ID).
			Str("user_id", userID.String()).
			Msg("no user info records found")
		return nil, fmt.Errorf("no user_info records found for user_id: %s", userID)
	}

	return infos, nil
}
