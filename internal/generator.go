package internal

import (
"bytes"
"fmt"
"os"
"path/filepath"
"strings"
"text/template"
)

type FileEntry struct {
Path    string
Content string
}

func collectFiles(c *ProjectConfig) []FileEntry {
var files []FileEntry

add := func(path, tmpl string) {
files = append(files, FileEntry{path, tmpl})
}

// ── Core ──
add("src/main.ts", TmplMainTS)
add("src/app.module.ts", TmplAppModule)
add("src/app.controller.ts", TmplAppController)
add("src/app.service.ts", TmplAppService)
add("src/app.controller.spec.ts", TmplAppControllerSpec)

// ── Config ──
add("src/config/index.ts", TmplConfigIndex)
add("src/config/env.ts", TmplConfigEnv)
add("src/config/swagger.config.ts", TmplSwaggerConfig)

// ── Pipes & Filters ──
add("src/pipes/zod.pipe.ts", TmplZodPipe)
add("src/filters/global-exception.filter.ts", TmplGlobalExceptionFilter)
if c.UsesPrisma() {
add("src/filters/prisma-exception.filter.ts", TmplPrismaExceptionFilter)
}

// ── Middleware ──
add("src/middleware/correlation-id.middleware.ts", TmplCorrelationIdMiddleware)

// ── API: constants ──
add("src/api/constants/index.ts", TmplConstantsIndex)
add("src/api/constants/errors.ts", TmplConstantsErrors)
add("src/api/constants/response.ts", TmplConstantsResponse)

// ── API: decorators ──
add("src/api/decorators/index.ts", TmplDecoratorsIndex)
add("src/api/decorators/public.decorator.ts", TmplPublicDecorator)
add("src/api/decorators/roles.decorator.ts", TmplRolesDecorator)
add("src/api/decorators/current-user.decorator.ts", TmplCurrentUserDecorator)

// ── API: guards ──
add("src/api/guards/index.ts", TmplGuardsIndex)
add("src/api/guards/auth.guard.ts", TmplAuthGuard)
add("src/api/guards/roles.guard.ts", TmplRolesGuard)

// ── API: interceptors ──
add("src/api/interceptors/index.ts", TmplInterceptorsIndex)
add("src/api/interceptors/response.interceptor.ts", TmplResponseInterceptor)

// ── API: utils ──
add("src/api/utils/index.ts", TmplUtilsIndex)
add("src/api/utils/hash.utils.ts", TmplHashUtils)
add("src/api/utils/ulid.utils.ts", TmplUlidUtils)
add("src/api/utils/jwt.utils.ts", TmplJWTUtils)
add("src/api/utils/pagination.utils.ts", TmplPaginationUtils)

// ── Providers: Redis ──
if c.EnableRedis {
add("src/providers/redis/redis.module.ts", TmplRedisModule)
add("src/providers/redis/redis.service.ts", TmplRedisService)
}

// ── Providers: S3 ──
if c.HasS3() {
add("src/providers/s3/s3.module.ts", TmplS3Module)
add("src/providers/s3/s3.service.ts", TmplS3Service)
}

// ── Database ──
if c.HasDatabase() {
switch {
case c.UsesPrisma():
add("src/database/prisma.module.ts", TmplPrismaModule)
add("src/database/prisma.service.ts", TmplPrismaService)
add("prisma/schema.prisma", TmplPrismaSchema)

case c.IsPostgresDrizzle():
add("src/database/db.ts", TmplDrizzleDB)
add("src/database/schema/users.schema.ts", TmplDrizzleUsersSchema)
add("src/database/schema/index.ts", TmplDrizzleSchemaIndex)
add("src/database/types.ts", TmplDrizzleTypes)
add("drizzle.config.ts", TmplDrizzleConfig)

case c.IsMySQLDrizzle():
add("src/database/db.ts", TmplDrizzleDBMySQL)
add("src/database/schema/users.schema.ts", TmplDrizzleUsersSchemaMySQL)
add("src/database/schema/index.ts", TmplDrizzleSchemaIndex)
add("src/database/types.ts", TmplDrizzleTypesMySQL)
add("drizzle.config.ts", TmplDrizzleConfigMySQL)

case c.IsSQLiteDrizzle():
add("src/database/db.ts", TmplDrizzleDBSQLite)
add("src/database/schema/users.schema.ts", TmplDrizzleUsersSchemaSQLite)
add("src/database/schema/index.ts", TmplDrizzleSchemaIndex)
add("src/database/types.ts", TmplDrizzleTypesSQLite)
add("drizzle.config.ts", TmplDrizzleConfigSQLite)

case c.UsesMongoose():
add("src/database/database.module.ts", TmplMongooseModule)
add("src/database/schemas/user.schema.ts", TmplMongooseUserSchema)
}
}

// ── Routes: health ──
if c.InitHealth {
add("src/api/routes/health/health.module.ts", TmplHealthModule)
add("src/api/routes/health/health.controller.ts", TmplHealthController)
add("src/api/routes/health/health.service.ts", TmplHealthService)
}

// ── Routes: auth ──
if c.InitAuth {
add("src/api/routes/auth/auth.module.ts", TmplAuthModule)
add("src/api/routes/auth/auth.controller.ts", TmplAuthController)
add("src/api/routes/auth/auth.service.ts", TmplAuthService)
add("src/api/routes/auth/dto/index.ts", TmplAuthDTOs)
}

// ── Project config ──
add("package.json", TmplPackageJSON)
add("tsconfig.json", TmplTSConfig)
add("tsconfig.build.json", TmplTSConfigBuild)
add("nest-cli.json", TmplNestCLIConfig)
add("biome.json", TmplBiomeConfig)
if c.UsesVitest() {
add("vitest.config.ts", TmplVitestConfig)
} else {
add("test/jest-e2e.json", TmplJestE2EConfig)
}
add("test/mock/pino-mock.ts", TmplPinoMock)
add("README.md", TmplReadme)

// ── Developer tooling ──
if c.IncludeMakefile {
		add("Makefile", TmplMakefile)
	}
	if c.IncludeCI {
		add(".github/workflows/ci.yml", TmplGitHubActionsCI)
	}
	if c.IncludeVSCode {
		add(".vscode/extensions.json", TmplVSCodeExtensions)
		add(".vscode/settings.json", TmplVSCodeSettings)
	}
	if c.IncludeHusky {
		add("lintstagedrc.json", TmplLintStagedRC)
		add(".husky/pre-commit", TmplHuskyPreCommit)
	}

// ── Docker & env ──
add("Dockerfile", TmplDockerfile)
add("docker-compose.yml", TmplDockerCompose)
add(".env", TmplEnvFile)
add(".env.example", TmplEnvExampleFile)
add(".gitignore", TmplGitignore)
add(".dockerignore", TmplDockerignore)

return files
}

