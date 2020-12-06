package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nagymarci/story-teller/controllers"
	"github.com/sirupsen/logrus"
)

func StoryTellerCreateHandler(router *mux.Router, storyTeller *controllers.StoryTeller) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		result := storyTeller.NewGame(emojiCount())

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}).Methods(http.MethodOptions, http.MethodPost)
}

func StoryTellerUseHandler(router *mux.Router, storyTeller *controllers.StoryTeller) {
	router.HandleFunc("/{gameId}/{emojiId}", func(w http.ResponseWriter, r *http.Request) {
		gameID := mux.Vars(r)["gameId"]
		emojiID, err := strconv.Atoi(mux.Vars(r)["emojiId"])

		if err != nil {
			logrus.WithField("emojiID", mux.Vars(r)["emojiId"]).Errorln(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid ID")
			return
		}

		result, err := storyTeller.Use(gameID, emojiID)

		if err != nil {
			logrus.WithFields(logrus.Fields{"gameId": gameID, "emojiID": emojiID}).Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}).Methods(http.MethodPost, http.MethodOptions)
}

func StoryTellerGetHandler(router *mux.Router, storyTeller *controllers.StoryTeller) {
	router.HandleFunc("/{gameId}", func(w http.ResponseWriter, r *http.Request) {
		gameID := mux.Vars(r)["gameId"]

		log := logrus.WithField("gameId", gameID)
		result, err := storyTeller.Get(gameID)

		if err != nil {
			logrus.WithFields(logrus.Fields{"gameId": gameID}).Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		log.Infoln(result)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	})
}

func emojiCount() int {
	env := os.Getenv("EMOJI_COUNT")
	emojiCount, err := strconv.Atoi(env)
	if env == "" || err != nil {
		return 9
	}
	return emojiCount
}
