package internal

var TmplDockerfile = `# ──────────── builder ────────────
FROM node:22-alpine AS builder
WORKDIR /app
{{if .UsesBun}}
RUN npm install -g bun
COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile
{{else}}
RUN npm install -g pnpm
COPY package.json pnpm-lock.yaml* ./
RUN pnpm install --frozen-lockfile
{{end}}
COPY . .
{{if .UsesPrisma}}RUN {{if .UsesBun}}bun run{{else}}pnpm{{end}} db:generate
{{end}}RUN {{if .UsesBun}}bun run{{else}}pnpm{{end}} build

# ──────────── runner ────────────
FROM node:22-alpine AS runner
WORKDIR /app

RUN addgroup --system --gid 1001 nodejs \
  && adduser --system --uid 1001 nestjs
{{if .UsesBun}}
RUN npm install -g bun
COPY --from=builder /app/node_modules ./node_modules
{{else}}
RUN npm install -g pnpm
COPY --from=builder /app/node_modules ./node_modules
{{end}}
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./package.json
{{if .UsesPrisma}}COPY --from=builder /app/prisma ./prisma
{{end}}
USER nestjs

EXPOSE 3000
ENV NODE_ENV=production

CMD ["node", "dist/main"]
`

var TmplDockerCompose = `services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ${APP_NAME:-{{.ProjectName}}}-app
    environment:
      NODE_ENV: ${NODE_ENV:-production}
      PORT: ${PORT:-3000}
      DATABASE_URL: ${DATABASE_URL}{{if .EnableRedis}}
      REDIS_URL: ${REDIS_URL}{{end}}{{if .HasS3}}
      S3_BUCKET_NAME: ${S3_BUCKET_NAME}
      S3_REGION: ${S3_REGION}
      S3_ACCESS_KEY_ID: ${S3_ACCESS_KEY_ID}
      S3_SECRET_ACCESS_KEY: ${S3_SECRET_ACCESS_KEY}{{if .IsGarage}}
      S3_ENDPOINT: ${S3_ENDPOINT}{{end}}{{end}}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - "${PORT:-3000}:3000"
    depends_on:{{if .NeedsPostgresDocker}}
      postgres-db:
        condition: service_healthy{{end}}{{if .NeedsMySQLDocker}}
      mysql-db:
        condition: service_healthy{{end}}{{if .NeedsMongoDocker}}
      mongo-db:
        condition: service_healthy{{end}}{{if .EnableRedis}}
      redis:
        condition: service_healthy{{end}}
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{if .NeedsPostgresDocker}}
  postgres-db:
    image: postgres:17-alpine
    container_name: ${APP_NAME:-{{.ProjectName}}}-postgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      POSTGRES_DB: ${POSTGRES_DB:-{{.ProjectName}}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "${POSTGRES_HOST_PORT:-5432}:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{end}}{{if .NeedsMySQLDocker}}
  mysql-db:
    image: mysql:8.4
    container_name: ${APP_NAME:-{{.ProjectName}}}-mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-root}
      MYSQL_DATABASE: ${MYSQL_DB:-{{.ProjectName}}}
      MYSQL_USER: ${MYSQL_USER:-mysql}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD:-mysql}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "${MYSQL_HOST_PORT:-3306}:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p${MYSQL_ROOT_PASSWORD:-root}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{end}}{{if .NeedsMongoDocker}}
  mongo-db:
    image: mongo:8.0
    container_name: ${APP_NAME:-{{.ProjectName}}}-mongo
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER:-mongo}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD:-mongo}
      MONGO_INITDB_DATABASE: ${MONGO_DB:-{{.ProjectName}}}
    volumes:
      - mongo_data:/data/db
    ports:
      - "${MONGO_HOST_PORT:-27017}:27017"
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{end}}{{if .EnableRedis}}
  redis:
    image: redis:7-alpine
    container_name: ${APP_NAME:-{{.ProjectName}}}-redis
    command: redis-server --requirepass ${REDIS_PASSWORD:-redis}
    volumes:
      - redis_data:/data
    ports:
      - "${REDIS_HOST_PORT:-6379}:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD:-redis}", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{end}}{{if .IsGarage}}
  garage:
    image: dxflrs/garage:v1.0.0
    container_name: ${APP_NAME:-{{.ProjectName}}}-garage
    environment:
      GARAGE_REPLICATION_MODE: "none"
      GARAGE_RPC_BIND_ADDR: "[::]:3901"
      GARAGE_RPC_PUBLIC_ADDR: "127.0.0.1:3901"
      GARAGE_S3_API_BIND_ADDR: "[::]:3900"
      GARAGE_ADMIN_API_BIND_ADDR: "[::]:3903"
      GARAGE_METADATA_DIR: /var/lib/garage/meta
      GARAGE_DATA_DIR: /var/lib/garage/data
    volumes:
      - garage_data:/var/lib/garage/data
      - garage_meta:/var/lib/garage/meta
    ports:
      - "${GARAGE_S3_PORT:-3900}:3900"
      - "${GARAGE_ADMIN_PORT:-3903}:3903"
    networks:
      - {{.ProjectName}}-network
    restart: unless-stopped
{{end}}
networks:
  {{.ProjectName}}-network:
    driver: bridge

volumes:{{if .NeedsPostgresDocker}}
  postgres_data:{{end}}{{if .NeedsMySQLDocker}}
  mysql_data:{{end}}{{if .NeedsMongoDocker}}
  mongo_data:{{end}}{{if .EnableRedis}}
  redis_data:{{end}}{{if .IsGarage}}
  garage_data:
  garage_meta:{{end}}
`

