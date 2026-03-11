package internal

type PackageManager string
type Adapter string
type S3Provider string
type DatabaseChoice string
type DatabaseORM string
type TestFramework string

const (
PMPnpm PackageManager = "pnpm"
PMBun  PackageManager = "bun"

AdapterFastify Adapter = "fastify"
AdapterExpress Adapter = "express"

S3Garage S3Provider = "garage"
S3AWS    S3Provider = "aws"
S3None   S3Provider = "none"

DBPostgres DatabaseChoice = "postgres"
DBMySQL    DatabaseChoice = "mysql"
DBSQLite   DatabaseChoice = "sqlite"
DBMongoDB  DatabaseChoice = "mongodb"
DBNone     DatabaseChoice = "none"

ORMPrisma   DatabaseORM = "prisma"
ORMDrizzle  DatabaseORM = "drizzle"
ORMMongoose DatabaseORM = "mongoose"

TestVitest TestFramework = "vitest"
TestJest   TestFramework = "jest"
)

type ProjectConfig struct {
// Core
ProjectName    string
PackageManager PackageManager
Adapter        Adapter

// Providers
EnableRedis bool
S3Provider  S3Provider

// Database
Database DatabaseChoice
ORM      DatabaseORM

// Tests + endpoints
TestFramework TestFramework
InitHealth    bool
InitAuth      bool

// Tooling (all optional, asked in form)
IncludeMakefile bool
IncludeCI       bool
IncludeVSCode   bool
IncludeHusky    bool

// Post-scaffold actions
InstallDeps bool
InitGit     bool
}

func (t PackageManager) String() string { return string(t) }
func (t Adapter) String() string        { return string(t) }
func (t S3Provider) String() string     { return string(t) }
func (t DatabaseChoice) String() string { return string(t) }
func (t DatabaseORM) String() string    { return string(t) }
func (t TestFramework) String() string  { return string(t) }

func (c *ProjectConfig) UsesFastify() bool  { return c.Adapter == AdapterFastify }
func (c *ProjectConfig) UsesExpress() bool  { return c.Adapter == AdapterExpress }
func (c *ProjectConfig) HasDatabase() bool  { return c.Database != DBNone }
func (c *ProjectConfig) UsesPrisma() bool   { return c.ORM == ORMPrisma }
func (c *ProjectConfig) UsesDrizzle() bool  { return c.ORM == ORMDrizzle }
func (c *ProjectConfig) UsesMongoose() bool { return c.ORM == ORMMongoose }
func (c *ProjectConfig) HasS3() bool        { return c.S3Provider != S3None }
func (c *ProjectConfig) IsGarage() bool     { return c.S3Provider == S3Garage }
func (c *ProjectConfig) IsAWSS3() bool      { return c.S3Provider == S3AWS }
func (c *ProjectConfig) UsesVitest() bool   { return c.TestFramework == TestVitest }
func (c *ProjectConfig) UsesJest() bool     { return c.TestFramework == TestJest }
func (c *ProjectConfig) UsesBun() bool      { return c.PackageManager == PMBun }
func (c *ProjectConfig) UsesPnpm() bool     { return c.PackageManager == PMPnpm }

func (c *ProjectConfig) NeedsPostgresDocker() bool { return c.Database == DBPostgres }
func (c *ProjectConfig) NeedsMySQLDocker() bool    { return c.Database == DBMySQL }
func (c *ProjectConfig) NeedsMongoDocker() bool    { return c.Database == DBMongoDB }
func (c *ProjectConfig) IsSQLite() bool            { return c.Database == DBSQLite }

func (c *ProjectConfig) IsPostgresDrizzle() bool {
return c.Database == DBPostgres && c.ORM == ORMDrizzle
}
func (c *ProjectConfig) IsMySQLDrizzle() bool {
return c.Database == DBMySQL && c.ORM == ORMDrizzle
}
func (c *ProjectConfig) IsSQLiteDrizzle() bool {
return c.Database == DBSQLite && c.ORM == ORMDrizzle
}
func (c *ProjectConfig) DrizzleDialect() string {
switch c.Database {
case DBMySQL:
return "mysql"
case DBSQLite:
return "sqlite"
default:
return "postgresql"
}
}
