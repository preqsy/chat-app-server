# Chat App Server

A real-time chat application server built with Go, GraphQL, PostgreSQL, Redis, and Neo4j. This server provides a complete backend solution for a chat application with user authentication, friend requests, real-time messaging, and message history.

## 🚀 Features

- **User Authentication**: JWT-based authentication with secure login/registration
- **Real-time Messaging**: WebSocket-based GraphQL subscriptions for instant messaging
- **Friend System**: Send, accept, and manage friend requests
- **Message History**: Retrieve conversation history between users
- **Recent Chats**: Get recent chat conversations for a user
- **Multi-database Architecture**: 
  - PostgreSQL for user data and messages
  - Redis for caching and session management
  - Neo4j for social graph and friend relationships
- **GraphQL API**: Single endpoint with powerful querying capabilities
- **CORS Support**: Configured for web client integration

## 🛠️ Tech Stack

- **Language**: Go 1.23.2
- **API**: GraphQL with gqlgen
- **Databases**: 
  - PostgreSQL (primary data storage)
  - Redis (caching/sessions)
  - Neo4j (social graph)
- **Authentication**: JWT tokens
- **Real-time**: WebSocket subscriptions
- **Containerization**: Docker

## 📋 Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Redis server
- Neo4j database
- Docker (optional)

## ⚙️ Installation

### 1. Clone the repository
```bash
git clone https://github.com/preqsy/chat-app-server.git
cd chat-app-server
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Environment Configuration
Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_PORT=5432
DB_USER=your_db_user
PASSWORD=your_db_password
HOST=localhost
DBNAME=chat_app_server

# Server Configuration
PORT=8000
JWT_SECRET=your_jwt_secret

# Redis Configuration
REDIS_URL=redis://localhost:6379

# Neo4j Configuration
NEO4J_URI=bolt://localhost:7687
NEO4J_PASSWORD=your_neo4j_password
NEO4J_USER=neo4j
```

### 4. Database Setup

#### PostgreSQL
Ensure your PostgreSQL database is running and create the required database:
```sql
CREATE DATABASE chat_app_server;
```

#### Redis
Start your Redis server:
```bash
redis-server
```

#### Neo4j
Start your Neo4j instance and ensure it's accessible at the configured URI.

### 5. Run the application

#### Development
```bash
go run main.go
```

#### Production Build
```bash
go build -o app .
./app
```

#### Docker
```bash
docker build -t chat-app-server .
docker run -p 8000:8000 --env-file .env chat-app-server
```

## 🔧 API Usage

The server runs on `http://localhost:8000` by default.

### GraphQL Playground
Access the GraphQL playground at: `http://localhost:8000/`

### GraphQL Endpoint
Send queries and mutations to: `http://localhost:8000/query`

## 📝 GraphQL Schema

### User Authentication

#### Register a new user
```graphql
mutation {
  createAuthUser(input: {
    username: "johndoe"
    firstName: "John"
    lastName: "Doe"
    email: "john@example.com"
    password: "securepassword"
  }) {
    authUser {
      id
      username
      email
    }
    token
  }
}
```

#### Login
```graphql
mutation {
  loginAuthUser(input: {
    email: "john@example.com"
    password: "securepassword"
  }) {
    token
  }
}
```

### Messaging

#### Send a message
```graphql
mutation {
  sendMessage(input: {
    sender_id: "1"
    receiver_id: 2
    content: "Hello there!"
  }) {
    id
    content
    createdAt
  }
}
```

#### Get messages between users
```graphql
query {
  retrieveMessages(sender_id: 1, receiver_id: 2) {
    id
    content
    createdAt
  }
}
```

#### Subscribe to new messages
```graphql
subscription {
  newMessage(receiver_id: 1) {
    id
    content
    sender {
      username
    }
    createdAt
  }
}
```

### Friends

#### Send friend request
```graphql
mutation {
  sendFriendRequest(receiver_id: 2) {
    id
    username
  }
}
```

#### Accept friend request
```graphql
mutation {
  acceptFriendRequest(sender_id: 1) {
    id
    username
  }
}
```

#### List friends
```graphql
query {
  listFriends(filters: { skip: 0, limit: 10 }) {
    id
    username
    firstName
    lastName
  }
}
```

## 🏗️ Project Structure

```
├── config/             # Configuration management
├── core/              # Business logic layer
├── database/          # Database connections and CRUD operations
├── external/          # External service integrations (Redis, Neo4j)
├── graph/             # GraphQL schema and resolvers
├── jwt_utils/         # JWT token utilities
├── middleware/        # HTTP middleware (auth, CORS)
├── model/             # Data models
├── utils/             # Utility functions
├── main.go            # Application entry point
├── Dockerfile         # Docker configuration
└── .env               # Environment variables
```

## 🔐 Authentication

The server uses JWT (JSON Web Tokens) for authentication. Include the token in your requests:

```http
Authorization: Bearer <your_jwt_token>
```

## 🚦 Health Check

The server provides a GraphQL playground at the root endpoint (`/`) which can be used to verify the server is running correctly.

## 🐳 Docker Deployment

Build and run with Docker:

```bash
# Build the image
docker build -t chat-app-server .

# Run the container
docker run -d \
  --name chat-app \
  -p 8000:8000 \
  --env-file .env \
  chat-app-server
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Issues

If you encounter any issues or have questions, please file an issue on the [GitHub Issues](https://github.com/preqsy/chat-app-server/issues) page.

## 📞 Support

For support and questions, please reach out through:
- GitHub Issues
- Email: [your-email@example.com]

---

Built with ❤️ using Go and GraphQL
