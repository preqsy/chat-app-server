package external

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
)

type NEO4JService struct {
	client neo4j.DriverWithContext
	logger *logrus.Logger
	ctx    context.Context
}

func InitNEO4J(ctx context.Context, logger *logrus.Logger) (*NEO4JService, error) {
	fmt.Println("Connecting to NEO4J...")
	dbUri := "neo4j://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "50610903"
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

	return &NEO4JService{client: driver, logger: logger, ctx: ctx}, nil
}
