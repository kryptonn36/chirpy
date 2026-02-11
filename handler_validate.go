package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// var value []string
	words := strings.Split(params.Body, " ")
	for i, word := range words {
		check := strings.ToLower(word)
		if check == "kerfuffle" || check == "sharbert" || check == "fornax" {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: cleaned,
	})
}
