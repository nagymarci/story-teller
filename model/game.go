package model

import (
	"math/rand"
)

var defaultEmojiSet = [...]string{"🐒", "🐕", "🐕", "🦌", "🐘", "🤙", "🧠", "🕵🏼‍♂️"}

type Player struct {
	name string
}

type Game struct {
	Id     string   `json:"id"`
	Emojis []Emoji  `json:"emojis"`
	Story  []string `json:"story"`
}

type Emoji struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	IsUsed bool   `json:"isUsed"`
}

func New(emojiCount int) *Game {
	return &Game{
		Id:     generateID(4),
		Emojis: generateEmojis(emojiCount),
		Story:  []string{},
	}
}

func (g *Game) Use(id int) {
	if g.Emojis[id].IsUsed {
		return
	}

	g.Emojis[id].IsUsed = true
	g.Story = append(g.Story, g.Emojis[id].Symbol)
}

func generateEmojis(length int) []Emoji {
	var emojis []Emoji

	for i := 0; i < length; i++ {
		emojis = append(emojis, Emoji{ID: i, Symbol: defaultEmojiSet[rand.Intn(len(defaultEmojiSet))], IsUsed: false})
	}

	return emojis
}

func generateID(length int) string {
	const idCharset = "abcdefghijklmnopqrstvwxyz0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = idCharset[rand.Intn(len(idCharset))]
	}
	return string(b)
}
