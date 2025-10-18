package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/EnduranNSU/end-user-info/internal/adapter/out/postgres/gen"
	"github.com/EnduranNSU/end-user-info/internal/domain"
)

type UserInfoRepositoryImpl struct {
	db      *sql.DB
	queries *gen.Queries
}

// NewUserInfoRepository теперь возвращает только репозиторий без ошибки
func NewUserInfoRepository(db *sql.DB) domain.UserInfoRepository {
	return &UserInfoRepositoryImpl{
		db:      db,
		queries: gen.New(db),
	}
}

// Остальные методы остаются без изменений...
func (r *UserInfoRepositoryImpl) CreateUserInfo(ctx context.Context, info *domain.UserInfo) error {
	params := gen.CreateUserInfoParams{
		Weight: info.Weight,
		Height: info.Height,
		Date:   info.Date,
		Age:    info.Age,
		UserID: info.UserID,
	}

	err := r.queries.CreateUserInfo(ctx, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("operation", "CreateUserInfo").
			Float64("weight", info.Weight).
			Int64("height", info.Height).
			Int64("age", info.Age).
			Str("user_id", info.UserID.String()).
			Msg("failed to create user info")
		return err
	}

	log.Debug().
		Str("operation", "CreateUserInfo").
		Float64("weight", info.Weight).
		Int64("height", info.Height).
		Int64("age", info.Age).
		Str("user_id", info.UserID.String()).
		Msg("successfully created user info")

	return nil
}

func (r *UserInfoRepositoryImpl) GetLatestUserInfoByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserInfo, error) {
	info, err := r.queries.GetLatestUserInfoByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		log.Warn().
			Str("operation", "GetLatestUserInfoByUserID").
			Str("user_id", userID.String()).
			Msg("user info not found")
		return nil, fmt.Errorf("user_info not found for user_id: %s", userID)
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("operation", "GetLatestUserInfoByUserID").
			Str("user_id", userID.String()).
			Msg("failed to get latest user info")
		return nil, err
	}

	domainInfo := &domain.UserInfo{
		ID:     info.ID,
		Weight: info.Weight,
		Height: info.Height,
		Date:   info.Date,
		Age:    info.Age,
		UserID: info.UserID,
	}

	log.Debug().
		Str("operation", "GetLatestUserInfoByUserID").
		Str("user_id", userID.String()).
		Int64("info_id", info.ID).
		Time("date", info.Date).
		Msg("successfully retrieved latest user info")

	return domainInfo, nil
}

func (r *UserInfoRepositoryImpl) GetAllUserInfoByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserInfo, error) {
	infos, err := r.queries.GetAllUserInfoByUserID(ctx, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("operation", "GetAllUserInfoByUserID").
			Str("user_id", userID.String()).
			Msg("failed to query user info")
		return nil, err
	}

	if len(infos) == 0 {
		log.Warn().
			Str("operation", "GetAllUserInfoByUserID").
			Str("user_id", userID.String()).
			Msg("no user info records found")
		return nil, fmt.Errorf("no user_info records found for user_id: %s", userID)
	}

	domainInfos := make([]*domain.UserInfo, len(infos))
	for i, info := range infos {
		domainInfos[i] = &domain.UserInfo{
			ID:     info.ID,
			Weight: info.Weight,
			Height: info.Height,
			Date:   info.Date,
			Age:    info.Age,
			UserID: info.UserID,
		}
	}

	log.Debug().
		Str("operation", "GetAllUserInfoByUserID").
		Str("user_id", userID.String()).
		Int("records_count", len(domainInfos)).
		Msg("successfully retrieved user info records")

	return domainInfos, nil
}
