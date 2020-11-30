package controllers

import (
	"github.com/nagymarci/story-teller/service"
)

type StoryTeller struct {
	store map[string]*service.Game
}

func New() *StoryTeller {
	return &StoryTeller{
		store: map[string]*service.Game{},
	}
}

func (st *StoryTeller) NewGame(emojiCount int) *service.Game {
	game := service.New(emojiCount)
	st.store[game.Id] = game
	return game
}

func (st *StoryTeller) Use(gameID string, emojiID int) *service.Game {
	st.store[gameID].Use(emojiID)
	return st.store[gameID]
}

func (st *StoryTeller) Get(gameID string) *service.Game {
	return st.store[gameID]
}
