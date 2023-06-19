package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nullblocks/rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello Rss feed reader ")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in environment")
	}
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB-URL not found in environment")
	}

	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("can't connect to DB ", err)
	}

	apiCFG := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.HandleFunc("/ready", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCFG.handlerCreateUser)
	v1Router.Get("/users", apiCFG.middlewareAuth(apiCFG.handlerUsersGet))
	v1Router.Post("/feeds", apiCFG.middlewareAuth(apiCFG.handlerFeedCreate))
	v1Router.Get("/feeds", apiCFG.handlerGetFeeds)

	// v1Router.Post("/feed_follow", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowCreate))
	v1Router.Get("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCFG.middlewareAuth(apiCFG.handlerFeedFollowDelete))

	v1Router.Get("/feeds", apiCFG.handlerGetFeeds)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Println("Port running on  : ", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
