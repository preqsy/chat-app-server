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
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"

	// "github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"

	// "github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	secret := config.GetSecrets()

	datastore, err := database.ConnectDB(secret.Host, secret.Db_User, secret.Password, secret.DbName, secret.Port)
	if err != nil {
		log.Fatal(err)
	}

	coreService := core.CoreService(datastore)
	jwtService := jwt_utils.InitDB(datastore)
	ctx := context.Background()

	redisService, err := external.InitRedis(ctx)
	if err != nil {
		logrus.Error("Redis connection failed:", err)
	}
	_, err = external.InitNEO4J(ctx, &logrus.Logger{})
	if err != nil {
		logrus.Error("NEO4J connection failed", err)
	}
	resolver := graph.NewResolver(coreService, jwtService, redisService)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
				// origin := r.Header.Get("Origin")
				// if origin == "" || origin == r.Header.Get("Host") {
				// 	return true
				// }
				// log.Printf("WebSocket connection attempt from origin: %s", origin)

				// return slices.Contains([]string{":5173", "http://localhost:5173"}, origin)

			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
