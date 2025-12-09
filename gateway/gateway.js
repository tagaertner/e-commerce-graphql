const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");

async function startServer() {
  try {
    console.log("🔄 Starting E-Commerce Federation Gateway...");

    // Create the federation gateway
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: "products", url: "http://products:4001/query" },
          { name: "users", url: "http://users:4002/query" },
          { name: "orders", url: "http://orders:4003/query" },
        ],
        introspectionHeaders: {
          "User-Agent": "ApolloGateway/2.5.5",
        },
      }),
      // Poll every 10 seconds for schema changes
      pollIntervalInMs: 10000,
      // Add explicit federation configuration
      buildService: ({ url }) => {
        return new RemoteGraphQLDataSource({
          url,
          willSendRequest: ({ request, context }) => {
            request.http.headers.set("apollo-federation-include-trace", "ftv1");
          },
        });
      },
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      introspection: true, // Enable introspection for development
      csrfPrevention: {
        requestHeaders: ["apollo-required-preflight"],
      },
      plugins: [
        // Enable GraphQL Playground
        ApolloServerPluginLandingPageLocalDefault({
          embed: true,
          settings: {
            "editor.theme": "dark",
            "editor.fontSize": 14,
          },
        }),

        // Custom plugin for health checks and logging
        {
          requestDidStart() {
            return {
              didResolveOperation(requestContext) {
                console.log(`📊 Query: ${requestContext.request.operationName || "Anonymous"}`);
              },
              didEncounterErrors(requestContext) {
                console.error("❌ GraphQL errors:", requestContext.errors);
              },
            };
          },
        },
      ],
      // Custom error formatting
      formatError: (error) => {
        console.error("🚨 Gateway Error:", error);
        return {
          message: error.message,
          code: error.extensions?.code,
          path: error.path,
        };
      },
    });

    // Start the server
    const { url } = await startStandaloneServer(server, {
      listen: { port: 4000 },
    });

    // Success logging
    console.log("✅ Federation Gateway Successfully Started!");
    console.log("");
    console.log(`🚀 Gateway ready at ${url}`);
    console.log(`🎮 GraphQL Playground: ${url}`);
    console.log("");
    console.log("📋 Connected Services:");
    console.log("  🛍️  Products Service: http://localhost:4001/query");
    console.log("  👥 Users Service: http://localhost:4002/query");
    console.log("  📦 Orders Service: http://localhost:4003/query");
    console.log("");
    console.log("🔗 Example Federated Query:");
    console.log(`
    query ExampleCrossServiceQuery {
      user(id: "1") {
        name
        email
        password
        # This will automatically resolve across services
      }
      products {
        name
        price
      }
      orders {
        id
        status
      }
    }`);
  } catch (error) {
    console.error("💥 Failed to start federation gateway:", error);

    // Error messages
    if (error.message.includes("ECONNREFUSED")) {
      console.error("");
      console.error("🚨 Connection Error: Make sure all services are running:");
      console.error("  - Products service on port 4001");
      console.error("  - Users service on port 4002");
      console.error("  - Orders service on port 4003");
      console.error("");
      console.error("💡 Start each service with: go run main.go");
    }

    process.exit(1);
  }
}

// Graceful shutdown handling
process.on("SIGINT", () => {
  console.log("\n🛑 Shutting down gateway gracefully...");
  process.exit(0);
});

process.on("SIGTERM", () => {
  console.log("\n🛑 Shutting down gateway gracefully...");
  process.exit(0);
});

// Start the gateway
startServer().catch((error) => {
  console.error("💥 Unhandled gateway startup error:", error);
  process.exit(1);
});

// Export for testing purposes
module.exports = { startServer };
