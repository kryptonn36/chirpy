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