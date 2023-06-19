package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nullblocks/rss-aggregator/internal/database"
)

func (apiCFG *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %s", err))
		return
	}

	user, err := apiCFG.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create new user %v", err))
		return
	}

	// respondWithJSON(w, 200, user)
	respondWithJSON(w, 200, databaseUserToUser(user))

}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	// apiKey, err := auth.GetAPIKey(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
	// 	return
	// }

	// user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	// if err != nil {
	// 	respondWithError(w, http.StatusNotFound, "Couldn't ge t user")
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
