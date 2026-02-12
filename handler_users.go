package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "github.com/kryptonn36/chirpy/internal/database"
)

type email_r struct{
	Email string `json:"email"`
}
type returnVals struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request){
	// decode the error from json to struct to get email
	req_email := email_r{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req_email)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError,"Error in email request", err)
	}

	// create the user with the help of api config
	user, err := cfg.queries.CreateUser(r.Context(), req_email.Email)
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