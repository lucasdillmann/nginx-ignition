import fs from "node:fs"
import os from "node:os"
import path from "node:path"
import { baseUrl, credentials, httpPort, httpsPort } from "./config.mjs"

const log = message => console.log(`[seed] ${message}`)

class ApiClient {
    constructor(token) {
        this.token = token
    }

    async request(method, endpoint, body) {
        const response = await fetch(`${baseUrl}${endpoint}`, {
            method,
            headers: {
                ...(body === undefined ? {} : { "Content-Type": "application/json" }),
                ...(this.token === undefined ? {} : { Authorization: `Bearer ${this.token}` }),
            },
            body: body === undefined ? undefined : JSON.stringify(body),
        })

        if (!response.ok) {
            const text = await response.text()
            throw new Error(`${method} ${endpoint} failed with ${response.status}: ${text}`)
        }

        if (response.status === 204) return undefined

        const text = await response.text()
        return text === "" ? undefined : JSON.parse(text)
    }

    get(endpoint) {
        return this.request("GET", endpoint)
    }

    post(endpoint, body) {
        return this.request("POST", endpoint, body)
    }

    put(endpoint, body) {
        return this.request("PUT", endpoint, body)
    }
}

export async function authenticate() {
    const anonymous = new ApiClient()

    const onboarding = await anonymous.post("/api/users/onboarding/finish", {
        name: credentials.name,
        username: credentials.username,
        password: credentials.password,
    })

    const api = new ApiClient(onboarding.token)
    log(`created the "${credentials.username}" administrator account`)

    return api
}

function findDockerSocket() {
    const candidates = [
        process.env.CAPTURE_DOCKER_SOCKET,
        "/var/run/docker.sock",
        path.join(os.homedir(), ".docker", "run", "docker.sock"),
    ].filter(Boolean)

    return candidates.find(candidate => fs.existsSync(candidate))
}

const routeSettings = {
    includeForwardHeaders: false,
    proxySslServerName: false,
    ignoreSslErrors: false,
    keepOriginalDomainName: false,
    directoryListingEnabled: false,
}

const hostFeatureSet = {
    websocketsSupport: false,
    http2Support: true,
    redirectHttpToHttps: false,
    statsEnabled: false,
}

const emptyCache = {
    inactiveSeconds: 600,
    maximumSizeMb: 128,
    allowedMethods: ["GET", "HEAD"],
    minimumUsesBeforeCaching: 1,
    useStale: ["ERROR", "TIMEOUT", "UPDATING"],
}

