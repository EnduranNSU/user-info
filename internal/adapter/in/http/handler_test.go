package httpin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dto "github.com/EnduranNSU/end-user-info/internal/adapter/in/http/dto"
	"github.com/EnduranNSU/end-user-info/internal/domain"
	"github.com/EnduranNSU/end-user-info/internal/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestUserInfoHandler_Create_Success(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.POST("/api/v1/user-info", handler.Create)

	userID := uuid.New()
	reqBody := dto.CreateUserInfoRequest{
		UserID: userID.String(),
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   "2024-05-10",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/user-info", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expected := &domain.CreateUserInfoCmd{
		UserID: userID,
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   "2024-05-10",
	}

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(cmd domain.CreateUserInfoCmd) bool {
		return cmd.UserID == expected.UserID &&
			cmd.Weight == expected.Weight &&
			cmd.Height == expected.Height &&
			cmd.Age == expected.Age &&
			cmd.Date == expected.Date
	})).Return(&domain.UserInfo{
		UserID: userID,
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
	}, nil)

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserInfoResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "70.5", fmt.Sprint(resp.Weight))
	assert.Equal(t, "2024-05-10", resp.Date)
	mockSvc.AssertExpectations(t)
}

func TestUserInfoHandler_Create_InvalidUUID(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.POST("/api/v1/user-info", handler.Create)

	reqBody := dto.CreateUserInfoRequest{
		UserID: "invalid-uuid", // Теперь это string
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   "2024-05-10",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/user-info", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ginEngine.ServeHTTP(w, req)
	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "bad json")
	mockSvc.AssertNotCalled(t, "Create")
}

func TestUserInfoHandler_Create_BadJSON(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.POST("/api/v1/user-info", handler.Create)

	req := httptest.NewRequest("POST", "/api/v1/user-info", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "bad json")
	mockSvc.AssertNotCalled(t, "Create")
}

func TestUserInfoHandler_GetLatest_Success(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info/latest", handler.GetLatest)

	userID := uuid.New()
	expected := &domain.UserInfo{
		UserID: userID,
		Weight: 70.5,
		Height: 175,
		Age:    25,
		Date:   time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
	}

	req := httptest.NewRequest("GET", "/api/v1/user-info/latest?user_id="+userID.String(), nil)
	w := httptest.NewRecorder()

	mockSvc.On("GetLatest", mock.Anything, userID).Return(expected, nil)

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserInfoResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "70.5", fmt.Sprint(resp.Weight))
	assert.Equal(t, "175", fmt.Sprint(resp.Height))
	assert.Equal(t, "25", fmt.Sprint(resp.Age))
	assert.Equal(t, "2024-05-10", resp.Date)
	mockSvc.AssertExpectations(t)
}

func TestUserInfoHandler_GetLatest_InvalidUUID(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info/latest", handler.GetLatest)

	req := httptest.NewRequest("GET", "/api/v1/user-info/latest?user_id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "invalid user_id")
	mockSvc.AssertNotCalled(t, "GetLatest")
}

func TestUserInfoHandler_GetLatest_NotFound(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info/latest", handler.GetLatest)

	userID := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/user-info/latest?user_id="+userID.String(), nil)
	w := httptest.NewRecorder()

	mockSvc.On("GetLatest", mock.Anything, userID).Return((*domain.UserInfo)(nil), errors.New("not found"))

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "not found")
	mockSvc.AssertExpectations(t)
}

func TestUserInfoHandler_List_Success(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info", handler.List)

	userID := uuid.New()
	items := []*domain.UserInfo{
		{
			UserID: userID,
			Weight: 70.0,
			Height: 170,
			Age:    20,
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			UserID: userID,
			Weight: 75.0,
			Height: 175,
			Age:    21,
			Date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	req := httptest.NewRequest("GET", "/api/v1/user-info?user_id="+userID.String(), nil)
	w := httptest.NewRecorder()

	mockSvc.On("List", mock.Anything, userID).Return(items, nil)

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []dto.UserInfoResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, "70", fmt.Sprint(resp[0].Weight))
	assert.Equal(t, "2024-01-01", resp[0].Date)
	assert.Equal(t, "75", fmt.Sprint(resp[1].Weight))
	assert.Equal(t, "2024-02-01", resp[1].Date)
	mockSvc.AssertExpectations(t)
}

func TestUserInfoHandler_List_InvalidUUID(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info", handler.List)

	req := httptest.NewRequest("GET", "/api/v1/user-info?user_id=invalid-uuid", nil)
	w := httptest.NewRecorder()

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "invalid user_id")
	mockSvc.AssertNotCalled(t, "List")
}

func TestUserInfoHandler_List_NotFound(t *testing.T) {
	// Arrange
	mockSvc := new(mocks.Service)
	handler := NewUserInfoHandler(mockSvc)

	ginEngine := gin.New()
	ginEngine.GET("/api/v1/user-info", handler.List)

	userID := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/user-info?user_id="+userID.String(), nil)
	w := httptest.NewRecorder()

	mockSvc.On("List", mock.Anything, userID).Return([]*domain.UserInfo(nil), errors.New("not found"))

	// Act
	ginEngine.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp dto.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp.Error, "not found")
	mockSvc.AssertExpectations(t)
}
