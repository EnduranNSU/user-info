package dto

type CreateUserInfoRequest struct {
	Weight float64 `json:"weight"`
	Height int64   `json:"height"`
	Date   string  `json:"date,omitempty"`
	Age    int64   `json:"age"`
	UserID string  `json:"user_id"`
}

type UserInfoResponse struct {
	Weight float64 `json:"weight"`
	Height int64   `json:"height"`
	Date   string  `json:"date"`
	Age    int64   `json:"age"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
