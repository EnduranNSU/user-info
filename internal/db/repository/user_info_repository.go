package repository

import (
	"context"

	"github.com/EnduranNSU/end-user-info/internal/db/entity"
	"github.com/google/uuid"
)

type UserInfoRepository interface {
	CreateUserInfo(ctx context.Context, info *entity.UserInfo) error
	GetLatestUserInfoByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error)
	GetAllUserInfoByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserInfo, error)
}