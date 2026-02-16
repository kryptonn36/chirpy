package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings" 

	// "github.com/google/uuid"
	"github.com/kryptonn36/chirpy/internal/auth"
	"github.com/kryptonn36/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		// UserId uuid.UUID `json:"user_id"`
	}
	token_string, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, 404, "Bear token not found", err)
		return 
	}

	id, err := auth.ValidateJWT(token_string, cfg.secret)
	if err!=nil{
		respondWithError(w, 404, "error in getting id in handler chirps", err)
		return 
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// if id!=params.UserId{
	// 	respondWithError(w, http.StatusUnauthorized, "Unothorized access", nil)
	// 	return 
	// }

	cleaned, err := cleanedUp(params.Body)
	if err!=nil{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: id,
	})
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "error in creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp_return{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	})
}


func cleanedUp(body string) (string, error){
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", fmt.Errorf("to much long chirp")
	}

	// var value []string
	words := strings.Split(body, " ")
	for i, word := range words {
		check := strings.ToLower(word)
		if check == "kerfuffle" || check == "sharbert" || check == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words," "), nil
}