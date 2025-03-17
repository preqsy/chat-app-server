package core

import (
	datastore "chat_app_server/database"
	"chat_app_server/external"
)

type Service struct {
	datastore    datastore.Datastore
	neo4jService *external.NEO4JService
}

func CoreService(datastore datastore.Datastore, neo4jService *external.NEO4JService) *Service {
	return &Service{
		datastore:    datastore,
		neo4jService: neo4jService,
	}
}