var TmplEnvFile = `# Application
NODE_ENV=development
PORT=3000
APP_NAME={{.ProjectName}}
FRONTEND_URL=http://localhost:5173

# Auth
JWT_SECRET=change-me-to-a-secure-random-string-at-least-32-chars
JWT_EXPIRES_IN=7d
{{if .NeedsPostgresDocker}}
# Database (PostgreSQL)
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/{{.ProjectName}}
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB={{.ProjectName}}
{{end}}{{if .NeedsMySQLDocker}}
# Database (MySQL)
DATABASE_URL=mysql://mysql:mysql@localhost:3306/{{.ProjectName}}
MYSQL_ROOT_PASSWORD=root
MYSQL_DB={{.ProjectName}}
MYSQL_USER=mysql
MYSQL_PASSWORD=mysql
{{end}}{{if .IsSQLite}}
# Database (SQLite)
DATABASE_URL=./dev.db
{{end}}{{if .NeedsMongoDocker}}
# Database (MongoDB)
DATABASE_URL=mongodb://mongo:mongo@localhost:27017/{{.ProjectName}}?authSource=admin
MONGO_USER=mongo
MONGO_PASSWORD=mongo
MONGO_DB={{.ProjectName}}
{{end}}{{if .EnableRedis}}
# Redis
REDIS_URL=redis://:redis@localhost:6379
REDIS_PASSWORD=redis
{{end}}{{if .IsGarage}}
# Garage S3
S3_BUCKET_NAME={{.ProjectName}}-bucket
S3_REGION=garage
S3_ACCESS_KEY_ID=your-garage-key-id
S3_SECRET_ACCESS_KEY=your-garage-secret-key
S3_ENDPOINT=http://localhost:3900
{{end}}{{if .IsAWSS3}}
# AWS S3
S3_BUCKET_NAME={{.ProjectName}}-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=your-aws-access-key
S3_SECRET_ACCESS_KEY=your-aws-secret-key
{{end}}`

var TmplEnvExampleFile = `# Application
NODE_ENV=development
PORT=3000
APP_NAME={{.ProjectName}}
FRONTEND_URL=http://localhost:5173

# Auth
JWT_SECRET=change-me-to-a-secure-random-string-at-least-32-chars
JWT_EXPIRES_IN=7d
{{if .NeedsPostgresDocker}}
# Database (PostgreSQL)
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/{{.ProjectName}}
{{end}}{{if .NeedsMySQLDocker}}
# Database (MySQL)
DATABASE_URL=mysql://user:password@localhost:3306/{{.ProjectName}}
{{end}}{{if .IsSQLite}}
# Database (SQLite)
DATABASE_URL=./dev.db
{{end}}{{if .NeedsMongoDocker}}
# Database (MongoDB)
DATABASE_URL=mongodb://user:password@localhost:27017/{{.ProjectName}}?authSource=admin
{{end}}{{if .EnableRedis}}
# Redis
REDIS_URL=redis://:password@localhost:6379
{{end}}{{if .HasS3}}
# S3 / Object Storage
S3_BUCKET_NAME=
S3_REGION=
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY={{if .IsGarage}}
S3_ENDPOINT={{end}}
{{end}}`

var TmplGitignore = `# Dependencies
node_modules/
.pnp
.pnp.js

# Build
dist/
build/

# Environment
.env
.env.local
.env.*.local

# Logs
logs/
*.log
npm-debug.log*
pnpm-debug.log*

# Database
dev.db
dev.db-journal
drizzle/

# Coverage
coverage/

# Cache
.cache/
.turbo/

# OS
.DS_Store
Thumbs.db

# Editor
.vscode/
!.vscode/extensions.json
!.vscode/settings.json
.idea/

# Misc
*.tsbuildinfo
`

var TmplDockerignore = `node_modules
dist
coverage
.git
.env
*.log
.DS_Store
README.md
`
