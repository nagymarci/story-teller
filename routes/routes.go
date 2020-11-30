package routes

import (
	"net/http"

	"github.com/nagymarci/story-teller/handlers"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
	"github.com/nagymarci/story-teller/controllers"
)

func Route(controller *controllers.StoryTeller) http.Handler {
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	story := mux.NewRouter().PathPrefix("/story").Subrouter()
	handlers.StoryTellerCreateHandler(story, controller)
	handlers.StoryTellerGetHandler(story, controller)
	handlers.StoryTellerUseHandler(story, controller)

	router.PathPrefix("/story").Handler(story)

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)

	return n
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
