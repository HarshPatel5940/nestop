package internal

import (
"fmt"
"regexp"
"strings"

"charm.land/huh/v2"
)

func RunForm() (*ProjectConfig, error) {
config := &ProjectConfig{}

// ── defaults ──
projectName := "my-nest-app"
adapter := string(AdapterFastify)
pm := string(PMPnpm)
enableRedis := true
s3Choice := string(S3None)
dbChoice := "postgres-prisma"
testChoice := string(TestVitest)

// endpoints: pre-populate so Value() sees them selected when Options() runs
initEndpoints := []string{"health", "auth"}

// tooling (all default no — opt-in)
includeMakefile := false
includeCI := false
includeVSCode := false
includeHusky := false
installDeps := true
initGit := true

form := huh.NewForm(
// ── Group 1: basics ──
huh.NewGroup(
huh.NewInput().
Title("Project name").
Description("Lowercase letters, digits, hyphens only").
Placeholder("my-nest-app").
Validate(func(s string) error {
if len(strings.TrimSpace(s)) == 0 {
return fmt.Errorf("project name is required")
}
if matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]*$`, s); !matched {
return fmt.Errorf("use lowercase letters, digits and hyphens only")
}
return nil
}).
Value(&projectName),

huh.NewSelect[string]().
Title("Package manager").
Options(
huh.NewOption("pnpm (Recommended)", string(PMPnpm)),
huh.NewOption("bun", string(PMBun)),
).
Value(&pm),

huh.NewSelect[string]().
Title("HTTP adapter").
Options(
huh.NewOption("Fastify (Recommended)", string(AdapterFastify)),
huh.NewOption("Express", string(AdapterExpress)),
).
Value(&adapter),
),

// ── Group 2: providers ──
huh.NewGroup(
huh.NewConfirm().
Title("Enable Redis caching?").
Description("Adds ioredis client + CacheModule").
Affirmative("Yes").
Negative("No").
Value(&enableRedis),

huh.NewSelect[string]().
Title("S3 object storage").
Options(
huh.NewOption("Skip", string(S3None)),
huh.NewOption("Garage (self-hosted S3-compatible)", string(S3Garage)),
huh.NewOption("AWS S3", string(S3AWS)),
).
Value(&s3Choice),

huh.NewSelect[string]().
Title("Database + ORM").
Description("Choose your database and ORM combination").
Options(
huh.NewOption("PostgreSQL + Prisma (Recommended)", "postgres-prisma"),
huh.NewOption("PostgreSQL + Drizzle", "postgres-drizzle"),
huh.NewOption("MySQL + Prisma", "mysql-prisma"),
huh.NewOption("MySQL + Drizzle", "mysql-drizzle"),
huh.NewOption("SQLite + Prisma", "sqlite-prisma"),
huh.NewOption("SQLite + Drizzle", "sqlite-drizzle"),
huh.NewOption("MongoDB + Prisma", "mongodb-prisma"),
huh.NewOption("MongoDB + Mongoose", "mongodb-mongoose"),
huh.NewOption("Skip (empty folders only)", "none"),
).
Value(&dbChoice),
),

// ── Group 3: tests ──
huh.NewGroup(
huh.NewSelect[string]().
Title("Test framework").
Options(
huh.NewOption("Jest (Recommended)", string(TestJest)),
huh.NewOption("Vitest", string(TestVitest)),
).
Value(&testChoice),
),

// ── Group 4: starter endpoints — isolated group fixes the empty-options bug ──
huh.NewGroup(
huh.NewMultiSelect[string]().
Title("Starter endpoints").
Description("Space to toggle, Enter to confirm. Both selected by default.").
// Value MUST come before Options — selectOptions() checks the accessor
Value(&initEndpoints).
Options(
huh.NewOption("Health (GET /health, /health/live, /health/ready)", "health"),
huh.NewOption("Auth (POST /auth/register, /auth/login, GET /auth/session)", "auth"),
).
Height(6).
Filterable(false),
),

// ── Group 5: tooling ──
huh.NewGroup(
huh.NewConfirm().
Title("Include Makefile?").
Description("Adds dev/build/test/lint/db-up/docker-up targets").
Affirmative("Yes").Negative("No").
Value(&includeMakefile),

huh.NewConfirm().
Title("Include GitHub Actions CI?").
Description("Adds .github/workflows/ci.yml for push + PR checks").
Affirmative("Yes").Negative("No").
Value(&includeCI),

huh.NewConfirm().
Title("Include VS Code config?").
Description("Adds .vscode/extensions.json + settings.json (Biome formatter)").
Affirmative("Yes").Negative("No").
Value(&includeVSCode),

huh.NewConfirm().
Title("Include Husky + lint-staged?").
Description("Adds pre-commit hook that runs Biome on staged files").
Affirmative("Yes").Negative("No").
Value(&includeHusky),
),

// ── Group 6: post-scaffold actions ──
huh.NewGroup(
huh.NewConfirm().
Title("Initialize git repository?").
Description("Runs git init and creates an initial commit").
Affirmative("Yes").Negative("No").
Value(&initGit),

huh.NewConfirm().
Title("Install dependencies now?").
Description("Runs pnpm install / bun install after scaffold").
Affirmative("Yes").Negative("No").
Value(&installDeps),
),
).
WithTheme(huh.ThemeFunc(huh.ThemeCatppuccin)).
WithAccessible(false)

if err := form.Run(); err != nil {
return nil, err
}

// ── map → config ──
config.ProjectName = strings.TrimSpace(projectName)
config.PackageManager = PackageManager(pm)
config.Adapter = Adapter(adapter)
config.EnableRedis = enableRedis
config.S3Provider = S3Provider(s3Choice)
config.TestFramework = TestFramework(testChoice)
config.IncludeMakefile = includeMakefile
config.IncludeCI = includeCI
config.IncludeVSCode = includeVSCode
config.IncludeHusky = includeHusky
config.InstallDeps = installDeps
config.InitGit = initGit

// parse db+orm combo
switch dbChoice {
case "postgres-prisma":
config.Database, config.ORM = DBPostgres, ORMPrisma
case "postgres-drizzle":
config.Database, config.ORM = DBPostgres, ORMDrizzle
case "mysql-prisma":
config.Database, config.ORM = DBMySQL, ORMPrisma
case "mysql-drizzle":
config.Database, config.ORM = DBMySQL, ORMDrizzle
case "sqlite-prisma":
config.Database, config.ORM = DBSQLite, ORMPrisma
case "sqlite-drizzle":
config.Database, config.ORM = DBSQLite, ORMDrizzle
case "mongodb-prisma":
config.Database, config.ORM = DBMongoDB, ORMPrisma
case "mongodb-mongoose":
config.Database, config.ORM = DBMongoDB, ORMMongoose
default:
config.Database = DBNone
}

// parse endpoint multiselect
for _, ep := range initEndpoints {
switch ep {
case "health":
config.InitHealth = true
case "auth":
config.InitAuth = true
}
}

return config, nil
}
