package controllers

import (
	"github.com/nagymarci/story-teller/events"
	"github.com/nagymarci/story-teller/model"
	"github.com/nagymarci/story-teller/store"
)

type StoryTeller struct {
	store  *store.Default
	events *events.InApp
}

func New(store *store.Default, sub *events.InApp) *StoryTeller {
	return &StoryTeller{
		store:  store,
		events: sub,
	}
}

func (st *StoryTeller) NewGame(emojiCount int) *model.Game {
	game := model.New(emojiCount)
	st.store.Save(game.Id, game)
	return game
}

func (st *StoryTeller) Use(gameID string, emojiID int) (*UseResponse, error) {
	game, err := st.store.Load(gameID)
	if err != nil {
		return nil, err
	}
	game.Use(emojiID)
	st.store.Save(gameID, game)

	res := NewUseResponse(game)

	st.events.Emit(gameID, events.Use, res)
	return res, nil
}

func (st *StoryTeller) Get(gameID string) (*model.Game, error) {
	return st.store.Load(gameID)
}
