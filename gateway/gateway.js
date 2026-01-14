const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");

// ===== subgraph URL selection =====
// Default path matches your Go services (change to "/graphql" if you switch your handlers)
const DEFAULT_PATH = process.env.SUBGRAPH_PATH || "/query";

function pickUrl(name, port) {
  const U = name.toUpperCase();

  // 1) Render Free / manual: full URL like https://...onrender.com/query (use as-is)
  const full = process.env[`${U}_URL`];
  if (full) return full;

  // 2) Render Paid (private): PRODUCTS_HOSTPORT=products-xxxx:4001 (compose full URL here)
  const hostport = process.env[`${U}_HOSTPORT`];
  if (hostport) return `http://${hostport}${DEFAULT_PATH}`;

  console.log(`🔍 Checking ${name}: FULL=${full}, HOSTPORT=${hostport}`);

  // 3) Local Docker fallback (Docker Compose service DNS)
  return `http://${name}:${port}${DEFAULT_PATH}`;
}

async function startServer() {
  try {
    console.log("🔄 Starting E-Commerce Federation Gateway...");

    const subgraphs = [
      { name: "products", url: pickUrl("products", 4001) },

      { name: "users", url: pickUrl("users", 4002) },
      { name: "orders", url: pickUrl("orders", 4003) },
    ];

    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs,
        introspectionHeaders: { "User-Agent": "ApolloGateway/2.5.5" },
      }),
      pollIntervalInMs: 10000,
      buildService: ({ url }) =>
        new RemoteGraphQLDataSource({
          url,
          willSendRequest: ({ request }) => {
            request.http.headers.set("apollo-federation-include-trace", "ftv1");
          },
        }),
    });

    const server = new ApolloServer({
      gateway,
      introspection: true,
      csrfPrevention: { requestHeaders: ["apollo-required-preflight"] },
      plugins: [
        ApolloServerPluginLandingPageLocalDefault({
          embed: true,
          settings: { "editor.theme": "dark", "editor.fontSize": 14 },
        }),
        {
          requestDidStart() {
            return {
              didResolveOperation(rc) {
                console.log(`📊 Query: ${rc.request.operationName || "Anonymous"}`);
              },
              didEncounterErrors(rc) {
                console.error("❌ GraphQL errors:", rc.errors);
              },
            };
          },
        },
      ],
      formatError: (error) => ({
        message: error.message,
        code: error.extensions?.code,
        path: error.path,
      }),
    });

    const { url } = await startStandaloneServer(server, {
      listen: { port: Number(process.env.PORT) || 4000, host: "0.0.0.0" },
    });

    console.log("✅ Federation Gateway Successfully Started!");
    console.log(`🚀 Gateway ready at ${url}`);
    console.log("📋 Connected Services:");
    for (const sg of subgraphs) console.log(`  • ${sg.name}: ${sg.url}`);
  } catch (error) {
    console.error("💥 Failed to start federation gateway:", error);
    if (/ECONNREFUSED|ENOTFOUND/.test(error.message)) {
      console.error("💡 Check subgraph URLs (env vars) and that services are reachable.");
    }
    process.exit(1);
  }
}

// Graceful shutdown
process.on("SIGINT", () => {
  console.log("\n🛑 Shutting down gateway gracefully...");
  process.exit(0);
});
process.on("SIGTERM", () => {
  console.log("\n🛑 Shutting down gateway gracefully...");
  process.exit(0);
});

startServer().catch((error) => {
  console.error("💥 Unhandled gateway startup error:", error);
  process.exit(1);
});

module.exports = { startServer };

// const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");
// const { ApolloServer } = require("@apollo/server");
// const { startStandaloneServer } = require("@apollo/server/standalone");
// const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");

// async function startServer() {
//   try {
//     console.log("🔄 Starting E-Commerce Federation Gateway...");

