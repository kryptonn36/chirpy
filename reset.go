package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// if cfg.platform != "dev"{
	// 	respondWithError(w, 403, "Forbidden", http.ErrNotSupported)
	// }
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

	err := cfg.queries.DeleteUsers(context.Background())
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "error in deleting users", err)
	}
}
