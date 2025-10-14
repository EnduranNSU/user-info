package domain

import (
    "time"
    "github.com/google/uuid"
)

type UserInfo struct {
    ID     int64     `db:"id" json:"id"`
    Weight float64   `db:"weight" json:"weight"`
    Height int64     `db:"height" json:"height"`
    Date   time.Time `db:"date" json:"date"`
    Age    int64     `db:"age" json:"age"`
    UserID uuid.UUID `db:"user_id" json:"user_id"`
}