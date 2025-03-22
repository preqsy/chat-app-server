package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	auth "chat_app_server/core"
	"chat_app_server/external"
	"chat_app_server/jwt_utils"

	"github.com/sirupsen/logrus"
)

type Resolver struct {
	service       *auth.Service
	jwt_utils     *jwt_utils.JWTUtils
	redis_service *external.RedisService
	// neo4jService  *external.NEO4JService
	logger *logrus.Logger
}

func NewResolver(service *auth.Service, jwt_utils *jwt_utils.JWTUtils, redis_service *external.RedisService, logger *logrus.Logger) *Resolver {
	return &Resolver{
		service:       service,
		jwt_utils:     jwt_utils,
		redis_service: redis_service,
		// neo4jService:  neo4jService,
		logger: logger,
	}
}
