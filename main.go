package main

import (
	"fmt"
	"html/template"
	"log"
	"movies/connection"
	"movies/models"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var movies []models.Movie

func main() {
	db := connection.Connect()

	defer db.Close()

	movies = models.QueryMovies(db)

	// http.HandleFunc("/movie", handler)
	// http.ListenAndServe(":3000", nil)

	router := mux.NewRouter()
	router.HandleFunc("/movie/{id}", handler).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var MovieVars models.Movie

	params := mux.Vars(r)
	movieID := params["id"]

	// movieID := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(movieID)

	fmt.Println("Movie ID: ", movieID)

	movie, err := models.FindMovie(movieID)

	if err == models.ErrNoMovie {
		log.Println("Movie not found on redis and will query from database")

		// Query movie from database
		db := connection.Connect()

		defer db.Close()

		mv, err := models.QueryMovie(db, id)
		if err != nil {
			fmt.Println("Movie not found on database")
		} else {
			MovieVars = models.Movie{
				ID:      mv.ID,
				Title:   mv.Title,
				Year:    mv.Year,
				Ratings: mv.Ratings,
				Likes:   mv.Likes,
			}

			// Add movie to redis
			models.AddMovie(mv)
		}

		// http.NotFound(w, r)
		// return
	} else if err != nil {

		http.Error(w, http.StatusText(500), 500)
		return
	} else {
		MovieVars = models.Movie{
			ID:      id,
			Title:   movie.Title,
			Year:    movie.Year,
			Ratings: movie.Ratings,
			Likes:   movie.Likes,
		}
	}

	t, err := template.ParseFiles("movies.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, MovieVars)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
