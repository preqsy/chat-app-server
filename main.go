package main

import (
	"chat_app_server/config"
	"chat_app_server/core"
	database "chat_app_server/database/crud"
	"chat_app_server/external"
	"chat_app_server/graph"
	"chat_app_server/jwt_utils"
	"chat_app_server/middleware"
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
)

// const defaultPort = "8080"

func main() {
	secrets := config.GetSecrets()
	logger := logrus.New()
	// logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	datastore, err := database.ConnectDB(secrets.Host, secrets.Db_User, secrets.Password, secrets.DbName, secrets.DbPort)
	if err != nil {
		logger.Fatal(err)
	}

	jwtService := jwt_utils.InitializeJWTUtils(datastore, logger)
	ctx := context.Background()

	redisService, err := external.InitRedis(ctx, logger, secrets.RedisURL)
	if err != nil {
		logrus.Error("Redis connection failed:", err)
	}
	neo4jService, err := external.InitNEO4J(ctx, logger, secrets)
	if err != nil {
		logrus.Error("NEO4J connection failed", err)
	}
	defer neo4jService.CloseNEO4J(ctx)

	coreService := core.CoreService(datastore, neo4jService, logger)

	resolver := graph.NewResolver(coreService, jwtService, redisService, logger)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	queryHandler := middleware.AuthMiddleware(srv)

	http.Handle("/", corsHandler(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", corsHandler(queryHandler))

	logger.Printf("connect to http://0.0.0.0:%s/ for GraphQL playground", secrets.DefaultPort)
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+secrets.DefaultPort, nil))
}
