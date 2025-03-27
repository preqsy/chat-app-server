package core

import (
	datastore "chat_app_server/database"
	"chat_app_server/external"

	"github.com/sirupsen/logrus"
)

type Service struct {
	datastore    datastore.Datastore
	neo4jService *external.NEO4JService
	logger       *logrus.Logger
}

func CoreService(datastore datastore.Datastore, neo4jService *external.NEO4JService, logger *logrus.Logger) *Service {
	return &Service{
		datastore:    datastore,
		neo4jService: neo4jService,
		logger:       logger,
	}
}