//     // Create the federation gateway
//     const gateway = new ApolloGateway({
//       supergraphSdl: new IntrospectAndCompose({
//         // subgraphs: [
//         //   { name: "products", url: "http://products:4001/query" },
//         //   { name: "users", url: "http://users:4002/query" },
//         //   { name: "orders", url: "http://orders:4003/query" },
//         // ],
//         subgraphs: [
//           { name: "products", url: process.env.PRODUCTS_URL || "http://products:4001/graphql" },
//           { name: "users", url: process.env.USERS_URL || "http://users:4002/graphql" },
//           { name: "orders", url: process.env.ORDERS_URL || "http://orders:4003/graphql" },
//         ],
//         introspectionHeaders: {
//           "User-Agent": "ApolloGateway/2.5.5",
//         },
//       }),

//       // Poll every 10 seconds for schema changes
//       pollIntervalInMs: 10000,
//       // Add explicit federation configuration
//       buildService: ({ url }) => {
//         return new RemoteGraphQLDataSource({
//           url,
//           willSendRequest: ({ request, context }) => {
//             request.http.headers.set("apollo-federation-include-trace", "ftv1");
//           },
//         });
//       },
//     });

//     // Create Apollo Server with the gateway
//     const server = new ApolloServer({
//       gateway,
//       introspection: true, // Enable introspection for development
//       csrfPrevention: {
//         requestHeaders: ["apollo-required-preflight"],
//       },
//       plugins: [
//         // Enable GraphQL Playground
//         ApolloServerPluginLandingPageLocalDefault({
//           embed: true,
//           settings: {
//             "editor.theme": "dark",
//             "editor.fontSize": 14,
//           },
//         }),

//         // Custom plugin for health checks and logging
//         {
//           requestDidStart() {
//             return {
//               didResolveOperation(requestContext) {
//                 console.log(`📊 Query: ${requestContext.request.operationName || "Anonymous"}`);
//               },
//               didEncounterErrors(requestContext) {
//                 console.error("❌ GraphQL errors:", requestContext.errors);
//               },
//             };
//           },
//         },
//       ],
//       // Custom error formatting
//       formatError: (error) => {
//         console.error("🚨 Gateway Error:", error);
//         return {
//           message: error.message,
//           code: error.extensions?.code,
//           path: error.path,
//         };
//       },
//     });

//     // Start the server
//     // const { url } = await startStandaloneServer(server, {
//     //   listen: { port: 4000 },
//     // });
//     const { url } = await startStandaloneServer(server, {
//       listen: {
//         port: process.env.PORT || 4000,
//         host: "0.0.0.0",
//       },
//     });
//     // Success logging
//     console.log("✅ Federation Gateway Successfully Started!");
//     console.log("");
//     console.log(`🚀 Gateway ready at ${url}`);
//     console.log(`🎮 GraphQL Playground: ${url}`);
//     console.log("");
//     console.log("📋 Connected Services:");
//     console.log("  🛍️  Products Service: http://localhost:4001/query");
//     console.log("  👥 Users Service: http://localhost:4002/query");
//     console.log("  📦 Orders Service: http://localhost:4003/query");
//     console.log("");
//     console.log("🔗 Example Federated Query:");
//     console.log(`
//     query ExampleCrossServiceQuery {
//       user(id: "1") {
//         name
//         email
//         password
//         # This will automatically resolve across services
//       }
//       products {
//         name
//         price
//       }
//       orders {
//         id
//         status
//       }
//     }`);
//   } catch (error) {
//     console.error("💥 Failed to start federation gateway:", error);

//     // Error messages
//     if (error.message.includes("ECONNREFUSED")) {
//       console.error("");
//       console.error("🚨 Connection Error: Make sure all services are running:");
//       console.error("  - Products service on port 4001");
//       console.error("  - Users service on port 4002");
//       console.error("  - Orders service on port 4003");
//       console.error("");
//       console.error("💡 Start each service with: go run main.go");
//     }

//     process.exit(1);
//   }
// }

// // Graceful shutdown handling
// process.on("SIGINT", () => {
//   console.log("\n🛑 Shutting down gateway gracefully...");
//   process.exit(0);
// });

// process.on("SIGTERM", () => {
//   console.log("\n🛑 Shutting down gateway gracefully...");
//   process.exit(0);
// });

// // Start the gateway
// startServer().catch((error) => {
//   console.error("💥 Unhandled gateway startup error:", error);
//   process.exit(1);
// });

// // Export for testing purposes
// module.exports = { startServer };
