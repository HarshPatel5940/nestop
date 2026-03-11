package internal

// ─── main.ts ───
var TmplMainTS = `{{if .UsesFastify}}import { NestFactory } from "@nestjs/core";
import {
  FastifyAdapter,
  type NestFastifyApplication,
} from "@nestjs/platform-fastify";
import compress from "@fastify/compress";
import multipart from "@fastify/multipart";
import fastifyCookie from "@fastify/cookie";
import cors from "@fastify/cors";
import { SwaggerModule } from "@nestjs/swagger";
import { Logger, LoggerErrorInterceptor } from "nestjs-pino";
import { AppModule } from "./app.module";
import { env } from "./config";
import { swaggerConfig } from "./config/swagger.config";
import { GlobalExceptionFilter } from "./filters/global-exception.filter";{{if .UsesPrisma}}
import { PrismaExceptionFilter } from "./filters/prisma-exception.filter";{{end}}

async function bootstrap() {
  const adapter = new FastifyAdapter({
    logger: false,
    trustProxy: true,
    bodyLimit: 50 * 1024 * 1024,
  });

  await adapter.register(compress, { encodings: ["gzip", "deflate"] });
  await adapter.register(multipart, {
    limits: { fileSize: 50 * 1024 * 1024, files: 1 },
    attachFieldsToBody: false,
  });

  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    adapter,
  );

  await app.register(cors, {
    origin: env.frontendUrl,
    credentials: true,
    methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization", "Accept"],
    exposedHeaders: ["Content-Disposition", "x-correlation-id"],
  });

  app.useLogger(app.get(Logger));

  await app.register(fastifyCookie, {
    secret: env.jwtSecret,
  });

  app.useGlobalFilters({{if .UsesPrisma}}
    new PrismaExceptionFilter(),{{end}}
    new GlobalExceptionFilter(),
  );
  app.useGlobalInterceptors(new LoggerErrorInterceptor());

  const document = SwaggerModule.createDocument(app, swaggerConfig);
  SwaggerModule.setup("api/swagger", app, document, {
    swaggerOptions: {
      persistAuthorization: true,
      withCredentials: true,
    },
  });

  await app.listen(env.port ?? 3000, "0.0.0.0", () => {
    console.log(` + "`" + `🚀 ${env.appName} running on port ${env.port}` + "`" + `);
  });
}

void bootstrap();
{{else}}import { NestFactory } from "@nestjs/core";
import { SwaggerModule } from "@nestjs/swagger";
import { Logger, LoggerErrorInterceptor } from "nestjs-pino";
import { AppModule } from "./app.module";
import { env } from "./config";
import { swaggerConfig } from "./config/swagger.config";
import { GlobalExceptionFilter } from "./filters/global-exception.filter";{{if .UsesPrisma}}
import { PrismaExceptionFilter } from "./filters/prisma-exception.filter";{{end}}

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.enableCors({
    origin: env.frontendUrl,
    credentials: true,
    methods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization", "Accept"],
    exposedHeaders: ["x-correlation-id"],
  });

  app.useLogger(app.get(Logger));

  app.useGlobalFilters({{if .UsesPrisma}}
    new PrismaExceptionFilter(),{{end}}
    new GlobalExceptionFilter(),
  );
  app.useGlobalInterceptors(new LoggerErrorInterceptor());

  const document = SwaggerModule.createDocument(app, swaggerConfig);
  SwaggerModule.setup("api/swagger", app, document, {
    swaggerOptions: {
      persistAuthorization: true,
      withCredentials: true,
    },
  });

  await app.listen(env.port ?? 3000, () => {
    console.log(` + "`" + `🚀 ${env.appName} running on port ${env.port}` + "`" + `);
  });
}

void bootstrap();
{{end}}`

