package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/yigitsadic/minigame/graph/generated"
	"github.com/yigitsadic/minigame/graph/model"
)

func (r *queryResolver) CurrentGame(ctx context.Context) (*model.Game, error) {
	return &model.Game{
		ID:           r.Game.Id,
		CurrentPrize: r.Game.CurrentPrize,
		CreatedAt:    r.Game.CreatedAt.Format(time.RFC3339),

		LastWinnerCheck: r.Game.LastWinnerCheck.Format(time.RFC3339),
		NextWinnerCheck: r.Game.NextWinnerCheck.Format(time.RFC3339),
	}, nil
}

func (r *subscriptionResolver) JoinGame(ctx context.Context, gameID string, identifier string) (<-chan *model.Message, error) {
	c, err := r.Game.JoinPlayer(gameID, identifier)
	if err != nil {
		return nil, err
	}

	go r.Game.PublishClaimedNumber(identifier)

	return c, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