export async function seed(api) {
    const accessList = await api.post("/api/access-lists", {
        name: "Public endpoints",
        realm: "nginx ignition",
        satisfyAll: true,
        defaultOutcome: "ALLOW",
        entries: [{ priority: 1, outcome: "ALLOW", sourceAddresses: ["0.0.0.0/0", "::/0"] }],
        forwardAuthenticationHeader: false,
    })
    const protectedAccessList = await api.post("/api/access-lists", {
        name: "Administration only",
        realm: "Restricted area",
        satisfyAll: true,
        defaultOutcome: "DENY",
        entries: [
            { priority: 1, outcome: "ALLOW", sourceAddresses: ["10.0.0.0/8", "192.168.0.0/16"] },
            { priority: 2, outcome: "DENY", sourceAddresses: ["0.0.0.0/0", "::/0"] },
        ],
        credentials: [{ username: "demo", password: "demo-password" }],
        forwardAuthenticationHeader: false,
    })
    log("created 2 access lists")

    const existingCaches = (await api.get("/api/caches?pageSize=100")).contents
    let cache = existingCaches[0]

    if (cache === undefined) {
        cache = await api.post("/api/caches", emptyCache)
    }

    await api.put(`/api/caches/${cache.id}`, {
        ...emptyCache,
        name: cache.name,
        backgroundUpdate: true,
        concurrencyLock: { enabled: true, timeoutSeconds: 5, ageSeconds: 10 },
        revalidate: true,
        cacheStatusResponseHeaderEnabled: true,
        ignoreUpstreamCacheHeaders: false,
        bypassRules: ["/api/*"],
        noCacheRules: ["/logout"],
        fileExtensions: ["css", "js", "png", "jpg", "svg", "woff2"],
        durations: [{ statusCodes: ["200", "301", "302"], validTimeSeconds: 3600 }],
    })
    log(`reused the cache configuration created by the migrations (${cache.name})`)

    const certificate = await api.post("/api/certificates/issue", {
        providerId: "SELF_SIGNED",
        domainNames: ["docs.example.com", "example.com"],
        parameters: {},
    })
    log("issued 1 self-signed certificate")

    const dockerSocket = findDockerSocket()
    let integration
    if (dockerSocket !== undefined) {
        integration = await api.post("/api/integrations", {
            name: "Local Docker",
            driver: "DOCKER",
            enabled: true,
            parameters: {
                connectionMode: "SOCKET",
                socketPath: dockerSocket,
                hostUrl: "",
                swarmMode: false,
                swarmServiceMesh: false,
                swarmDnsResolvers: "",
                useContainerNameAsId: true,
                proxyUrl: "",
            },
        })
        log(`connected the "Local Docker" integration through ${dockerSocket}`)
    }

    const vpn = await api.post("/api/vpns", {
        name: "Corporate network",
        driver: "TAILSCALE",
        enabled: false,
        parameters: { authKey: "tskey-auth-0000000000000000000000000000" },
    })
    log("created 1 disabled VPN")

    const existingHosts = await api.get("/api/hosts?pageSize=100")
    const defaultHost = existingHosts.contents.find(host => host.defaultServer)

    if (defaultHost !== undefined) {
        await api.put(`/api/hosts/${defaultHost.id}`, {
            ...defaultHost,
            enabled: true,
            domainNames: [],
            routes: [
                {
                    priority: 0,
                    enabled: true,
                    type: "STATIC_RESPONSE",
                    protocol: "HTTP_1_1",
                    sourcePath: "/",
                    settings: routeSettings,
                    response: {
                        statusCode: 200,
                        payload: "<h1>nginx ignition</h1><p>Your server is up and running.</p>",
                        headers: { "Content-Type": "text/html; charset=utf-8" },
                    },
                },
            ],
        })
        log("reconfigured the default virtual host")
    }

    const blogHost = await api.post("/api/hosts", {
        enabled: true,
        defaultServer: false,
        useGlobalBindings: true,
        domainNames: ["blog.example.com"],
        featureSet: { ...hostFeatureSet, websocketsSupport: true, redirectHttpToHttps: true },
        accessListId: accessList.id,
        routes: [
            {
                priority: 0,
                enabled: true,
                type: "PROXY",
                protocol: "HTTP_1_1",
                sourcePath: "/",
                settings: { ...routeSettings, keepOriginalDomainName: true },
                targetUri: "http://127.0.0.1:3000",
            },
        ],
    })

    const integrationOptions =
        integration === undefined
            ? []
            : (await api.get(`/api/integrations/${integration.id}/options?pageSize=20`)).contents
    const integrationOption = integrationOptions[0]

    const dashboardHost = await api.post("/api/hosts", {
        enabled: true,
        defaultServer: false,
        useGlobalBindings: false,
        domainNames: ["dashboard.example.com", "metrics.example.com"],
        featureSet: { ...hostFeatureSet, redirectHttpToHttps: true, statsEnabled: true },
        accessListId: protectedAccessList.id,
        cacheId: cache.id,
        bindings: [{ type: "HTTPS", ip: "0.0.0.0", port: 8443, certificateId: certificate.certificateId }],
        routes: [
            {
                priority: 0,
                enabled: true,
                type: "PROXY",
                protocol: "HTTP_1_1",
                sourcePath: "/",
                settings: { ...routeSettings, includeForwardHeaders: true, keepOriginalDomainName: true },
                targetUri: "https://metrics.example.com",
            },
            {
                priority: 1,
                enabled: true,
                type: "REDIRECT",
                protocol: "HTTP_1_1",
                sourcePath: "/legacy",
                settings: routeSettings,
                redirectCode: 301,
                targetUri: "https://metrics.example.com/dashboards",
            },
            ...(integrationOption === undefined
                ? []
                : [
                      {
                          priority: 2,
                          enabled: true,
                          type: "INTEGRATION",
                          protocol: "HTTP_1_1",
                          sourcePath: "/apps",
                          settings: routeSettings,
                          integration: {
                              integrationId: integration.id,
                              optionId: integrationOption.id,
                              useHttps: false,
                          },
                      },
                  ]),
        ],
    })
    log("created 3 virtual hosts")

    const postgresStream = await api.post("/api/streams", {
        enabled: true,
        name: "PostgreSQL",
        type: "SIMPLE",
        featureSet: {
            useProxyProtocol: false,
            socketKeepAlive: false,
            tcpKeepAlive: true,
            tcpNoDelay: true,
            tcpDeferred: true,
        },
        defaultBackend: { target: { protocol: "TCP", address: "postgres.example.com", port: 5432 } },
        binding: { protocol: "TCP", address: "0.0.0.0", port: 15432 },
    })

    const internalStream = await api.post("/api/streams", {
        enabled: true,
        name: "Internal services",
        type: "SNI_ROUTER",
        featureSet: {
            useProxyProtocol: true,
            socketKeepAlive: true,
            tcpKeepAlive: true,
            tcpNoDelay: true,
            tcpDeferred: true,
        },
        defaultBackend: {
            target: { protocol: "TCP", address: "fallback.example.com", port: 443 },
            circuitBreaker: { maxFailures: 3, openSeconds: 30 },
        },
        binding: { protocol: "TCP", address: "0.0.0.0", port: 443 },
        routes: [
            {
                domainNames: ["api.example.com"],
                backends: [{ target: { protocol: "TCP", address: "api.example.com", port: 8443 }, weight: 3 }],
            },
            {
                domainNames: ["auth.example.com"],
                backends: [{ target: { protocol: "UDP", address: "auth.example.com", port: 51820 } }],
            },
        ],
    })
    log("created 2 streams")

    const settings = await api.get("/api/settings")
    await api.put("/api/settings", {
        ...settings,
        globalBindings: [
            { type: "HTTP", ip: "127.0.0.1", port: httpPort },
            {
                type: "HTTPS",
                ip: "127.0.0.1",
                port: httpsPort,
                certificateId: certificate.certificateId,
            },
        ],
        nginx: {
            ...settings.nginx,
            gzipEnabled: true,
            serverTokensEnabled: false,
            maximumBodySizeMb: 2048,
            timeouts: { ...settings.nginx.timeouts, read: 600, send: 600, keepalive: 65 },
        },
    })
    log(`adjusted the server settings to listen on ${httpPort} and ${httpsPort}`)

    await api.post("/api/nginx/reload").catch(() => api.post("/api/nginx/start"))
    log(`reloaded nginx on ${httpPort} and ${httpsPort}`)

    await generateTraffic()

    const user = (await api.get("/api/users?pageSize=10")).contents.find(item => item.username === credentials.username)

    return {
        user,
        hosts: { blog: blogHost.id, dashboard: dashboardHost.id, default: defaultHost?.id },
        streams: { postgres: postgresStream.id, internal: internalStream.id },
        accessLists: { public: accessList.id, restricted: protectedAccessList.id },
        cache: { id: cache.id, name: cache.name },
        certificate: { id: certificate.certificateId },
        integration: { id: integration?.id },
        vpn: { id: vpn.id },
    }
}

async function generateTraffic() {
    const origin = `http://127.0.0.1:${httpPort}`
    const paths = ["/", "/", "/", "/about", "/assets/app.css", "/api/health", "/legacy"]
    const proxied = { headers: { host: "blog.example.com" } }
    const requests = []

    for (let round = 0; round < 3; round++) {
        for (const item of paths) {
            requests.push(fetch(`${origin}${item}`, { redirect: "manual" }))
        }

        requests.push(fetch(origin, proxied), fetch(`${origin}/legacy`, proxied))
    }

    await Promise.all(requests.map(request => request.catch(() => undefined)))
    await new Promise(resolve => setTimeout(resolve, 500))

    log("generated sample traffic through the seeded nginx")
}
