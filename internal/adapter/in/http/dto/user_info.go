package dto

type CreateUserInfoRequest struct {
	Weight float64 `json:"weight"`
	Height int64   `json:"height"`
	Date   string  `json:"date,omitempty"`
	Age    int64   `json:"age"`
}
