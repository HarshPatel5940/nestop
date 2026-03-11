package internal

// ─── Makefile ───
// Note: uses \t for tab indentation (required by make)
var TmplMakefile = "{{if .UsesBun}}PM=bun{{else}}PM=pnpm{{end}}\n\n.PHONY: dev build start test lint fmt db-up db-down docker-up docker-down\n\ndev:\n\t$(PM) start:dev\n\nbuild:\n\t$(PM) build\n\nstart:\n\t$(PM) start:prod\n\ntest:\n\t$(PM) test\n\ntest-cov:\n\t$(PM) test:cov\n\nlint:\n\t$(PM) lint\n\nfmt:\n\tbiome format --write .\n\ndb-up:{{if .NeedsPostgresDocker}}\n\tdocker compose up postgres-db -d{{end}}{{if .NeedsMySQLDocker}}\n\tdocker compose up mysql-db -d{{end}}{{if .NeedsMongoDocker}}\n\tdocker compose up mongo-db -d{{end}}{{if .EnableRedis}}\n\tdocker compose up redis -d{{end}}\n\ndb-down:{{if .NeedsPostgresDocker}}\n\tdocker compose rm postgres-db -s -f -v{{end}}{{if .NeedsMySQLDocker}}\n\tdocker compose rm mysql-db -s -f -v{{end}}{{if .NeedsMongoDocker}}\n\tdocker compose rm mongo-db -s -f -v{{end}}\n\ndb-restart: db-down\n\tsleep 1\n\t$(MAKE) db-up{{if .UsesPrisma}}\n\t$(PM) db:deploy{{end}}{{if .UsesDrizzle}}\n\t$(PM) db:push{{end}}\n\ndocker-up:\n\tdocker compose up -d\n\ndocker-down:\n\tdocker compose down\n\ndocker-logs:\n\tdocker compose logs -f app\n"

// ─── .github/workflows/ci.yml ───
var TmplGitHubActionsCI = `name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  lint-build-test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
{{if .UsesBun}}
      - name: Setup Bun
        uses: oven-sh/setup-bun@v2
        with:
          bun-version: latest

      - name: Install dependencies
        run: bun install --frozen-lockfile

      - name: Lint
        run: bun run lint

      - name: Build
        run: bun run build

      - name: Test
        run: bun run test
{{else}}
      - name: Setup pnpm
        uses: pnpm/action-setup@v4
        with:
          version: latest

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: pnpm

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Lint
        run: pnpm lint

      - name: Build
        run: pnpm build

      - name: Test
        run: pnpm test
{{end}}
`

// ─── .vscode/extensions.json ───
var TmplVSCodeExtensions = `{
  "recommendations": [
    "biomejs.biome",
    "firsttris.vscode-jest-runner",
    "dbaeumer.vscode-eslint",
    "ms-vscode.vscode-typescript-next",
    "streetsidesoftware.code-spell-checker",
    "bradlc.vscode-tailwindcss",
    "ms-azuretools.vscode-docker",
    "humao.rest-client",
    "arjun.swagger-viewer",
    "42crunch.vscode-openapi",
    "mikestead.dotenv"
  ]
}
`

// ─── .vscode/settings.json ───
var TmplVSCodeSettings = `{
  "editor.defaultFormatter": "biomejs.biome",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports.biome": "explicit"
  },
  "[typescript]": {
    "editor.defaultFormatter": "biomejs.biome"
  },
  "[json]": {
    "editor.defaultFormatter": "biomejs.biome"
  }
}
`

// ─── lintstagedrc.json ───
var TmplLintStagedRC = `{
  "*.{ts,js,json,md}": [
    "biome lint --write",
    "biome format --write"
  ]
}
`

// ─── .husky/pre-commit ───
var TmplHuskyPreCommit = `#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

npx lint-staged
`

// ─── src/api/decorators/index.ts ───
var TmplDecoratorsIndex = `export * from "./public.decorator";
export * from "./roles.decorator";
export * from "./current-user.decorator";
`

// ─── src/api/decorators/public.decorator.ts ───
var TmplPublicDecorator = `import { SetMetadata } from "@nestjs/common";

export const IS_PUBLIC_KEY = "isPublic";
export const Public = () => SetMetadata(IS_PUBLIC_KEY, true);
`

// ─── src/api/decorators/roles.decorator.ts ───
var TmplRolesDecorator = `import { SetMetadata } from "@nestjs/common";

export const ROLES_KEY = "roles";
export const Roles = (...roles: string[]) => SetMetadata(ROLES_KEY, roles);
`

// ─── src/api/decorators/current-user.decorator.ts ───
var TmplCurrentUserDecorator = `import { createParamDecorator, ExecutionContext } from "@nestjs/common";

export const CurrentUser = createParamDecorator(
  (_data: unknown, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    return request.user;
  },
);
`

// ─── src/api/guards/roles.guard.ts ───
var TmplRolesGuard = `import { Injectable, CanActivate, ExecutionContext } from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { ROLES_KEY } from "../decorators/roles.decorator";

@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(context: ExecutionContext): boolean {
    const requiredRoles = this.reflector.getAllAndOverride<string[]>(ROLES_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);

    if (!requiredRoles || requiredRoles.length === 0) {
      return true;
    }

    const { user } = context.switchToHttp().getRequest();
    return requiredRoles.some((role) => user?.role === role);
  }
}
`

// ─── src/middleware/correlation-id.middleware.ts ───
var TmplCorrelationIdMiddleware = `import { Injectable, NestMiddleware } from "@nestjs/common";
import { generateULID } from "../api/utils";

@Injectable()
export class CorrelationIdMiddleware implements NestMiddleware {
  use(req: any, res: any, next: () => void) {
    const correlationId =
      (req.headers["x-correlation-id"] as string) || generateULID();
    req.headers["x-correlation-id"] = correlationId;
    res.setHeader("x-correlation-id", correlationId);
    next();
  }
}
`

// ─── src/api/utils/pagination.utils.ts ───
var TmplPaginationUtils = `export interface PaginationQuery {
  page?: number | string;
  limit?: number | string;
}

export interface PaginationMeta {
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

export interface PaginatedResult<T> {
  data: T[];
  meta: PaginationMeta;
}

export const normalizePagination = (
  query: PaginationQuery,
): { page: number; limit: number; offset: number } => {
  const page = Math.max(1, Number(query.page ?? 1));
  const limit = Math.min(100, Math.max(1, Number(query.limit ?? 20)));
  const offset = (page - 1) * limit;
  return { page, limit, offset };
};

export const paginate = <T>(
  data: T[],
  total: number,
  params: { page: number; limit: number },
): PaginatedResult<T> => {
  const totalPages = Math.ceil(total / params.limit);
  return {
    data,
    meta: {
      total,
      page: params.page,
      limit: params.limit,
      totalPages,
      hasNextPage: params.page < totalPages,
      hasPreviousPage: params.page > 1,
    },
  };
};
`

// ─── src/api/utils/jwt.utils.ts ───
var TmplJWTUtils = `import * as jwt from "jsonwebtoken";
import { env } from "../../config";

export interface JwtPayload {
  sub: string;
  email?: string;
  role?: string;
  iat?: number;
  exp?: number;
}

export const signJwt = (
  payload: Omit<JwtPayload, "iat" | "exp">,
): string => {
  return jwt.sign(payload, env.jwtSecret, {
    expiresIn: env.jwtExpiresIn as any,
  });
};

export const verifyJwt = (token: string): JwtPayload => {
  return jwt.verify(token, env.jwtSecret) as JwtPayload;
};
`
