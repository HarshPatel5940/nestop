package internal

// ─── Health Module ───

var TmplHealthModule = `import { Module } from "@nestjs/common";
import { HealthController } from "./health.controller";
import { HealthService } from "./health.service";

@Module({
  controllers: [HealthController],
  providers: [HealthService],
})
export class HealthModule {}
`

var TmplHealthController = `import { Controller, Get } from "@nestjs/common";
import { ApiTags, ApiOperation, ApiResponse } from "@nestjs/swagger";
import { HealthService } from "./health.service";
import { Public } from "../../decorators/public.decorator";

@ApiTags("health")
@Public()
@Controller("health")
export class HealthController {
  constructor(private readonly healthService: HealthService) {}

  @ApiOperation({ summary: "Health check" })
  @ApiResponse({ status: 200, description: "Service is healthy" })
  @Get()
  check() {
    return this.healthService.check();
  }

  @ApiOperation({ summary: "Liveness probe" })
  @Get("live")
  live() {
    return this.healthService.live();
  }

  @ApiOperation({ summary: "Readiness probe" })
  @Get("ready")
  ready() {
    return this.healthService.ready();
  }
}
`

var TmplHealthService = `import { Injectable } from "@nestjs/common";
import { PinoLogger } from "nestjs-pino";

@Injectable()
export class HealthService {
  constructor(private readonly logger: PinoLogger) {
    this.logger.setContext(HealthService.name);
  }

  check() {
    return {
      status: "ok",
      uptime: process.uptime(),
      timestamp: new Date().toISOString(),
      memory: process.memoryUsage(),
      version: process.env.npm_package_version ?? "unknown",
    };
  }

  live() {
    return { status: "ok" };
  }

  ready() {
    return { status: "ok", timestamp: new Date().toISOString() };
  }
}
`

// ─── Auth Module ───

var TmplAuthModule = `import { Module } from "@nestjs/common";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";

@Module({
  controllers: [AuthController],
  providers: [AuthService],
})
export class AuthModule {}
`

var TmplAuthController = `import { Controller, Post, Get, Body } from "@nestjs/common";
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
} from "@nestjs/swagger";
import { AuthService } from "./auth.service";
import { RegisterDto, LoginDto } from "./dto";
import { Public } from "../../decorators/public.decorator";
import { CurrentUser } from "../../decorators/current-user.decorator";
import type { JwtPayload } from "../../../api/utils/jwt.utils";

@ApiTags("auth")
@Controller("auth")
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @ApiOperation({ summary: "Register a new user" })
  @ApiResponse({ status: 201, description: "User registered successfully" })
  @ApiResponse({ status: 409, description: "Email already in use" })
  @Public()
  @Post("register")
  register(@Body() body: RegisterDto) {
    return this.authService.register(body);
  }

  @ApiOperation({ summary: "Login with credentials" })
  @ApiResponse({ status: 200, description: "Login successful, returns JWT" })
  @ApiResponse({ status: 401, description: "Invalid credentials" })
  @Public()
  @Post("login")
  login(@Body() body: LoginDto) {
    return this.authService.login(body);
  }

  @ApiOperation({ summary: "Get current session info" })
  @ApiResponse({ status: 200, description: "Session is valid" })
  @ApiResponse({ status: 401, description: "Unauthorized" })
  @ApiBearerAuth()
  @Get("session")
  session(@CurrentUser() user: JwtPayload) {
    return this.authService.checkSession(user);
  }
}
`

var TmplAuthService = `import { Injectable } from "@nestjs/common";
import { PinoLogger } from "nestjs-pino";
import { hashPassword, comparePassword, signJwt } from "../../../api/utils";
import { errors, type ReturnError, type ReturnResponse } from "../../../api/constants";
import type { RegisterDto, LoginDto } from "./dto";

@Injectable()
export class AuthService {
  constructor(private readonly logger: PinoLogger) {
    this.logger.setContext(AuthService.name);
  }

  async register(body: RegisterDto): Promise<ReturnResponse | ReturnError> {
    const { name, email, password } = body;

    // TODO: check if user already exists (query your database here)
    // const existing = await this.db.findUser({ email });
    // if (existing) return errors.CONFLICT;

    const hashedPassword = await hashPassword(password);

    // TODO: save user to database
    // await this.db.createUser({ name, email, password: hashedPassword });

    this.logger.info({ email }, "User registered");

    return {
      status: 201,
      message: "OK",
      prettyMessage: "User created successfully",
    };
  }

  async login(body: LoginDto): Promise<ReturnResponse | ReturnError> {
    const { email, password } = body;

    // TODO: fetch user from database
    // const user = await this.db.findUser({ email });
    // if (!user) return errors.INVALID_CREDENTIALS;

    // TODO: verify password
    // const valid = await comparePassword(password, user.password);
    // if (!valid) return errors.INVALID_CREDENTIALS;

    const token = signJwt({ sub: "replace-with-real-user-id", email });

    this.logger.info({ email }, "User logged in");

    return {
      status: 200,
      message: "OK",
      prettyMessage: "User logged in successfully",
      data: { token },
    };
  }

  async checkSession(user: any): Promise<ReturnResponse | ReturnError> {
    if (!user) return errors.UNAUTHORIZED;

    return {
      status: 200,
      message: "OK",
      prettyMessage: "Session is valid",
      data: { user },
    };
  }
}
`

// ─── Auth DTOs ───

var TmplAuthDTOs = `import { ApiProperty } from "@nestjs/swagger";
import { z } from "zod";

export const registerSchema = z
  .object({
    name: z.string().min(2).max(100),
    email: z.string().email(),
    password: z.string().min(8).max(100),
  })
  .strict();

export const loginSchema = z
  .object({
    email: z.string().email(),
    password: z.string().min(1),
  })
  .strict();

export type RegisterDto = z.infer<typeof registerSchema>;
export type LoginDto = z.infer<typeof loginSchema>;

export class RegisterDtoClass implements RegisterDto {
  @ApiProperty({ example: "John Doe" })
  name!: string;

  @ApiProperty({ example: "john@example.com" })
  email!: string;

  @ApiProperty({ example: "super-secret-password" })
  password!: string;
}

export class LoginDtoClass implements LoginDto {
  @ApiProperty({ example: "john@example.com" })
  email!: string;

  @ApiProperty({ example: "super-secret-password" })
  password!: string;
}
`
