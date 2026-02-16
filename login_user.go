package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kryptonn36/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	params := paramater{}
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&params)
	if err!= nil{
		respondWithError(w, http.StatusInternalServerError, "error in decoding post request for login", err)
		return
	}

	user ,err := cfg.queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error in getting user by email: %v",err)))
		return 
	}

	boolean,err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error in checking password for authentication: %v",err)))
		return 
	}
	if boolean==false{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect Password"))
		return
	}

	expiring_time := time.Hour
	if params.Expire_in_seconds != nil{
		expiring_time = time.Duration(*params.Expire_in_seconds) * time.Second
	}
	jwtToken, err := auth.MakeJWT(user.ID, cfg.secret, expiring_time)
	if err!=nil{
		respondWithError(w,404, "Error in creating JWT Token", err)
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: jwtToken,
	})
}