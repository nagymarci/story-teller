package controllers

import (
	"github.com/nagymarci/story-teller/model"
	"github.com/nagymarci/story-teller/store"
)

type StoryTeller struct {
	store *store.Default
}

func New(store *store.Default) *StoryTeller {
	return &StoryTeller{
		store: store,
	}
}

func (st *StoryTeller) NewGame(emojiCount int) *model.Game {
	game := model.New(emojiCount)
	st.store.Save(game.Id, game)
	return game
}

func (st *StoryTeller) Use(gameID string, emojiID int) (*model.Game, error) {
	game, err := st.store.Load(gameID)
	if err != nil {
		return game, err
	}
	game.Use(emojiID)
	st.store.Save(gameID, game)
	return game, nil
}

func (st *StoryTeller) Get(gameID string) (*model.Game, error) {
	return st.store.Load(gameID)
}
