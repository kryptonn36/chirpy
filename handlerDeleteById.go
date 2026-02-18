package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kryptonn36/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request){
	tokenString, err := auth.GetBearerToken(r.Header)
	if err!= nil {
		respondWithError(w, 401, "error in getting bearer token to delete chirp" ,err)
		return 
	}
	
	userId, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err!= nil {
		respondWithError(w, 403, "Not a valid user to perform this action" ,err)
		return 
	}
	
	pathValue := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(pathValue)
	if err!= nil {
		respondWithError(w, 403, "parsing problem in path value" ,err)
		return
	}
	chirp, err := cfg.queries.GetChirpById(r.Context(), chirpId)
	if err!= nil {
		respondWithError(w, 404, "Chirp not found" ,err)
		return
	}
	if userId != chirp.UserID{
		respondWithError(w, 403, "Not a valid user to delete this chirp", nil)
		return 
	}
	err = cfg.queries.DeleteChirpById(r.Context(), chirp.ID)
	if err!= nil {
		respondWithError(w, 403, "error in deleting chirp by id" ,err)
		return
	}
	respondWithJSON(w, 204, nil)
}