package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/yigitsadic/minigame/graph"
	"github.com/yigitsadic/minigame/graph/generated"
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal"
	"github.com/yigitsadic/minigame/internal/game"
	"github.com/yigitsadic/minigame/internal/random_generator"
	"log"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	g := &game.Game{
		Id:             uuid.NewString(),
		CreatedAt:      time.Now(),
		Players:        make(map[string]int),
		PlayerChannels: make(map[string]chan *model.Message),
		WinnerNumber:   random_generator.GenerateRandomNumber(),

		CurrentPrize:    internal.StartingPrize,
		NextWinnerCheck: time.Now().Add(time.Minute * internal.TryMinute),
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Game: g,
	}}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	if err := http.ListenAndServe(":8080", r); err != nil {
		for _, c := range g.PlayerChannels {
			close(c)
		}

		log.Fatal(err)
	}
}
