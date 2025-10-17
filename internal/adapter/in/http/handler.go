package httpin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EnduranNSU/end-user-info/internal/adapter/in/http/dto"
	"github.com/EnduranNSU/end-user-info/internal/domain"
)

type UserInfoHandler struct {
	repo domain.UserInfoRepository
}

func NewUserInfoHandler(repo domain.UserInfoRepository) *UserInfoHandler {
	return &UserInfoHandler{repo: repo}
}

func (h *UserInfoHandler) Create(c *gin.Context) {
	var req dto.CreateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}
	if req.UserID == "" || req.Weight <= 0 || req.Height <= 0 || req.Age <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "missing or invalid fields"})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	now := time.Now().UTC()

	m := &domain.UserInfo{
		UserID: uid,
		Weight: req.Weight,
		Height: req.Height,
		Age:    req.Age,
		Date:   now,
	}
	if err := h.repo.CreateUserInfo(c.Request.Context(), m); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create"})
		return
	}

	c.JSON(http.StatusOK, dto.UserInfoResponse{
		Weight: m.Weight,
		Height: m.Height,
		Age:    m.Age,
		Date:   m.Date.Format("2006-01-02"),
	})
}

func (h *UserInfoHandler) GetLatest(c *gin.Context) {
	uidStr := c.Query("user_id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}
	m, err := h.repo.GetLatestUserInfoByUserID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
		return
	}
	c.JSON(http.StatusOK, dto.UserInfoResponse{
		Weight: m.Weight,
		Height: m.Height,
		Age:    m.Age,
		Date:   m.Date.Format("2006-01-02"),
	})
}

func (h *UserInfoHandler) List(c *gin.Context) {
	uidStr := c.Query("user_id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}
	items, err := h.repo.GetAllUserInfoByUserID(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
		return
	}
	resp := make([]dto.UserInfoResponse, 0, len(items))
	for _, m := range items {
		resp = append(resp, dto.UserInfoResponse{
			Weight: m.Weight,
			Height: m.Height,
			Age:    m.Age,
			Date:   m.Date.Format("2006-01-02"),
		})
	}
	c.JSON(http.StatusOK, resp)
}