// ─── app.module.ts ───
var TmplAppModule = `import { Module, MiddlewareConsumer, NestModule } from "@nestjs/common";
import { APP_GUARD } from "@nestjs/core";
import { LoggerModule } from "nestjs-pino";
import { ThrottlerModule, ThrottlerGuard } from "@nestjs/throttler";
import { AppController } from "./app.controller";
import { AppService } from "./app.service";
import { AuthGuard } from "./api/guards/auth.guard";
import { RolesGuard } from "./api/guards/roles.guard";
import { CorrelationIdMiddleware } from "./middleware/correlation-id.middleware";{{if .EnableRedis}}
import { RedisModule } from "./providers/redis/redis.module";{{end}}{{if .HasS3}}
import { S3Module } from "./providers/s3/s3.module";{{end}}{{if .UsesPrisma}}
import { PrismaModule } from "./database/prisma.module";{{end}}{{if .UsesMongoose}}
import { DatabaseModule } from "./database/database.module";{{end}}{{if .InitHealth}}
import { HealthModule } from "./api/routes/health/health.module";{{end}}{{if .InitAuth}}
import { AuthModule } from "./api/routes/auth/auth.module";{{end}}

@Module({
  imports: [
    LoggerModule.forRoot({
      pinoHttp: {
        customLogLevel: (_req, res, _err) => {
          if (res.statusCode >= 200 && res.statusCode < 300) {
            return "silent";
          }
          return "info";
        },
        transport:
          process.env.NODE_ENV !== "production"
            ? {
                target: "pino-pretty",
                options: {
                  colorize: true,
                  translateTime: "HH:MM:ss Z",
                  ignore: "hostname",
                },
              }
            : undefined,
        customSuccessMessage: (req, res) => {
          return ` + "`" + `${req.method} ${req.url} - ${res.statusCode}` + "`" + `;
        },
      },
    }),
    ThrottlerModule.forRoot([
      { name: "short", ttl: 1000, limit: 3 },
      { name: "medium", ttl: 10000, limit: 20 },
      { name: "long", ttl: 60000, limit: 100 },
    ]),{{if .EnableRedis}}
    RedisModule,{{end}}{{if .HasS3}}
    S3Module,{{end}}{{if .UsesPrisma}}
    PrismaModule,{{end}}{{if .UsesMongoose}}
    DatabaseModule,{{end}}{{if .InitHealth}}
    HealthModule,{{end}}{{if .InitAuth}}
    AuthModule,{{end}}
  ],
  controllers: [AppController],
  providers: [
    AppService,
    { provide: APP_GUARD, useClass: ThrottlerGuard },
    { provide: APP_GUARD, useClass: AuthGuard },
    { provide: APP_GUARD, useClass: RolesGuard },
  ],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(CorrelationIdMiddleware).forRoutes("*path");
  }
}
`

// ─── app.controller.ts ───
var TmplAppController = `import { Controller, Get } from "@nestjs/common";
import { AppService } from "./app.service";
import { PinoLogger } from "nestjs-pino";
import { Public } from "./api/decorators/public.decorator";

@Controller()
export class AppController {
  constructor(
    private readonly appService: AppService,
    private readonly logger: PinoLogger,
  ) {
    this.logger.setContext(AppController.name);
  }

  @Public()
  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

  @Public()
  @Get("health")
  getHealth() {
    return this.appService.getHealth();
  }
}
`

// ─── app.service.ts ───
var TmplAppService = `import { Injectable } from "@nestjs/common";
import { PinoLogger } from "nestjs-pino";

@Injectable()
export class AppService {
  constructor(private readonly logger: PinoLogger) {
    this.logger.setContext(AppService.name);
  }

  getHello(): string {
    return "Hello World!";
  }

  getHealth() {
    return {
      status: "ok",
      uptime: process.uptime(),
      timestamp: new Date().toISOString(),
    };
  }
}
`

// ─── app.controller.spec.ts ───
var TmplAppControllerSpec = `import { Test, TestingModule } from "@nestjs/testing";
import { AppController } from "./app.controller";
import { AppService } from "./app.service";
import { mockPinoLoggerProvider } from "../test/mock/pino-mock";

describe("AppController", () => {
  let appController: AppController;

  beforeEach(async () => {
    const app: TestingModule = await Test.createTestingModule({
      controllers: [AppController],
      providers: [AppService, mockPinoLoggerProvider],
    }).compile();

    appController = app.get<AppController>(AppController);
  });

  describe("root", () => {
    it("should return Hello World!", () => {
      expect(appController.getHello()).toBe("Hello World!");
    });
  });

  describe("health", () => {
    it("should return health status", () => {
      const result = appController.getHealth();
      expect(result.status).toBe("ok");
      expect(result.timestamp).toBeDefined();
    });
  });
});
`
