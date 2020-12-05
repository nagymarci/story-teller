package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nagymarci/story-teller/controllers"
	"github.com/nagymarci/story-teller/routes"
	"github.com/nagymarci/story-teller/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s := store.New()
	c := controllers.New(s)

	router := routes.Route(c)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
