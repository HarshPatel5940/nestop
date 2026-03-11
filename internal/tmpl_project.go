package internal

var TmplPackageJSON = `{
  "name": "{{.ProjectName}}",
  "version": "0.1.0",
  "description": "",
  "private": true,
  "license": "UNLICENSED",
  "scripts": {
    "prebuild": "rimraf dist",
    "build": "nest build",
    "format": "biome format --write .",
    "lint": "biome check --write .",
    "start:dev": "nest start --watch",
    "start:debug": "nest start --debug --watch",
    "start:prod": "node dist/main"{{if .UsesVitest}},
    "test": "vitest run",
    "test:watch": "vitest",
    "test:cov": "vitest run --coverage",
    "test:ui": "vitest --ui"{{else}},
    "test": "jest",
    "test:watch": "jest --watch",
    "test:cov": "jest --coverage",
    "test:e2e": "jest --config ./test/jest-e2e.json"{{end}}{{if .UsesPrisma}},
    "db:generate": "prisma generate",
    "db:push": "prisma db push",
    "db:migrate": "prisma migrate dev",
    "db:deploy": "prisma migrate deploy",
    "db:studio": "prisma studio",
    "db:seed": "tsx prisma/seed.ts"{{end}}{{if .UsesDrizzle}},
    "db:generate": "drizzle-kit generate",
    "db:push": "drizzle-kit push",
    "db:migrate": "drizzle-kit migrate",
    "db:studio": "drizzle-kit studio"{{end}}{{if .IncludeHusky}},
    "prepare": "husky"{{end}}
  },
  "dependencies": {
    "@nestjs/common": "^11.0.1",
    "@nestjs/core": "^11.0.1",
    "@nestjs/platform-{{.Adapter}}": "^11.0.1",{{if .UsesFastify}}
    "@fastify/compress": "^8.0.0",
    "@fastify/cookie": "^11.0.2",
    "@fastify/cors": "^10.0.2",
    "@fastify/multipart": "^9.0.2",{{end}}
    "@nestjs/swagger": "^11.0.3",
    "@nestjs/throttler": "^6.4.0",{{if .EnableRedis}}
    "@nestjs/cache-manager": "^3.0.1",
    "cache-manager": "^6.3.2",
    "cache-manager-ioredis-yet": "^2.1.1",
    "ioredis": "^5.3.2",{{end}}{{if .HasS3}}
    "@aws-sdk/client-s3": "^3.726.0",
    "@aws-sdk/s3-request-presigner": "^3.726.0",{{end}}{{if .UsesPrisma}}
    "@prisma/client": "^6.2.1",{{end}}{{if .IsPostgresDrizzle}}
    "drizzle-orm": "^0.38.3",
    "pg": "^8.13.1",{{end}}{{if .IsMySQLDrizzle}}
    "drizzle-orm": "^0.38.3",
    "mysql2": "^3.9.0",{{end}}{{if .IsSQLiteDrizzle}}
    "drizzle-orm": "^0.38.3",
    "better-sqlite3": "^9.4.3",{{end}}{{if .UsesMongoose}}
    "@nestjs/mongoose": "^11.0.1",
    "mongoose": "^8.9.3",{{end}}
    "argon2": "^0.43.0",
    "jsonwebtoken": "^9.0.2",
    "nestjs-pino": "^4.2.0",
    "pino": "^9.6.0",
    "pino-http": "^10.3.0",
    "pino-pretty": "^13.0.0",
    "reflect-metadata": "^0.2.2",
    "rimraf": "^6.0.1",
    "rxjs": "^7.8.1",
    "ulid": "^2.3.0",
    "zod": "^3.24.1"
  },
  "devDependencies": {
    "@nestjs/cli": "^11.0.0",
    "@nestjs/schematics": "^11.0.0",
    "@nestjs/testing": "^11.0.1",
    "@biomejs/biome": "^1.9.4",{{if .IsSQLiteDrizzle}}
    "@types/better-sqlite3": "^7.6.8",{{end}}
    "@types/jsonwebtoken": "^9.0.7",
    "@types/node": "^22.10.9",{{if .NeedsPostgresDocker}}
    "@types/pg": "^8.11.10",{{end}}
{{if .IncludeHusky}}    "husky": "^9.1.7",
    "lint-staged": "^16.2.4",{{end}}{{if .UsesPrisma}}
    "prisma": "^6.2.1",{{end}}{{if .UsesDrizzle}}
    "drizzle-kit": "^0.30.1",{{end}}
    "source-map-support": "^0.5.21",
    "tsx": "^4.19.2",{{if .UsesVitest}}
    "vitest": "^3.0.3",
    "@vitest/coverage-v8": "^3.0.3"{{else}}
    "@types/jest": "^29.5.14",
    "jest": "^29.7.0",
    "ts-jest": "^29.2.5"{{end}}
  }
}
`

var TmplTSConfig = `{
  "compilerOptions": {
    "module": "commonjs",
    "declaration": true,
    "removeComments": true,
    "emitDecoratorMetadata": true,
    "experimentalDecorators": true,
    "allowSyntheticDefaultImports": true,
    "target": "ES2021",
    "sourceMap": true,
    "outDir": "./dist",
    "baseUrl": "./",
    "incremental": true,
    "skipLibCheck": true,
    "strictNullChecks": true,
    "noImplicitAny": false,
    "strictBindCallApply": false,
    "forceConsistentCasingInFileNames": false,
    "noFallthroughCasesInSwitch": false,
    "resolveJsonModule": true,
    "paths": {
      "@/*": ["src/*"]
    }
  }
}
`

