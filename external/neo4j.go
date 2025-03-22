package external

import (
	"context"
	"encoding/json"
	"fmt"

	"chat_app_server/config"
	models "chat_app_server/model"

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

func StructToMap(data any) (map[string]any, error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	return result, err
}

func (n *NEO4JService) CreateUserNode(ctx context.Context, payload *models.AuthUser) {
	if n.driver == nil {
		n.logger.Error("Neo4j driver is nil")
		return
	}
	resultMap, err := StructToMap(&payload)
	if err != nil {
		n.logger.Error("Error converting struct to map")
		return
	}
	delete(resultMap, "password")

	_, err = neo4j.ExecuteQuery(
		ctx, n.driver, "CREATE (u:User {name: $username, user_id: $ID, firstName: $firstName, lastName: $lastName, email: $email}) RETURN u", resultMap, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		n.logger.Errorf("Error creating user: %v", err)
	}
	n.logger.Info("User created successfully")

}

func (n *NEO4JService) SendFriendRequest(ctx context.Context, sender, receiver *models.AuthUser) (*models.AuthUser, error) {
	if n.driver == nil {
		n.logger.Error("Neo4j driver is nil")
		return nil, nil
	}
	result, err := neo4j.ExecuteQuery(
		ctx, n.driver,
		`
		MERGE (sender:User {userId: $senderId})
		ON CREATE SET sender.name = $senderName
		ON MATCH SET sender.name = sender.name
		
		MERGE (receiver:User {userId: $receiverId})
		ON CREATE SET receiver.name = $receiverName
		ON MATCH SET receiver.name = receiver.name
		
		MERGE (sender)-[:FRIEND_REQUEST]->(receiver)
		
		`,
		map[string]any{"senderId": sender.ID, "receiverId": receiver.ID, "senderName": sender.Username, "receiverName": receiver.Username},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		n.logger.Errorf("Error sending friend request: %v", err)
		return nil, err
	}

	n.logger.Info("Friend request sent successfully", result.Summary)
	return receiver, err

}

func (n *NEO4JService) AcceptFriendRequest(ctx context.Context, sender, receiver *models.AuthUser) (*models.AuthUser, error) {
	if n.driver == nil {
		n.logger.Error("Neo4j driver is nil")
		return nil, nil
	}
	result, err := neo4j.ExecuteQuery(
		ctx, n.driver,
		`
		MATCH (sender:User {userId: $senderId})-[r:FRIEND_REQUEST]->(receiver:User {userId: $receiverId})
		DELETE r
		
		MERGE (sender)-[:FRIENDS]->(receiver)
		MERGE (receiver)-[:FRIENDS]->(sender)
		`, map[string]any{"senderId": sender.ID, "receiverId": receiver.ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		n.logger.Errorf("Error accepting friend request: %v", err)
		return nil, err
	}
	n.logger.Info("Friend request accepted successfully", result.Summary)
	return receiver, err
}

func (n *NEO4JService) CheckIfFriends(ctx context.Context, sender, receiver *models.AuthUser) (bool, error) {
	if n.driver == nil {
		n.logger.Error("Neo4j driver is nil")
		return false, nil
	}

	result, err := neo4j.ExecuteQuery(
		ctx, n.driver,
		`
		MATCH (sender:User {userId: $senderId})-[r:FRIENDS]-(receiver:User {userId: $receiverId})
		RETURN COUNT(r) AS friendCount
		
		`,
		map[string]any{"senderId": sender.ID, "receiverId": receiver.ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		n.logger.Errorf("error checking friendship: %s", err)
		return false, err
	}

	if len(result.Records) == 0 {
		return false, nil // No friendship found
	}

	friendCount, _ := result.Records[0].Get("friendCount")
	isFriends := friendCount.(int64) > 0

	return isFriends, nil
}

func (n *NEO4JService) ListFriendRequests(ctx context.Context, user *models.AuthUser) ([]int64, error) {
	if n.driver == nil {
		n.logger.Errorf("error connecting to driver")
		return nil, fmt.Errorf("error connecting to neo4j driver")
	}
	result, err := neo4j.ExecuteQuery(
		ctx, n.driver,
		`
		MATCH (sender: User)-[:FRIEND_REQUEST]-> (receiver:User {userId:$userId})
		RETURN sender.userId
		`,
		map[string]any{"userId": user.ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		n.logger.Errorf("failed to retrieve friend requests: %v", err)
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, err
	}

	var users []int64

	for _, record := range result.Records {
		value, ok := record.Get("sender.userId")
		if !ok {
			return nil, err
		}
		userId := value.(int64)
		users = append(users, userId)
	}
	return users, nil
}
func (n *NEO4JService) ListFriends(ctx context.Context, user *models.AuthUser) ([]int64, error) {
	if n.driver == nil {
		n.logger.Errorf("error connecting to driver")
		return nil, fmt.Errorf("error connecting to neo4j driver")
	}
	result, err := neo4j.ExecuteQuery(
		ctx, n.driver,
		`
		MATCH (sender: User)-[:FRIENDS]-> (receiver:User {userId:$userId})
		RETURN sender.userId
		`,
		map[string]any{"userId": user.ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		n.logger.Errorf("failed to retrieve friends: %v", err)
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, err
	}

	var users []int64

	for _, record := range result.Records {
		value, ok := record.Get("sender.userId")
		if !ok {
			return nil, err
		}
		userId := value.(int64)
		users = append(users, userId)
	}
	return users, nil
}
