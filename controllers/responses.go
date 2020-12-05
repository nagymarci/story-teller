package controllers

import (
	"github.com/nagymarci/story-teller/model"
)

type UseResponse struct {
	Emojis []*model.Emoji
	Story  []model.StoryItem
}

func NewUseResponse(g *model.Game) *UseResponse {
	return &UseResponse{
		Emojis: g.Emojis,
		Story:  g.Story,
	}
}
