package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/kryptonn36/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request){
	// var chirp_list []chirp_return
	authorParam := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

    var chirp_list []database.Chirp
    var err error

	if authorParam!=""{
        authorID, err := uuid.Parse(authorParam)
        if err != nil {
            respondWithError(w, 400, "invalid author_id", err)
            return
        }
		chirp_list,err = cfg.queries.ChirpByAuthor(r.Context(), authorID)
		if err!=nil{
			w.WriteHeader(500)
			w.Write([]byte (fmt.Sprintf("Error in getting chirps: %v",err)))
			return
		}
	}else if sortParam=="desc"{
		chirp_list,err = cfg.queries.GetAllChirp(r.Context())
		if err!=nil{
			w.WriteHeader(500)
			w.Write([]byte (fmt.Sprintf("Error in getting chirps: %v",err)))
			return
		}
		sort.Slice(chirp_list, func(i, j int) bool {return chirp_list[i].CreatedAt.After(chirp_list[j].CreatedAt)
    })
	}else{
		chirp_list,err = cfg.queries.GetAllChirp(r.Context())
		if err!=nil{
			w.WriteHeader(500)
			w.Write([]byte (fmt.Sprintf("Error in getting chirps: %v",err)))
			return
		}
	}

	response := make([]chirp_return, len(chirp_list))
	for i, chirp := range chirp_list{
		response[i] = chirp_return{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request){
	path_value := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(path_value)
	if err!=nil{
		w.WriteHeader(404)
		w.Write([]byte (fmt.Sprintf("Error in parsing chirpId: %v",err)))
		return
	}
	chirp, err := cfg.queries.GetChirpById(r.Context(), chirpId)
	if err!= nil{
		w.WriteHeader(404)
		w.Write([]byte (fmt.Sprintf("Error in gettin chirp by Id: %v",err)))
		return
	}
	response := chirp_return{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, response)
	
}