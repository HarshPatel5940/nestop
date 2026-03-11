package internal

// ─── config/index.ts ───
var TmplConfigIndex = `export * from "./env";
`

// ─── config/env.ts ───
var TmplConfigEnv = `import { config } from "dotenv";
import { z } from "zod";

config();

const envSchema = z.object({
  NODE_ENV: z
    .enum(["development", "production", "test"])
    .default("development"),
  PORT: z.coerce.number().default(3000),
  APP_NAME: z.string().default("{{.ProjectName}}"),
  FRONTEND_URL: z.string().default("http://localhost:5173"),
{{if .HasDatabase}}{{if eq .Database.String "mongodb"}}  DATABASE_URL: z.string({ required_error: "DATABASE_URL is required" }),
{{else}}  DATABASE_URL: z.string({ required_error: "DATABASE_URL is required" }),
  DIRECT_URL: z.string().optional(),
{{end}}{{end}}{{if .EnableRedis}}  REDIS_URL: z.string().default("redis://localhost:6379"),
{{end}}{{if .HasS3}}  S3_BUCKET_NAME: z.string().default("{{.ProjectName}}-bucket"),
  S3_REGION: z.string().default("us-east-1"),
  S3_ACCESS_KEY_ID: z.string({ required_error: "S3_ACCESS_KEY_ID is required" }),
  S3_SECRET_ACCESS_KEY: z.string({ required_error: "S3_SECRET_ACCESS_KEY is required" }),{{if .IsGarage}}
  S3_ENDPOINT: z.string().default("http://localhost:3900"),{{end}}
{{end}}  JWT_SECRET: z.string({ required_error: "JWT_SECRET is required" }),
  JWT_EXPIRES_IN: z.string().trim().min(1).default("7d"),

  LOG_LEVEL: z
    .enum(["error", "warn", "info", "debug", "verbose"])
    .default("info"),
});

const parsed = envSchema.safeParse(process.env);

if (!parsed.success) {
  console.error("❌ Invalid environment variables:", parsed.error.format());
  throw new Error("Invalid environment variables");
}

export const env = {
  nodeEnv: parsed.data.NODE_ENV,
  port: parsed.data.PORT,
  appName: parsed.data.APP_NAME,
  frontendUrl: parsed.data.FRONTEND_URL,
{{if .HasDatabase}}  databaseUrl: parsed.data.DATABASE_URL,
{{if ne .Database.String "mongodb"}}  directUrl: parsed.data.DIRECT_URL,
{{end}}{{end}}{{if .EnableRedis}}  redisUrl: parsed.data.REDIS_URL,
{{end}}{{if .HasS3}}  s3BucketName: parsed.data.S3_BUCKET_NAME,
  s3Region: parsed.data.S3_REGION,
  s3AccessKeyId: parsed.data.S3_ACCESS_KEY_ID,
  s3SecretAccessKey: parsed.data.S3_SECRET_ACCESS_KEY,{{if .IsGarage}}
  s3Endpoint: parsed.data.S3_ENDPOINT,{{end}}
{{end}}  jwtSecret: parsed.data.JWT_SECRET,
  jwtExpiresIn: parsed.data.JWT_EXPIRES_IN,
  logLevel: parsed.data.LOG_LEVEL,
} as const;
`

// ─── config/swagger.config.ts ───
var TmplSwaggerConfig = `import { DocumentBuilder } from "@nestjs/swagger";
import { env } from "./env";

export const swaggerConfig = new DocumentBuilder()
  .setTitle(env.appName)
  .setDescription(` + "`" + `API documentation for ${env.appName}` + "`" + `)
  .addCookieAuth("cookie", {
    type: "apiKey",
    in: "cookie",
    name: "sessionId",
  })
  .addBearerAuth()
  .setVersion("1.0")
  .build();
`
