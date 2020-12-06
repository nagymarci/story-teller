package controllers

import (
	"github.com/nagymarci/story-teller/model"
)

type UseResponse struct {
	Emojis []*model.Emoji    `json:"emojis"`
	Story  []model.StoryItem `json:"story"`
}

func NewUseResponse(g *model.Game) *UseResponse {
	return &UseResponse{
		Emojis: g.Emojis,
		Story:  g.Story,
	}
}
