package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/yigitsadic/minigame/graph"
	"github.com/yigitsadic/minigame/graph/generated"
	"github.com/yigitsadic/minigame/internal/game"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	g := game.NewGame()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Game: g}}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")

	go g.HandleGameTicker()

	if err := http.ListenAndServe(":8080", r); err != nil {
		g.CloseAllChannels()

		log.Fatal(err)
	}
}
