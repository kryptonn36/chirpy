package main

import (
	"encoding/json"
	"net/http"

	"github.com/kryptonn36/chirpy/internal/auth"
	"github.com/kryptonn36/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request){
	// decode the error from json to struct to get email
	params := paramater{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError,"Error in email request", err)
	}
	if params.Email == "" || params.Password == "" {
	respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
	return
}

	// creating hash and checking password
	hash_paswd, err := auth.HashPassword(params.Password)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "error in HashPassword function", err)
	}

	// create the user with the help of api config
	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hash_paswd,
		Email: params.Email,
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "error in creating user", err)
	}

	// now respond with a json file with details
	respondWithJSON(w, 201, returnVals{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request){
	accessToken, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, 401, "access token not get by bear token function", err)
		return
	}
	userId, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err!=nil{
		respondWithError(w, 401, "Not a valid user", err)
		return
	}

	params := paramater{}
	decode := json.NewDecoder(r.Body)
	err = decode.Decode(&params)
	if err!=nil{
		respondWithError(w, 401, "error in decoding json to update password", err)
		return
	}

	hashString, err := auth.HashPassword(params.Password)
	if err!=nil {
		respondWithError(w, 401, "error in creating hash password to update password", err)
		return
	}

	// tokenString,err := auth.MakeJWT(userId, cfg.secret, time.Hour)
	// if err!=nil{
	// 	respondWithError(w, 401, "error in making jwt to update password", err)
	// }

	err = cfg.queries.UpdateEmailPassword(r.Context(), database.UpdateEmailPasswordParams{
		HashedPassword: hashString,
		Email: params.Email,
		ID: userId,
	})
	if err!=nil{
		respondWithError(w, 401, "errror in updating email and password", err)
		return
	}
	user, err:= cfg.queries.GetUserByEmail(r.Context(),params.Email)
	respondWithJSON(w, 200, returnVals{
		Id: userId,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}