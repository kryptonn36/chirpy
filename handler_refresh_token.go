package main

import (
	"net/http"
	"time"

	"github.com/kryptonn36/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request){
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, 401, "error in getting bearer token", err)
		return
	}
	databaseRefreshToken, err := cfg.queries.GetRefreshToken(r.Context(), requestRefreshToken)
	if err!=nil{
		respondWithError(w, 401, "refresh token is not in database", err)
		return
	}
	if time.Now().After(databaseRefreshToken.ExpiresAt){
		respondWithError(w, 401, "the refresh token is expired", nil)
		return
	}
	
	if databaseRefreshToken.RevokedAt.Valid{
		respondWithError(w, 401, "the refresh token is revoked", nil)
		return
	}

	expiring_time := time.Hour
	jwtToken, err := auth.MakeJWT(databaseRefreshToken.UserID, cfg.secret, expiring_time)
	if err!=nil{
		respondWithError(w,404, "Error in creating JWT Token", err)
		return
	}
	respondWithJSON(w, 200, returnVals{
		Token: jwtToken,
	})
}


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "bearer token not error", err)
		return
	}
	databaseRefreshToken, err := cfg.queries.GetRefreshToken(r.Context(), requestRefreshToken)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "refresh token not found in database", err)
		return
	}
	err = cfg.queries.RevoketimeUpdate(r.Context(), databaseRefreshToken.Token)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "error in revoking refresh token", err)
		return
	}
	respondWithJSON(w, 204,nil)
}