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

// ─── Health Tests (Vitest) ───

var TmplHealthServiceVitestSpec = `import { describe, it, expect, beforeEach, vi } from "vitest";
import { Test, type TestingModule } from "@nestjs/testing";
import { HealthService } from "./health.service";
import { mockPinoLoggerProvider } from "../../../../test/mock/pino-mock";

describe("HealthService", () => {
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [HealthService, mockPinoLoggerProvider],
    }).compile();

    service = module.get<HealthService>(HealthService);
  });

  describe("check()", () => {
    it("should return status ok", () => {
      const result = service.check();
      expect(result.status).toBe("ok");
    });

    it("should include uptime as a number", () => {
      const result = service.check();
      expect(typeof result.uptime).toBe("number");
      expect(result.uptime).toBeGreaterThanOrEqual(0);
    });

    it("should include a valid ISO timestamp", () => {
      const result = service.check();
      expect(result.timestamp).toBeDefined();
      expect(() => new Date(result.timestamp)).not.toThrow();
      expect(new Date(result.timestamp).toISOString()).toBe(result.timestamp);
    });

    it("should include memory usage", () => {
      const result = service.check();
      expect(result.memory).toBeDefined();
      expect(typeof result.memory.heapUsed).toBe("number");
    });
  });

  describe("live()", () => {
    it("should return status ok", () => {
      expect(service.live()).toEqual({ status: "ok" });
    });
  });

  describe("ready()", () => {
    it("should return status ok with timestamp", () => {
      const result = service.ready();
      expect(result.status).toBe("ok");
      expect(result.timestamp).toBeDefined();
    });
  });
});
`

var TmplHealthControllerVitestSpec = `import { describe, it, expect, beforeEach, vi } from "vitest";
import { Test, type TestingModule } from "@nestjs/testing";
import { HealthController } from "./health.controller";
import { HealthService } from "./health.service";
import { mockPinoLoggerProvider } from "../../../../test/mock/pino-mock";

const mockHealthResult = {
  status: "ok",
  uptime: 100,
  timestamp: new Date().toISOString(),
  memory: { heapUsed: 1024, heapTotal: 2048, rss: 4096, external: 0, arrayBuffers: 0 },
  version: "0.1.0",
};

describe("HealthController", () => {
  let controller: HealthController;
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [HealthService, mockPinoLoggerProvider],
    }).compile();

    controller = module.get<HealthController>(HealthController);
    service = module.get<HealthService>(HealthService);
  });

  describe("check()", () => {
    it("should return health check result", () => {
      vi.spyOn(service, "check").mockReturnValue(mockHealthResult);
      const result = controller.check();
      expect(result).toEqual(mockHealthResult);
      expect(service.check).toHaveBeenCalledTimes(1);
    });

    it("should return status ok directly", () => {
      const result = controller.check();
      expect(result.status).toBe("ok");
    });
  });

  describe("live()", () => {
    it("should return liveness status", () => {
      vi.spyOn(service, "live").mockReturnValue({ status: "ok" });
      const result = controller.live();
      expect(result).toEqual({ status: "ok" });
    });
  });

  describe("ready()", () => {
    it("should return readiness status with timestamp", () => {
      const result = controller.ready();
      expect(result.status).toBe("ok");
      expect(result.timestamp).toBeDefined();
    });
  });
});
`

// ─── Health Tests (Jest) ───

var TmplHealthServiceJestSpec = `import { Test, TestingModule } from "@nestjs/testing";
import { HealthService } from "./health.service";
import { mockPinoLoggerProvider } from "../../../../test/mock/pino-mock";

describe("HealthService", () => {
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [HealthService, mockPinoLoggerProvider],
    }).compile();

    service = module.get<HealthService>(HealthService);
  });

  describe("check()", () => {
    it("should return status ok", () => {
      const result = service.check();
      expect(result.status).toBe("ok");
    });

    it("should include uptime as a number", () => {
      const result = service.check();
      expect(typeof result.uptime).toBe("number");
      expect(result.uptime).toBeGreaterThanOrEqual(0);
    });

    it("should include a valid ISO timestamp", () => {
      const result = service.check();
      expect(result.timestamp).toBeDefined();
      expect(() => new Date(result.timestamp)).not.toThrow();
      expect(new Date(result.timestamp).toISOString()).toBe(result.timestamp);
    });

    it("should include memory usage", () => {
      const result = service.check();
      expect(result.memory).toBeDefined();
      expect(typeof result.memory.heapUsed).toBe("number");
    });
  });

  describe("live()", () => {
    it("should return status ok", () => {
      expect(service.live()).toEqual({ status: "ok" });
    });
  });

  describe("ready()", () => {
    it("should return status ok with timestamp", () => {
      const result = service.ready();
      expect(result.status).toBe("ok");
      expect(result.timestamp).toBeDefined();
    });
  });
});
`

var TmplHealthControllerJestSpec = `import { Test, TestingModule } from "@nestjs/testing";
import { HealthController } from "./health.controller";
import { HealthService } from "./health.service";
import { mockPinoLoggerProvider } from "../../../../test/mock/pino-mock";

const mockHealthResult = {
  status: "ok",
  uptime: 100,
  timestamp: new Date().toISOString(),
  memory: { heapUsed: 1024, heapTotal: 2048, rss: 4096, external: 0, arrayBuffers: 0 },
  version: "0.1.0",
};

describe("HealthController", () => {
  let controller: HealthController;
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [HealthService, mockPinoLoggerProvider],
    }).compile();

    controller = module.get<HealthController>(HealthController);
    service = module.get<HealthService>(HealthService);
  });

  describe("check()", () => {
    it("should return health check result", () => {
      jest.spyOn(service, "check").mockReturnValue(mockHealthResult);
      const result = controller.check();
      expect(result).toEqual(mockHealthResult);
      expect(service.check).toHaveBeenCalledTimes(1);
    });

    it("should return status ok directly", () => {
      const result = controller.check();
      expect(result.status).toBe("ok");
    });
  });

  describe("live()", () => {
    it("should return liveness status", () => {
      jest.spyOn(service, "live").mockReturnValue({ status: "ok" });
      const result = controller.live();
      expect(result).toEqual({ status: "ok" });
    });
  });

  describe("ready()", () => {
    it("should return readiness status with timestamp", () => {
      const result = controller.ready();
      expect(result.status).toBe("ok");
      expect(result.timestamp).toBeDefined();
    });
  });
});
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
