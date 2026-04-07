const { ApolloGateway, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");
const { readFileSync } = require("fs");

async function startServer() {
  try {
    console.log("🔄 Starting E-Commerce Federation Gateway with Static SDL...");

    // 1. Load the supergraph file from your root directory
    const supergraphSdl = readFileSync("./supergraph.graphql").toString();

    // 2. Initialize Gateway (No retry needed - it's a local file!)
    const gateway = new ApolloGateway({
      supergraphSdl,
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
  } catch (error) {
    console.error("💥 Failed to start federation gateway:", error);
    process.exit(1);
  }
}

// Graceful shutdown
process.on("SIGINT", () => {
  console.log("\n🛑 Shutting down...");
  process.exit(0);
});
process.on("SIGTERM", () => {
  console.log("\n🛑 Shutting down...");
  process.exit(0);
});

startServer().catch((error) => {
  console.error("💥 Unhandled gateway startup error:", error);
  process.exit(1);
});

module.exports = { startServer };
