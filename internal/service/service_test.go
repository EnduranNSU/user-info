package service

import (
	"context"
	"testing"
	"time"

	"github.com/EnduranNSU/end-user-info/internal/domain"
	"github.com/EnduranNSU/end-user-info/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Create_ValidDate(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.UserInfoRepository)
	svc := New(mockRepo)

	userID := uuid.New()
	cmd := domain.CreateUserInfoCmd{
		UserID: userID,
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   "2024-05-10",
	}

	expected := &domain.UserInfo{
		UserID: userID,
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
	}

	mockRepo.On("CreateUserInfo", mock.Anything, mock.MatchedBy(func(info *domain.UserInfo) bool {
		return info.UserID == userID &&
			info.Weight == 70.5 &&
			info.Height == 175 &&
			info.Age == 25 &&
			info.Date.Format("2006-01-02") == "2024-05-10"
	})).Return(nil)

	result, err := svc.Create(context.Background(), cmd)

	assert.NoError(t, err)
	assert.Equal(t, expected.UserID, result.UserID)
	assert.Equal(t, expected.Weight, result.Weight)
	assert.Equal(t, expected.Date.Format("2006-01-02"), result.Date.Format("2006-01-02"))
	mockRepo.AssertExpectations(t)
}

func TestService_Create_EmptyDate_UseNow(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.UserInfoRepository)
	svc := New(mockRepo)

	userID := uuid.New()
	cmd := domain.CreateUserInfoCmd{
		UserID: userID,
		Weight: 70.0,
		Height: 180,
		Age:    30,
		Date:   "",
	}

	now := time.Now().UTC()
	mockRepo.On("CreateUserInfo", mock.Anything, mock.MatchedBy(func(info *domain.UserInfo) bool {
		diff := info.Date.Sub(now)
		return info.UserID == userID &&
			info.Weight == 70.0 &&
			info.Height == 180 &&
			info.Age == 30 &&
			diff < time.Second
	})).Return(nil)

	result, err := svc.Create(context.Background(), cmd)

	assert.NoError(t, err)
	assert.WithinDuration(t, now, result.Date, time.Second)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_InvalidDate_ReturnsError(t *testing.T) {
	mockRepo := new(mocks.UserInfoRepository)
	svc := New(mockRepo)

	cmd := domain.CreateUserInfoCmd{
		UserID: uuid.New(),
		Date:   "invalid-date",
	}

	result, err := svc.Create(context.Background(), cmd)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidDate, err)
	assert.Nil(t, result)
}

func TestService_GetLatest_Success(t *testing.T) {
	mockRepo := new(mocks.UserInfoRepository)
	svc := New(mockRepo)

	userID := uuid.New()
	expected := &domain.UserInfo{ID: 1, UserID: userID, Weight: 70.0}

	mockRepo.On("GetLatestUserInfoByUserID", mock.Anything, userID).Return(expected, nil)

	result, err := svc.GetLatest(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestService_List_Success(t *testing.T) {
	mockRepo := new(mocks.UserInfoRepository)
	svc := New(mockRepo)

	userID := uuid.New()
	expected := []*domain.UserInfo{
		{ID: 1, UserID: userID, Weight: 70.0},
		{ID: 2, UserID: userID, Weight: 75.0},
	}

	mockRepo.On("GetAllUserInfoByUserID", mock.Anything, userID).Return(expected, nil)

	result, err := svc.List(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}