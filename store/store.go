package store

import (
	"errors"

	"github.com/nagymarci/story-teller/model"
)

type Default struct {
	repo map[string]*model.Game
}

func New() *Default {
	return &Default{
		repo: map[string]*model.Game{},
	}
}

func (s *Default) Save(id string, game *model.Game) {
	s.repo[id] = game
}

func (s *Default) Load(id string) (*model.Game, error) {
	game, ok := s.repo[id]
	if !ok {
		return game, errors.New("not exists")
	}

	return game, nil
}
