package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"chat_app_server/graph/model"
	models "chat_app_server/model"
	"context"
	"fmt"
	"net/http"
	"time"
)

// CreateAuthUser is the resolver for the createAuthUser field.
func (r *mutationResolver) CreateAuthUser(ctx context.Context, input model.AuthUserCreate) (*model.AuthUserResponse, error) {
	newUser := models.AuthUser{
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
	}
	savedUser, err := r.service.SaveUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	return &model.AuthUserResponse{
		AuthUser: &model.AuthUser{
			ID:        int32(savedUser.AuthUser.ID),
			Email:     savedUser.AuthUser.Email,
			FirstName: savedUser.AuthUser.FirstName,
			LastName:  savedUser.AuthUser.LastName,
			CreatedAt: savedUser.AuthUser.CreatedAt.Format(time.RFC3339),
		},
		Token: savedUser.Token,
	}, nil
}

// LoginAuthUser is the resolver for the loginAuthUser field.
func (r *mutationResolver) LoginAuthUser(ctx context.Context, input model.AuthUserLogin) (*model.LoginResponse, error) {
	user := models.AuthUserLogin{
		Email:    input.Email,
		Password: input.Password,
	}
	token, err := r.service.LoginUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token}, nil
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.MessageInput) (*model.MessageResponse, error) {
	request, ok := ctx.Value("request").(*http.Request)
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	token := request.Header.Get("authorization")

	user, err := r.jwt_utils.GetCurrentAuthUser(ctx, token)
	if err != nil {
		return nil, err
	}
	fmt.Print(user)
	newMessage := models.Message{
		SenderID:   user.ID,
		ReceiverID: uint(input.Receiver),
		Content:    input.Content,
	}
	message, err := r.service.SaveMessage(ctx, &newMessage)

	if err != nil {
		return nil, err
	}
	messageResponse := &model.MessageResponse{
		Content:   message.Content,
		Sender:    int32(message.SenderID),
		Receiver:  int32(message.ReceiverID),
		CreatedAt: message.CreatedAt.String(),
		ID:        int32(message.ID),
	}
	return messageResponse, nil
}

// GetCurrentUser is the resolver for the getCurrentUser field.
func (r *queryResolver) GetCurrentUser(ctx context.Context, token string) (*model.AuthUser, error) {
	user, err := r.jwt_utils.GetCurrentAuthUser(ctx, token)
	if err != nil {
		return nil, err
	}
	return &model.AuthUser{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		ID:        int32(user.ID),
	}, nil
}

// CurrentTime is the resolver for the currentTime field.
func (r *subscriptionResolver) CurrentTime(ctx context.Context) (<-chan *model.Time, error) {
	ch := make(chan *model.Time)

	go func() {
		defer close(ch)

		count := 1
		for {
			time.Sleep(1)
			count = count + 1
			fmt.Println("Tick", count)

			currentTime := time.Now()

			t := &model.Time{
				UnixTime:  int32(currentTime.Unix()),
				TimeStamp: currentTime.Format(time.RFC3339),
			}

			select {
			case <-ctx.Done():
				fmt.Println("Subscription Closed at", count)
				return
			case ch <- t:

			}

		}
	}()
	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
