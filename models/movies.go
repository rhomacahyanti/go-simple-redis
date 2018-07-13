package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"movies/connection"
	"strconv"
)

// var db *pool.Pool

// func init() {
// 	var err error

// 	db, err = pool.New("tcp", "localhost:6379", 10)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

// ErrNoMovie var
var ErrNoMovie = errors.New("models: no movie found")

//Movie struct
type Movie struct {
	ID      int
	Title   string
	Year    string
	Ratings int
	Likes   int
}

func populateMovie(reply map[string]string) *Movie {
	var err error

	mv := new(Movie)
	mv.Title = reply["title"]
	mv.Year = reply["year"]
	mv.Ratings, err = strconv.Atoi(reply["ratings"])
	if err != nil {
		return nil
	}
	mv.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil
	}

	fmt.Println(mv)

	return mv
}

// FindMovie on redis
func FindMovie(id string) (*Movie, error) {
	conn, err := connection.ConnectRedis()

	defer conn.Close()

	reply, err := conn.Cmd("HGETALL", "movie:"+id).Map()
	if err != nil {
		fmt.Println("Error found: ", err)
		return nil, err
	} else if len(reply) == 0 {
		fmt.Println("Error found: ", ErrNoMovie)
		return nil, ErrNoMovie
	}

	fmt.Println("Movie ", id, " found on redis!")

	return populateMovie(reply), err
}

// AddMovie to redis
func AddMovie(movie Movie) {
	conn, err := connection.ConnectRedis()
	checkError(err)

	defer conn.Close()

	id := strconv.Itoa(movie.ID)
	resp := conn.Cmd("HMSET", "movie:"+id, "title", movie.Title, "year", movie.Year, "ratings", movie.Ratings, "likes", movie.Likes)

	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		log.Fatal(resp.Err)
	} else {
		fmt.Println(movie.Title, "added!")
	}
}

// QueryMovie from database
func QueryMovie(db *sql.DB, id int) (Movie, error) {
	var movie Movie

	queryStatement := `SELECT * FROM movies WHERE id = ?`

	row := db.QueryRow(queryStatement, id)

	err := row.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Ratings, &movie.Likes)
	checkError(err)

	fmt.Println(movie.Title)

	return movie, err
}

// QueryMovies
func QueryMovies(db *sql.DB) []Movie {
	var movies []Movie

	rows, err := db.Query("SELECT * FROM movies")
	checkError(err)

	var movie Movie

	for rows.Next() {
		err = rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Ratings, &movie.Likes)
		checkError(err)

		fmt.Println(movie.Title)

		movies = append(movies, movie)
	}

	return movies
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
