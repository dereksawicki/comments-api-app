/*
Derek Sawicki 2017
app.go

App stucture is essentially an interface for global utilities.
It maintains accesss to the router and the database

*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// verify connection
	if err := a.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

//  initialize routes for API
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/comments/{page:[0-9]+}", a.getComments).Methods("GET")
	a.Router.HandleFunc("/comment", a.createComment).Methods("POST")
	a.Router.HandleFunc("/comment/{id:[0-9]+}", a.getComment).Methods("GET")
	a.Router.HandleFunc("/comment/{id:[0-9]+}", a.updateComment).Methods("PUT")
	a.Router.HandleFunc("/comment/{id:[0-9]+}", a.deleteComment).Methods("DELETE")
}

// start listening on port @addr
func (a *App) Run(addr string) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	headers := []string{"Content-Type"}
	handler := handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers))(a.Router)
	log.Fatal(http.ListenAndServe(addr, handler))
}

/* HANDLERS */
func (a *App) getComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // retrieve route variables

	pagenum, err := strconv.Atoi(vars["page"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid page number")
		return
	}

	if pagenum < 0 {
		pagenum = 0
	}

	/* retrieve comments on page pagenum */
	comments, err := getComments(a.DB, (pagenum-1)*10, 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, comments)
}

func (a *App) getComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}

	if id < 0 {
		id = 0
	}

	c := comment{ID: id}
	if err := c.getComment(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Comment not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id number")
		return
	}

	if id < 0 {
		id = 0
	}

	c := comment{ID: id}
	if err := c.deleteComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) createComment(w http.ResponseWriter, r *http.Request) {
	var c comment

	fmt.Println("creating comment")

	// extract comment from http request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		fmt.Println("Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// insert into database
	if err := c.createComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) updateComment(w http.ResponseWriter, r *http.Request) {
	// extract id from request variables
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// extract comment from http request load
	var c comment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	c.ID = id

	if err := c.updateComment(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
