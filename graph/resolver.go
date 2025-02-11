package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import "chat_app_server/core"

type Resolver struct {
	service *auth.Service
}

func NewResolver(service *auth.Service) *Resolver {
	return &Resolver{
		service: service,
	}
}
