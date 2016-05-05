package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jameycribbs/cribbnotes/global_vars"
	"github.com/jameycribbs/cribbnotes/handlers/notes_handler"
	"github.com/justinas/nosurf"
	"net/http"
)

func main() {
	var port string

	port = ":8080"

	store := sessions.NewCookieStore([]byte("cribbnotes-is-awesome"))

	gv := global_vars.GlobalVars{SessionStore: store}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()
	r.HandleFunc("/", makeHandler(notes_handler.Index, &gv)).Methods("GET")
	r.HandleFunc("/search", makeHandler(notes_handler.Index, &gv)).Methods("POST")
	r.HandleFunc("/new", makeHandler(notes_handler.New, &gv)).Methods("GET")
	r.HandleFunc("/create", makeHandler(notes_handler.Create, &gv)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/edit", makeHandler(notes_handler.Edit, &gv)).Methods("GET")
	r.HandleFunc("/update", makeHandler(notes_handler.Update, &gv)).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}/delete", makeHandler(notes_handler.Delete, &gv)).Methods("GET")
	r.HandleFunc("/destroy", makeHandler(notes_handler.Destroy, &gv)).Methods("POST")

	http.Handle("/", r)

	csrfHandler := nosurf.New(http.DefaultServeMux)

	csrfHandler.SetFailureHandler(http.HandlerFunc(failHand))

	http.ListenAndServe(port, csrfHandler)
}

func failHand(w http.ResponseWriter, r *http.Request) {
	// will return the reason of the failure
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string, *global_vars.GlobalVars),
	gv *global_vars.GlobalVars) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		fn(w, r, vars["id"], gv)
	}
}
