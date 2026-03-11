package internal

// ──────────────────── PRISMA ────────────────────

var TmplPrismaModule = `import { Global, Module } from "@nestjs/common";
import { PrismaService } from "./prisma.service";

@Global()
@Module({
  providers: [PrismaService],
  exports: [PrismaService],
})
export class PrismaModule {}
`

var TmplPrismaService = `import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
} from "@nestjs/common";
import { PinoLogger } from "nestjs-pino";
import { PrismaClient } from "@prisma/client";

@Injectable()
export class PrismaService
  extends PrismaClient
  implements OnModuleInit, OnModuleDestroy
{
  constructor(private readonly logger: PinoLogger) {
    super({
      log: [
        { level: "query", emit: "event" },
        { level: "warn", emit: "stdout" },
        { level: "error", emit: "stdout" },
      ],
    });
    this.logger.setContext(PrismaService.name);
  }

  async onModuleInit() {
    await this.$connect();
    this.logger.info("Prisma connected");
  }

  async onModuleDestroy() {
    await this.$disconnect();
    this.logger.info("Prisma disconnected");
  }
}
`

var TmplPrismaSchema = `generator client {
  provider = "prisma-client-js"
}

datasource db {
{{- if .NeedsPostgresDocker}}
  provider = "postgresql"
  url      = env("DATABASE_URL")
{{- else if .NeedsMySQLDocker}}
  provider = "mysql"
  url      = env("DATABASE_URL")
{{- else if .IsSQLite}}
  provider = "sqlite"
  url      = env("DATABASE_URL")
{{- else if .NeedsMongoDocker}}
  provider = "mongodb"
  url      = env("DATABASE_URL")
{{- end}}
}

model User {
  id        String   @id @default(cuid())
  name      String
  email     String   @unique
  password  String
  role      String   @default("user")
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
`

// ──────────────────── DRIZZLE POSTGRESQL ────────────────────

var TmplDrizzleDB = `import { drizzle } from "drizzle-orm/node-postgres";
import { Pool } from "pg";
import { env } from "../config";

const pool = new Pool({
  connectionString: env.databaseUrl,
  max: 10,
});

export const db = drizzle(pool);
`

var TmplDrizzleConfig = `import type { Config } from "drizzle-kit";
import { env } from "./src/config";

export default {
  schema: "./src/database/schema/*.ts",
  out: "./drizzle",
  dialect: "postgresql",
  dbCredentials: {
    url: env.databaseUrl,
  },
} satisfies Config;
`

var TmplDrizzleUsersSchema = `import { pgTable, text, timestamp } from "drizzle-orm/pg-core";
import { generateULID } from "../api/utils/ulid.utils";

export const users = pgTable("users", {
  id: text("id")
    .primaryKey()
    .$defaultFn(() => generateULID()),
  name: text("name").notNull(),
  email: text("email").notNull().unique(),
  password: text("password").notNull(),
  role: text("role").notNull().default("user"),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type User = typeof users.$inferSelect;
export type NewUser = typeof users.$inferInsert;
`

var TmplDrizzleSchemaIndex = `export * from "./users.schema";
`

var TmplDrizzleTypes = `import type { NodePgDatabase } from "drizzle-orm/node-postgres";
import type * as schema from "./schema";

export type DrizzleDB = NodePgDatabase<typeof schema>;
`

// ──────────────────── DRIZZLE MYSQL ────────────────────

var TmplDrizzleDBMySQL = `import { drizzle } from "drizzle-orm/mysql2";
import mysql from "mysql2/promise";
import { env } from "../config";

const pool = mysql.createPool({
  uri: env.databaseUrl,
  connectionLimit: 10,
});

export const db = drizzle({ client: pool });
`

var TmplDrizzleConfigMySQL = `import type { Config } from "drizzle-kit";
import { env } from "./src/config";

export default {
  schema: "./src/database/schema/*.ts",
  out: "./drizzle",
  dialect: "mysql",
  dbCredentials: {
    url: env.databaseUrl,
  },
} satisfies Config;
`

var TmplDrizzleUsersSchemaMySQL = `import { mysqlTable, varchar, timestamp } from "drizzle-orm/mysql-core";
import { generateULID } from "../api/utils/ulid.utils";

export const users = mysqlTable("users", {
  id: varchar("id", { length: 26 })
    .primaryKey()
    .$defaultFn(() => generateULID()),
  name: varchar("name", { length: 100 }).notNull(),
  email: varchar("email", { length: 255 }).notNull().unique(),
  password: varchar("password", { length: 255 }).notNull(),
  role: varchar("role", { length: 50 }).notNull().default("user"),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type User = typeof users.$inferSelect;
export type NewUser = typeof users.$inferInsert;
`

var TmplDrizzleTypesMySQL = `import type { MySql2Database } from "drizzle-orm/mysql2";
import type * as schema from "./schema";

export type DrizzleDB = MySql2Database<typeof schema>;
`

// ──────────────────── DRIZZLE SQLITE ────────────────────

var TmplDrizzleDBSQLite = `import { drizzle } from "drizzle-orm/better-sqlite3";
import Database from "better-sqlite3";
import { env } from "../config";

const sqlite = new Database(env.databaseUrl);
export const db = drizzle(sqlite);
`

var TmplDrizzleConfigSQLite = `import type { Config } from "drizzle-kit";
import { env } from "./src/config";

export default {
  schema: "./src/database/schema/*.ts",
  out: "./drizzle",
  dialect: "sqlite",
  dbCredentials: {
    url: env.databaseUrl,
  },
} satisfies Config;
`

var TmplDrizzleUsersSchemaSQLite = `import { sqliteTable, text, integer } from "drizzle-orm/sqlite-core";
import { generateULID } from "../api/utils/ulid.utils";
import { sql } from "drizzle-orm";

export const users = sqliteTable("users", {
  id: text("id")
    .primaryKey()
    .$defaultFn(() => generateULID()),
  name: text("name").notNull(),
  email: text("email").notNull().unique(),
  password: text("password").notNull(),
  role: text("role").notNull().default("user"),
  createdAt: text("created_at").default(sql` + "`(datetime('now'))`" + `).notNull(),
  updatedAt: text("updated_at").default(sql` + "`(datetime('now'))`" + `).notNull(),
});

export type User = typeof users.$inferSelect;
export type NewUser = typeof users.$inferInsert;
`

var TmplDrizzleTypesSQLite = `import type { BetterSQLite3Database } from "drizzle-orm/better-sqlite3";
import type * as schema from "./schema";

export type DrizzleDB = BetterSQLite3Database<typeof schema>;
`

// ──────────────────── MONGOOSE ────────────────────

var TmplMongooseModule = `import { Module } from "@nestjs/common";
import { MongooseModule } from "@nestjs/mongoose";
import { env } from "../config";

@Module({
  imports: [
    MongooseModule.forRootAsync({
      useFactory: () => ({
        uri: env.databaseUrl,
      }),
    }),
  ],
  exports: [MongooseModule],
})
export class DatabaseModule {}
`

var TmplMongooseUserSchema = `import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { Document } from "mongoose";
import { generateULID } from "../api/utils/ulid.utils";

export type UserDocument = User & Document;

@Schema({ timestamps: true })
export class User {
  @Prop({ type: String, default: () => generateULID() })
  _id!: string;

  @Prop({ required: true })
  name!: string;

  @Prop({ required: true, unique: true, lowercase: true, trim: true })
  email!: string;

  @Prop({ required: true })
  password!: string;

  @Prop({ default: "user" })
  role!: string;
}

export const UserSchema = SchemaFactory.createForClass(User);
`
