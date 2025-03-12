package external

import (
	"context"

	"chat_app_server/config"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
)

type NEO4JService struct {
	driver neo4j.DriverWithContext
	logger *logrus.Logger
	ctx    context.Context
}

func InitNEO4J(ctx context.Context, logger *logrus.Logger, secrets config.Secrets) (*NEO4JService, error) {

	logger.Info("Connecting to NEO4J...")

	dbUri := secrets.Neo4jUri
	dbUser := secrets.Neo4jUser
	dbPassword := secrets.Password

	driver, err := neo4j.NewDriverWithContext(
		dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""),
	)

	if err != nil {
		logger.Errorf("Failed to create Neo4j driver: %v", err)
		return nil, err
	}
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		logger.Errorf("Neo4j connection failed: %v", err)
		return nil, err
	}
	logger.Info("Database connection established")

	return &NEO4JService{driver: driver, logger: logger, ctx: ctx}, nil
}

func (n *NEO4JService) CloseNEO4J(ctx context.Context) error {
	if n.driver != nil {
		if err := n.driver.Close(ctx); err != nil {
			n.driver.Close(ctx)
			n.logger.Errorf("Error closing Neo4j connection: %v", err)
			return err

		}
		n.logger.Info("Neo4j connection closed successfully")
	}
	return nil
}
