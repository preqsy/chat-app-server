package core

import datastore "chat_app_server/database"

type Service struct {
	datastore datastore.Datastore
}

func CoreService(datastore datastore.Datastore) *Service {
	return &Service{
		datastore: datastore,
	}
}