var TmplTSConfigBuild = `{
  "extends": "./tsconfig.json",
  "exclude": ["node_modules", "test", "dist", "**/*spec.ts", "vitest.config.ts"]
}
`

var TmplNestCLIConfig = `{
  "$schema": "https://json.schemastore.org/nest-cli",
  "collection": "@nestjs/schematics",
  "sourceRoot": "src",
  "compilerOptions": {
    "deleteOutDir": true
  }
}
`

var TmplBiomeConfig = `{
  "$schema": "https://biomejs.dev/schemas/1.9.4/schema.json",
  "vcs": {
    "enabled": false,
    "clientKind": "git",
    "useIgnoreFile": false
  },
  "files": {
    "ignoreUnknown": false,
    "ignore": ["node_modules", "dist", "coverage", ".husky", "drizzle"]
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2,
    "lineWidth": 80
  },
  "organizeImports": {
    "enabled": true
  },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true,
      "correctness": {
        "useExhaustiveDependencies": "warn"
      },
      "style": {
        "noNonNullAssertion": "off"
      },
      "suspicious": {
        "noExplicitAny": "off"
      }
    }
  },
  "javascript": {
    "formatter": {
      "quoteStyle": "double",
      "trailingCommas": "all",
      "semicolons": "always"
    }
  }
}
`

var TmplVitestConfig = `import { defineConfig } from "vitest/config";
import { pathsToModuleNameMapper } from "ts-jest";
import { compilerOptions } from "./tsconfig.json";

export default defineConfig({
  test: {
    environment: "node",
    globals: true,
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      exclude: [
        "node_modules/",
        "dist/",
        "**/*.spec.ts",
        "**/index.ts",
        "**/*.dto.ts",
        "src/main.ts",
        "vitest.config.ts",
      ],
    },
  },
  resolve: {
    alias: {
      "@": "/src",
    },
  },
});
`

var TmplJestE2EConfig = `{
  "moduleFileExtensions": ["js", "json", "ts"],
  "rootDir": ".",
  "testEnvironment": "node",
  "testRegex": ".e2e-spec.ts$",
  "transform": {
    "^.+\\.(t|j)s$": "ts-jest"
  }
}
`

var TmplPinoMock = `import { PinoLogger } from "nestjs-pino";

export const mockPinoLoggerProvider = {
  provide: PinoLogger,
  useValue: {
    setContext: () => {},
    trace: () => {},
    debug: () => {},
    info: () => {},
    warn: () => {},
    error: () => {},
    fatal: () => {},
    log: () => {},
  },
};
`

var TmplReadme = `# {{.ProjectName}}

A NestJS application scaffolded with [nestop](https://github.com/harshpatel5940/nestop).

## Stack

- **Framework**: NestJS with {{if .UsesFastify}}Fastify{{else}}Express{{end}}
- **Package Manager**: {{.PackageManager}}
- **Language**: TypeScript{{if .UsesPrisma}}
- **ORM**: Prisma{{end}}{{if .UsesDrizzle}}
- **ORM**: Drizzle ({{.Database}}){{end}}{{if .UsesMongoose}}
- **ODM**: Mongoose{{end}}{{if .EnableRedis}}
- **Caching**: Redis{{end}}{{if .HasS3}}
- **Storage**: {{if .IsGarage}}Garage (S3-compatible){{else}}AWS S3{{end}}{{end}}
- **Auth**: JWT (Bearer token)
- **Validation**: Zod
- **Logging**: Pino
- **Linting**: Biome
- **Tests**: {{if .UsesVitest}}Vitest{{else}}Jest{{end}}

## Getting started

` + "```" + `bash
# Install dependencies
{{if .UsesBun}}bun install{{else}}pnpm install{{end}}

# Copy environment variables
cp .env.example .env

# Start services (Docker required)
make docker-up

{{if .UsesPrisma}}# Run database migrations
{{if .UsesBun}}bun run db:deploy{{else}}pnpm db:deploy{{end}}

{{end}}{{if .UsesDrizzle}}# Push database schema
{{if .UsesBun}}bun run db:push{{else}}pnpm db:push{{end}}

{{end}}# Start development server
{{if .UsesBun}}bun run start:dev{{else}}pnpm start:dev{{end}}
` + "```" + `

## API docs

Swagger UI: [http://localhost:3000/api/swagger](http://localhost:3000/api/swagger)
{{if .InitAuth}}
## Auth endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /auth/register | Register a new user |
| POST | /auth/login | Login and get JWT token |
| GET | /auth/session | Check current session (protected) |
{{end}}{{if .InitHealth}}
## Health endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Full health check |
| GET | /health/live | Liveness probe |
| GET | /health/ready | Readiness probe |
{{end}}
{{if .IncludeMakefile}}## Development

` + "```" + `bash
make dev        # Start dev server
make test       # Run tests
make lint       # Lint + format
make db-up      # Start database containers
make db-down    # Stop database containers
` + "```" + `
{{end}}
`
