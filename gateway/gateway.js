const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");

// ===== subgraph URL selection =====
function pickUrl(name, port) {
  const U = name.toUpperCase();

  // Render free tier: full external URL like https://products.onrender.com
  const full = process.env[`${U}_URL`];
  if (full) {
    console.log(`✅ Using external URL for ${name}: ${full}/query`);
    return `${full}/query`;
  }

  // Local Docker fallback (Docker Compose service DNS)
  console.log(`🐳 Using local Docker URL for ${name}: http://${name}:${port}/query`);
  return `http://${name}:${port}/query`;
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
