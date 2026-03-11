# nestop 🪺

An interactive CLI tool to scaffold NestJS projects with best practices, built with Go and [charmbracelet/huh](https://github.com/charmbracelet/huh).

## Features

- **Interactive TUI** — beautiful form-based project setup
- **Package Managers** — pnpm or bun
- **HTTP Adapters** — Fastify (default, with compress/cookie/cors/multipart) or Express
- **Validation** — Zod pipes (object schema validation)
- **Logging** — Pino (nestjs-pino) with pino-pretty for dev
- **IDs** — ULID with Zod schema validation
- **Hashing** — Argon2 for password hashing
- **Providers** — Redis caching, S3-compatible storage (Garage / AWS)
- **Databases** — PostgreSQL/SQL/MongoDB with Prisma, Drizzle, or Mongoose
- **Testing** — Vitest or Jest
- **Starter Endpoints** — Health check and Auth (register/login/session)
- **Docker** — Multi-stage Dockerfile + docker-compose with all services
- **Code Quality** — Biome linter/formatter, Swagger API docs, throttling

## Install

```bash
go install github.com/harshpatel5940/nestop@latest
```

Or build from source:

```bash
git clone https://github.com/harshpatel5940/nestop.git
cd nestop
go build -o nestop .
```

## Usage

```bash
./nestop
```

The CLI will guide you through:

1. **Project name** and **package manager** (pnpm/bun)
2. **HTTP adapter** (Fastify/Express)
3. **Providers** — Redis caching, S3 storage (Garage/AWS/skip)
4. **Database** — PostgreSQL/SQL/MongoDB + ORM (Prisma/Drizzle/Mongoose/skip)
5. **Test framework** (Vitest/Jest)
6. **Starter endpoints** (Health, Auth)

## Generated Structure

```
my-project/
├── src/
│   ├── main.ts, app.module.ts, app.controller.ts, app.service.ts
│   ├── config/          # Zod-validated env, Swagger config
│   ├── pipes/           # Zod validation pipe
│   ├── filters/         # Global + Prisma exception filters
│   ├── api/
│   │   ├── constants/   # Error types, response types
│   │   ├── guards/      # Auth guard with JWT placeholder
│   │   ├── interceptors/ # Response interceptor
│   │   ├── utils/       # Argon2 hash, ULID, JWT helpers
│   │   └── routes/      # health/, auth/ endpoints
│   ├── providers/       # Redis, S3 modules
│   └── database/        # Prisma/Drizzle/Mongoose setup
├── test/                # Test config + pino mock
├── Dockerfile           # Multi-stage (pnpm or bun)
├── docker-compose.yml   # App + DB + Redis + S3
├── package.json, tsconfig.json, biome.json, nest-cli.json
└── .env, .env.example, .gitignore, .dockerignore
```

## License

MIT

