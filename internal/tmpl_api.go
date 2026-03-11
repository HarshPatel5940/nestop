package internal

// ─── api/constants/index.ts ───
var TmplConstantsIndex = `export * from "./errors";
export * from "./response";
`

// ─── api/constants/errors.ts ───
var TmplConstantsErrors = `export type ReturnError = {
  status: number;
  message: string;
  prettyMessage: string;
};

export const errors: Record<string, ReturnError> = {
  UNAUTHORIZED: {
    status: 401,
    message: "Unauthorized",
    prettyMessage: "You are not authorized to access this resource.",
  },
  FORBIDDEN: {
    status: 403,
    message: "Forbidden",
    prettyMessage: "You are forbidden from accessing this resource.",
  },
  NOT_FOUND: {
    status: 404,
    message: "Not Found",
    prettyMessage: "The resource you are looking for does not exist.",
  },
  CONFLICT: {
    status: 409,
    message: "Conflict",
    prettyMessage: "The resource you are trying to create already exists.",
  },
  INTERNAL_SERVER_ERROR: {
    status: 500,
    message: "Internal Server Error",
    prettyMessage: "An unexpected error occurred.",
  },
  BAD_REQUEST: {
    status: 400,
    message: "Bad Request",
    prettyMessage: "The request you made was invalid.",
  },
  INVALID_CREDENTIALS: {
    status: 401,
    message: "Invalid Credentials",
    prettyMessage: "The email or password you entered is incorrect.",
  },
};
`

// ─── api/constants/response.ts ───
var TmplConstantsResponse = `export type ReturnResponse = {
  status: number;
  message: string;
  prettyMessage?: string;
  data?: any;
};
`

// ─── api/guards/index.ts ───
var TmplGuardsIndex = `export * from "./auth.guard";
export * from "./roles.guard";
`

// ─── api/guards/auth.guard.ts ───
// Uses raw jsonwebtoken via jwt.utils — no @nestjs/jwt needed
var TmplAuthGuard = `import {
  Injectable,
  CanActivate,
  ExecutionContext,
  UnauthorizedException,
} from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { IS_PUBLIC_KEY } from "../decorators/public.decorator";
import { verifyJwt } from "../utils/jwt.utils";
import { errors } from "../constants";

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);

    if (isPublic) return true;

    const request = context.switchToHttp().getRequest();
    const token = this.extractToken(request);

    if (!token) {
      throw new UnauthorizedException(errors.UNAUTHORIZED);
    }

    try {
      const payload = verifyJwt(token);
      request.user = payload;
    } catch {
      throw new UnauthorizedException(errors.UNAUTHORIZED);
    }

    return true;
  }

  private extractToken(request: any): string | undefined {
    const [type, token] =
      request.headers.authorization?.split(" ") ?? [];
    return type === "Bearer" ? token : undefined;
  }
}
`

// ─── api/interceptors/index.ts ───
var TmplInterceptorsIndex = `export * from "./response.interceptor";
`

// ─── api/interceptors/response.interceptor.ts ───
var TmplResponseInterceptor = `import {
  Injectable,
  NestInterceptor,
  ExecutionContext,
  CallHandler,
} from "@nestjs/common";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

@Injectable()
export class ResponseInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    return next.handle().pipe(
      map((response) => {
        const httpResponse = context.switchToHttp().getResponse();
        if (response?.status) {
          httpResponse.status(response.status);
        }
        return response;
      }),
    );
  }
}
`

// ─── api/utils/index.ts ───
var TmplUtilsIndex = `export * from "./hash.utils";
export * from "./ulid.utils";
export * from "./jwt.utils";
export * from "./pagination.utils";
`

// ─── api/utils/hash.utils.ts ───
var TmplHashUtils = `import * as argon2 from "argon2";

export const hashPassword = async (password: string): Promise<string> => {
  return argon2.hash(password);
};

export const comparePassword = async (
  password: string,
  hashedPassword: string,
): Promise<boolean> => {
  return argon2.verify(hashedPassword, password);
};
`

// ─── api/utils/ulid.utils.ts ───
var TmplUlidUtils = `import { z } from "zod";
import { ulid, decodeTime, type ULID } from "ulid";

const isULID = (value: string): boolean => {
  if (typeof value !== "string") return false;
  if (value.length !== 26) return false;
  return /^[0-7][0-9A-HJKMNP-TV-Z]{25}$/.test(value);
};

export const createULIDSchema = () =>
  z
    .string()
    .refine(isULID, { message: "Invalid ULID format" })
    .refine(
      (val) => decodeTime(val) <= Date.now(),
      { message: "ULID timestamp cannot be in the future" },
    ) as unknown as z.ZodType<ULID>;

export const generateULID = (): string => ulid();
`
