const { ApolloGateway, IntrospectAndCompose, RemoteGraphQLDataSource } = require("@apollo/gateway");
const { ApolloServer } = require("@apollo/server");
const { startStandaloneServer } = require("@apollo/server/standalone");
const { ApolloServerPluginLandingPageLocalDefault } = require("@apollo/server/plugin/landingPage/default");

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

function parseRetryAfterSeconds(errMsg) {
  const m = /retry-after[:=]\s*(\d+)/i.exec(errMsg);
  return m ? Number(m[1]) : null;
}

function isTransientStartupError(err) {
  const msg = (err && (err.message || String(err))) || "";
  return (
    msg.includes("429") ||
    msg.includes("Too Many Requests") ||
    msg.includes("503") ||
    msg.includes("502") ||
    msg.includes("504") ||
    msg.includes("ECONNREFUSED") ||
    msg.includes("ENOTFOUND") ||
    msg.includes("ETIMEDOUT") ||
    msg.includes("socket hang up") ||
    msg.includes("Couldn't load service definitions")
  );
}

async function withRetry(fn, { maxAttempts = 8, baseDelayMs = 1500, maxDelayMs = 20000 } = {}) {
  let attempt = 0;
  while (true) {
    attempt += 1;
    try {
      return await fn(attempt);
    } catch (err) {
      const msg = err?.message || String(err);
      const transient = isTransientStartupError(err);

      if (!transient || attempt >= maxAttempts) {
        throw err;
      }

      const retryAfter = parseRetryAfterSeconds(msg);
      const backoff = Math.min(maxDelayMs, baseDelayMs * Math.pow(2, attempt - 1));
      const jitter = Math.floor(Math.random() * 400);
      const waitMs = retryAfter != null ? retryAfter * 1000 : backoff + jitter;

      console.warn(`⏳ Gateway startup retry ${attempt}/${maxAttempts} after error: ${msg}`);
      console.warn(`🔁 Waiting ${Math.round(waitMs / 1000)}s before retry...`);
      await sleep(waitMs);
    }
  }
}

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

    const started = await withRetry(async () => {
      const gateway = new ApolloGateway({
        supergraphSdl: new IntrospectAndCompose({
          subgraphs,
          introspectionHeaders: { "User-Agent": "ApolloGateway/2.5.5" },
        }),
        // Polling can cause repeated schema fetches (extra traffic).
        // Consider increasing or disabling for Render demos.
        pollIntervalInMs: 60000,
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

      return { url, subgraphs };
    });

    console.log("✅ Federation Gateway Successfully Started!");
    console.log(`🚀 Gateway ready at ${started.url}`);
    console.log("📋 Connected Services:");
    for (const sg of started.subgraphs) console.log(`  • ${sg.name}: ${sg.url}`);
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
