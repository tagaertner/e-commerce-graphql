package main

import (
	"log"
	"net/http"
	"os"
    "flag"
    "fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tagaertner/e-commerce-graphql/services/users/generated"
	"github.com/tagaertner/e-commerce-graphql/services/users/resolvers"
	"github.com/tagaertner/e-commerce-graphql/services/users/database" 
    "github.com/joho/godotenv"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/tagaertner/e-commerce-graphql/services/users/services"
)

// const defaultPort = "4002"

func main() {
    // Only load .env file when not in Docker
    if _, err := os.Stat(".env"); err == nil {
        err := godotenv.Load()
        if err != nil {
            log.Println("⚠️  Failed to load .env file")
        } else {
            log.Println("✅ Loaded .env file")
        }
    } else {
        log.Println("📦 Running in containerized environment, using system environment variables")
    }
    
    // Flag to check the database connection and exit
	testDB := flag.Bool("test-db", false, "Test DB connection and exit")
	flag.Parse()

	db := database.Connect()

	if *testDB {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("❌ Failed to get sql DB: %v", err)
            fmt.Println("🔍 DB_HOST =", os.Getenv("DB_HOST"))
		}
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("❌ Database ping failed: %v", err)
		}
		log.Println("✅ Connected to PostgreSQL successfully")
		return 
	}

    database.RunMigrations(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4002"
	}

	// Pass db into resolver
	userService := services.NewUserService(db)

	resolver := &resolvers.Resolver{
		UserService: userService,
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	// Enable introspection 
    srv.Use(extension.Introspection{})

	// Supported transport methods for GraphQL requests:
	// - POST and GET for queries/mutations
	// - WebSocket transport enables live data features like subscriptions
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Websocket{})

	// Routes
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "service": "users"}`))
	})

	log.Printf("🛍️ [users] service ready at http://users:%s/query", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

