package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yigitsadic/minigame/graph/generated"
	"github.com/yigitsadic/minigame/graph/model"
	"github.com/yigitsadic/minigame/internal"
	"github.com/yigitsadic/minigame/internal/random_generator"
)

func (r *mutationResolver) PublishMessage(ctx context.Context, message string) (*bool, error) {
	var ok bool

	for _, c := range r.Game.PlayerChannels {
		c <- &model.Message{
			ID:          uuid.NewString(),
			Text:        message,
			MessageType: model.MessageTypeDoublePrize,
		}
	}

	ok = true

	return &ok, nil
}

func (r *queryResolver) GameSession(ctx context.Context) (*model.GameSession, error) {
	return &model.GameSession{
		ID:           r.Game.Id,
		CurrentPrize: r.Game.CurrentPrize,
		CreatedAt:    r.Game.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (r *subscriptionResolver) MessageFeed(ctx context.Context, gameID string, identifier string) (<-chan *model.Message, error) {
	if gameID != r.Game.Id {
		return nil, errors.New("game is no longer valid")
	}

	if len(r.Game.Players) >= internal.PlayerLimit {
		return nil, errors.New("maximum user limit reached")
	}

	r.Game.Players[identifier] = random_generator.GenerateRandomNumber()
	r.Game.PlayerChannels[identifier] = make(chan *model.Message)

	return r.Game.PlayerChannels[identifier], nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
