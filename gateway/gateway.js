const { ApolloGateway, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");
const { readFileSync } = require("fs");

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function getStatusFromError(err) {
  return err?.extensions?.response?.status || err?.response?.status || err?.extensions?.status || null;
}

class RetryableDataSource extends RemoteGraphQLDataSource {
  async process(options) {
    const maxAttempts = 3;

    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
      try {
        return await super.process(options);
      } catch (err) {
        const status = getStatusFromError(err);
        const retryable = [429, 502, 503, 504].includes(Number(status));
        const subgraphUrl = this.url || "unknown-subgraph";

        console.error(
          `❌ Subgraph request failed: ${subgraphUrl} | status=${status} | attempt=${attempt}/${maxAttempts}`,
        );

        if (!retryable || attempt === maxAttempts) {
          throw err;
        }

        const delayMs = Math.min(10000, attempt * 3000);
        console.warn(`⏳ Retrying ${subgraphUrl} in ${delayMs}ms...`);
        await sleep(delayMs);
      }
    }
  }
}

async function startServer() {
  try {
    console.log("🔄 Starting E-Commerce Federation Gateway with Static SDL...");

    const supergraphSdl = readFileSync("./supergraph.graphql").toString();

    const gateway = new ApolloGateway({
      supergraphSdl,
      buildService: ({ url }) =>
        new RetryableDataSource({
          url,
          willSendRequest: ({ request }) => {
            request.http.headers.set("apollo-federation-include-trace", "ftv1");
          },
        }),
    });

    const isProd = process.env.NODE_ENV === "production";

    const server = new ApolloServer({
      gateway,
      introspection: true,
      csrfPrevention: { requestHeaders: ["apollo-required-preflight"] },
      plugins: isProd
        ? [
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
          ]
        : [
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