// collectEmptyDirs returns dirs to create even when no files go inside
func collectEmptyDirs(c *ProjectConfig) []string {
var dirs []string
if !c.HasDatabase() {
dirs = append(dirs, "src/database")
}
if !c.HasS3() {
dirs = append(dirs, "src/providers/s3")
}
if !c.EnableRedis {
dirs = append(dirs, "src/providers/redis")
}
if !c.InitHealth && !c.InitAuth {
dirs = append(dirs, "src/api/routes")
}
return dirs
}

// ─── Template rendering ───

var funcMap = template.FuncMap{
"lower":    strings.ToLower,
"title":    strings.Title,
"contains": strings.Contains,
}

func renderTemplate(name, tmplStr string, data *ProjectConfig) (string, error) {
t, err := template.New(name).Funcs(funcMap).Parse(tmplStr)
if err != nil {
return "", fmt.Errorf("parse template %s: %w", name, err)
}
var buf bytes.Buffer
if err := t.Execute(&buf, data); err != nil {
return "", fmt.Errorf("execute template %s: %w", name, err)
}
return buf.String(), nil
}

// ─── Generate ───

func Generate(config *ProjectConfig) error {
root := config.ProjectName

files := collectFiles(config)
emptyDirs := collectEmptyDirs(config)

// Create all directories
dirSet := map[string]bool{}
for _, f := range files {
dirSet[filepath.Dir(filepath.Join(root, f.Path))] = true
}
for _, d := range emptyDirs {
dirSet[filepath.Join(root, d)] = true
}
for dir := range dirSet {
if err := os.MkdirAll(dir, 0o755); err != nil {
return fmt.Errorf("mkdir %s: %w", dir, err)
}
}

// Render and write files
for _, f := range files {
rendered, err := renderTemplate(f.Path, f.Content, config)
if err != nil {
return err
}
dest := filepath.Join(root, f.Path)
if err := os.WriteFile(dest, []byte(rendered), 0o644); err != nil {
return fmt.Errorf("write %s: %w", dest, err)
}
// Make husky hook executable
if strings.HasPrefix(f.Path, ".husky/") {
os.Chmod(dest, 0o755)
}
}

fmt.Printf("  📁 Generated %d files in ./%s\n", len(files), root)
return nil
}
