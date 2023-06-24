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

func (apiCFG *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
func (apiCFG *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiCFG.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(10),
	})
	println("POSTS :=  ", posts)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("COuldn't get Posts : %v", err))
		return
	}
	fmt.Println("POST at handler are", posts)
	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
