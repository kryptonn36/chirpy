package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/kryptonn36/chirpy/internal/auth"
)

type reqWebhook struct{
	Event string `json:"event"`
	Data struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) hanlderChirpRed(w http.ResponseWriter, r *http.Request){
	apiKey, err := auth.GetAPIKey(r.Header)
	if err!=nil{
		w.WriteHeader(401)
		return 
	}
	if apiKey!= os.Getenv("POLKA_KEY"){
		w.WriteHeader(401)
		return 
	}
	params := reqWebhook{}
	decode := json.NewDecoder(r.Body)
	err = decode.Decode(&params)
	if err!=nil{
		respondWithError(w, 404, "error in webhook request body", err)
		return
	}
	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return 
	}
	err = cfg.queries.UpdateToChirpyRed(r.Context(), params.Data.UserId)
	if err!=nil{
		w.WriteHeader(http.StatusNoContent)
		return 
	}
	respondWithJSON(w, 204, nil)
}