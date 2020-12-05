package events

import (
	"errors"
	"sync"
)

type game struct {
	sync.Mutex
	clients map[interface{}]chan interface{}
}

func newGame() *game {
	return &game{
		clients: map[interface{}]chan interface{}{},
	}
}

type InApp struct {
	sync.Mutex
	games map[string]*game
}

func New() *InApp {
	return &InApp{
		games: map[string]*game{},
	}
}

func (b *InApp) Subscribe(gameID string, clientID interface{}) (chan interface{}, error) {
	c := make(chan interface{})

	var g *game

	g, ok := b.games[gameID]
	if !ok {
		b.Lock()
		defer b.Unlock()

		g = newGame()
		b.games[gameID] = g
	}

	g.Lock()
	defer g.Unlock()

	g.clients[clientID] = c

	return c, nil
}

func (b *InApp) Unsubscribe(gameID string, clientID interface{}) error {
	g, ok := b.games[gameID]
	if !ok {
		return errors.New("no game found")
	}

	g.Lock()
	defer g.Unlock()

	if c, ok := g.clients[clientID]; ok {
		close(c)
		delete(g.clients, clientID)
	}

	if len(g.clients) == 0 {
		delete(b.games, gameID)
	}

	return nil
}

// Type tells which kind of events happened
type Type string

// Available types
const (
	Use Type = "use"
)

type Event struct {
	Action Type
	Data   interface{}
}

func (b *InApp) Emit(gameID string, t Type, body interface{}) {
	g, ok := b.games[gameID]
	if !ok {
		return
	}

	g.Lock()
	defer g.Unlock()

	for _, s := range g.clients {
		s <- Event{t, body}
	}
}
