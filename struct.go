package main

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/kryptonn36/chirpy/internal/database"
)

// from main.go
type apiConfig struct {
	fileserverHits atomic.Int32
	queries *database.Queries
	platform string
	secret string
}

// handler_get_chirps.go
type paramater struct{
	Password string `json:"password"`
	Email string `json:"email"`
	Expire_in_seconds *int `json:"expires_in_seconds`
}
type returnVals struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// handler_validate.go
type chirp_return struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}