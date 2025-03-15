package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"chat_app_server/graph/model"
	models "chat_app_server/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

	user, err := r.jwt_utils.GetCurrentAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	newMessage := models.Message{
		SenderID:   user.ID,
		ReceiverID: uint(input.ReceiverID),
		Content:    input.Content,
	}
	fmt.Println("This is the first message", newMessage)

	msgJson, err := json.Marshal(newMessage)
	if err != nil {
		return nil, err
	}

	receiverIdString := strconv.Itoa(int(input.ReceiverID))
	channel := "chat:" + receiverIdString
	fmt.Println("This is the channel", channel)
	err = r.redis_service.PublishMessage(channel, string(msgJson))

	if err != nil {
		return nil, err
	}

	message, err := r.service.SaveMessage(ctx, &newMessage)

	if err != nil {
		return nil, err
	}
	messageResponse := &model.MessageResponse{
		Content:    message.Content,
		SenderID:   int32(message.SenderID),
		ReceiverID: int32(message.ReceiverID),
		CreatedAt:  message.CreatedAt.String(),
		ID:         int32(message.ID),
	}
	return messageResponse, nil
}

// SendFriendRequest is the resolver for the sendFriendRequest field.
func (r *mutationResolver) SendFriendRequest(ctx context.Context, receiverID int32) (*model.AuthUser, error) {
	authUser, err := r.jwt_utils.GetCurrentAuthUser(ctx)
	if err != nil {
		return nil, err
	}
	r.neo4jService.CreateUser(ctx, authUser)
	fmt.Println("This is the authUser", authUser)
	return &model.AuthUser{Email: authUser.Email}, nil
}

// GetCurrentUser is the resolver for the getCurrentUser field.
func (r *queryResolver) GetCurrentUser(ctx context.Context, token string) (*model.AuthUser, error) {
	user, err := r.jwt_utils.GetCurrentAuthUser(ctx)
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

// NewMessage is the resolver for the newMessage field.
func (r *subscriptionResolver) NewMessage(ctx context.Context, receiverID int32) (<-chan *model.MessageResponse, error) {
	msgChan := make(chan *model.MessageResponse, 1)
	channel := "chat:" + strconv.Itoa(int(receiverID))

	pubSub := r.redis_service.SubscribeToChannel(channel)
	go func() {
		defer pubSub.Close()
		for msg := range pubSub.Channel() {

			var message model.MessageResponse
			err := json.Unmarshal([]byte(msg.Payload), &message)
			fmt.Println("Subsscription Message", message)
			if err == nil {
				msgChan <- &message
			}
		}
	}()
	go func() {
		<-ctx.Done()
		pubSub.Close()
		close(msgChan)
	}()
	return msgChan, nil
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
