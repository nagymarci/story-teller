package itest

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/nagymarci/story-teller/controllers"
	"github.com/nagymarci/story-teller/events"
	"github.com/nagymarci/story-teller/model"
	"github.com/nagymarci/story-teller/routes"
	"github.com/nagymarci/story-teller/store"
)

var router http.Handler
var s *store.Default

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	s = store.New()
	sub := events.New()
	controller := controllers.New(s, sub)

	router = routes.Route(controller, sub, s)

	res := m.Run()

	os.Exit(res)
}

func TestStorytellerCreateHandler(t *testing.T) {
	t.Run("sends 201Created with new game state", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/story", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected [%d], got [%d]", http.StatusCreated, res.StatusCode)
		}

		var result model.Game
		json.NewDecoder(res.Body).Decode(&result)

		_, err := s.Load(result.Id)

		if err != nil {
			t.Fatalf("failed to read game from store [%v]", err)
		}

		if len(result.Emojis) != 2 {
			t.Fatalf("expected 2 emoji, got [%d]", len(result.Emojis))
		}

		for i, emoji := range result.Emojis {
			if i != emoji.ID {
				t.Fatalf("expected id [%d], got [%d]", i, emoji.ID)
			}
			if emoji.IsUsed {
				t.Fatalf("expected emoji not used")
			}
		}

		if len(result.Story) != 0 {
			t.Fatalf("expected no story, got [%d]", len(result.Story))
		}
	})
}

func TestStorytellerUseHandler(t *testing.T) {
	t.Run("sends 200Ok with new game state", func(t *testing.T) {

		game := model.New(2)
		s.Save(game.Id, game)

		req := httptest.NewRequest(http.MethodPost, "/story/"+game.Id+"/0", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		var result controllers.UseResponse
		json.NewDecoder(res.Body).Decode(&result)

		if len(result.Emojis) != 2 {
			t.Fatalf("expected 2 emoji, got [%d]", len(result.Emojis))
		}

		for i, emoji := range result.Emojis {
			if i != emoji.ID {
				t.Fatalf("expected id [%d], got [%d]", i, emoji.ID)
			}
			if i != 0 && emoji.IsUsed {
				t.Fatalf("expected emoji not used")
			}
		}

		if len(result.Story) != 1 {
			t.Fatalf("expected no story, got [%d]", len(result.Story))
		}

		if result.Story[0] != result.Emojis[0].Symbol {
			t.Fatalf("expected [%v], got [%v]", result.Emojis[0].Symbol, result.Story[0])
		}
	})
	t.Run("stores new game state", func(t *testing.T) {

		game := model.New(2)
		s.Save(game.Id, game)

		req := httptest.NewRequest(http.MethodPost, "/story/"+game.Id+"/0", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		result, err := s.Load(game.Id)

		if err != nil {
			t.Fatalf("failed to read game [%s] from store [%v]", game.Id, err)
		}

		if len(result.Emojis) != 2 {
			t.Fatalf("expected 2 emoji, got [%d]", len(result.Emojis))
		}

		for i, emoji := range result.Emojis {
			if i != emoji.ID {
				t.Fatalf("expected id [%d], got [%d]", i, emoji.ID)
			}
			if i != 0 && emoji.IsUsed {
				t.Fatalf("expected emoji not used")
			}
		}

		if len(result.Story) != 1 {
			t.Fatalf("expected no story, got [%d]", len(result.Story))
		}

		if result.Story[0] != result.Emojis[0].Symbol {
			t.Fatalf("expected [%v], got [%v]", result.Emojis[0].Symbol, result.Story[0])
		}
	})
}
