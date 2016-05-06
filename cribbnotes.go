package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jameycribbs/cribbnotes/handlers/notes_handler"
	"github.com/justinas/nosurf"
	"io/ioutil"
	"net/http"
)

type Config struct {
	DataDir string `json:"data_dir"`
	Port    string `json:"port"`
}

func main() {
	var config Config

	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Cannot read config file!!!")
		return
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Invalid config file!!!")
		return
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()
	r.HandleFunc("/", makeHandler(notes_handler.Index, config.DataDir)).Methods("GET")
	r.HandleFunc("/search", makeHandler(notes_handler.Index, config.DataDir)).Methods("POST")
	r.HandleFunc("/new", makeHandler(notes_handler.New, config.DataDir)).Methods("GET")
	r.HandleFunc("/create", makeHandler(notes_handler.Create, config.DataDir)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/edit", makeHandler(notes_handler.Edit, config.DataDir)).Methods("GET")
	r.HandleFunc("/update", makeHandler(notes_handler.Update, config.DataDir)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/delete", makeHandler(notes_handler.Delete, config.DataDir)).Methods("GET")
	r.HandleFunc("/destroy", makeHandler(notes_handler.Destroy, config.DataDir)).Methods("POST")

	http.Handle("/", r)

	csrfHandler := nosurf.New(http.DefaultServeMux)

	csrfHandler.SetFailureHandler(http.HandlerFunc(failHand))

	fmt.Println("CribbNotes data directory is: ", config.DataDir)
	fmt.Println("CribbNotes server is running on port: ", config.Port)

	http.ListenAndServe(":"+config.Port, csrfHandler)
}

func failHand(w http.ResponseWriter, r *http.Request) {
	// will return the reason of the failure
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string, string), dataDir string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		fn(w, r, vars["id"], dataDir)
	}
}
