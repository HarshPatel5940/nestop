package internal

// ─── pipes/index.ts ───
var TmplPipesIndex = `export * from "./zod.pipe";
`

// ─── pipes/zod.pipe.ts ───
var TmplZodPipe = `import { HttpException, HttpStatus, type PipeTransform } from "@nestjs/common";
import type { ZodObject, ZodRawShape } from "zod";

export class ZodValidationPipe implements PipeTransform {
  private readonly schema;

  constructor(schema: ZodObject<ZodRawShape>) {
    this.schema = schema;
  }

  transform(value: unknown) {
    try {
      this.schema.parse(value);
    } catch (error) {
      throw new HttpException(
        {
          statusCode: HttpStatus.BAD_REQUEST,
          error: "Validation Error",
          message: error,
        },
        HttpStatus.BAD_REQUEST,
      );
    }
    return value;
  }
}
`

// ─── filters/global-exception.filter.ts ───
var TmplGlobalExceptionFilter = `import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
  Logger,
} from "@nestjs/common";
{{if .UsesFastify}}import { FastifyReply, FastifyRequest } from "fastify";{{else}}import { Request, Response } from "express";{{end}}

@Catch()
export class GlobalExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(GlobalExceptionFilter.name);

  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
{{if .UsesFastify}}    const reply = ctx.getResponse<FastifyReply>();
    const request = ctx.getRequest<FastifyRequest>();
{{else}}    const reply = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();
{{end}}
    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = "Internal server error";
    let error = "Internal Server Error";

    if (exception instanceof HttpException) {
      status = exception.getStatus();
      const response = exception.getResponse();

      if (typeof response === "string") {
        message = response;
      } else if (typeof response === "object" && response !== null) {
        message = (response as any).message || message;
        error = (response as any).error || error;
      }

      if (status >= 500) {
        this.logger.error(
          ` + "`" + `[${request.method}] ${request.url} - ${status} ${message}` + "`" + `,
        );
      }
    } else if (exception instanceof Error) {
      message = exception.message;
      this.logger.error(
        ` + "`" + `Unhandled exception: ${exception.message}` + "`" + `,
        exception.stack,
      );
    } else {
      this.logger.error("Unknown exception occurred", exception);
    }

{{if .UsesFastify}}    reply.status(status).send({
{{else}}    reply.status(status).json({
{{end}}      statusCode: status,
      message,
      error,
      timestamp: new Date().toISOString(),
      path: request.url,
    });
  }
}
`

// ─── filters/prisma-exception.filter.ts ───
var TmplPrismaExceptionFilter = `import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpStatus,
  Logger,
} from "@nestjs/common";
{{if .UsesFastify}}import { FastifyReply } from "fastify";{{else}}import { Response } from "express";{{end}}
import { PrismaClientKnownRequestError } from "@prisma/client/runtime/library";

@Catch(PrismaClientKnownRequestError)
export class PrismaExceptionFilter implements ExceptionFilter {
  private readonly logger = new Logger(PrismaExceptionFilter.name);

  catch(exception: PrismaClientKnownRequestError, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
{{if .UsesFastify}}    const reply = ctx.getResponse<FastifyReply>();
{{else}}    const reply = ctx.getResponse<Response>();
{{end}}
    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = "Internal server error";

    switch (exception.code) {
      case "P2002":
        status = HttpStatus.CONFLICT;
        message = "Duplicate entry found";
        break;
      case "P2003":
        status = HttpStatus.BAD_REQUEST;
        message = "Invalid reference provided";
        break;
      case "P2025":
        status = HttpStatus.NOT_FOUND;
        message = "Record not found";
        break;
      case "P2021":
        status = HttpStatus.INTERNAL_SERVER_ERROR;
        message = "Database connection error";
        this.logger.error("Database connection failed", exception);
        break;
      default:
        this.logger.error(` + "`" + `Unhandled Prisma error: ${exception.code}` + "`" + `, {
          code: exception.code,
          message: exception.message,
          meta: exception.meta,
        });
        message = "Database operation failed";
        break;
    }

{{if .UsesFastify}}    reply.status(status).send({
{{else}}    reply.status(status).json({
{{end}}      statusCode: status,
      message,
      timestamp: new Date().toISOString(),
    });
  }
}
`
