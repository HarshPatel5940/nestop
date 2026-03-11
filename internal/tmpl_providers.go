package internal

// ─── providers/redis/redis.module.ts ───
var TmplRedisModule = `import { Module, Global, Logger } from "@nestjs/common";
import { createClient, type RedisClientType } from "redis";
import { env } from "../../config";
import { RedisService } from "./redis.service";

@Global()
@Module({
  providers: [
    {
      provide: "REDIS_CLIENT",
      useFactory: async (): Promise<RedisClientType> => {
        const logger = new Logger("RedisModule");
        const maxRetries = 5;
        let retryCount = 0;

        const redisClient: RedisClientType = createClient({
          url: ` + "`" + `redis://${env.redisHost}:${env.redisPort}` + "`" + `,
          socket: {
            reconnectStrategy: (retries) => {
              if (retries >= maxRetries) {
                logger.error(
                  ` + "`" + `Failed to connect to Redis after ${maxRetries} attempts` + "`" + `,
                );
                return new Error("Max retries reached");
              }
              return Math.min(1000 * Math.pow(2, retries), 16000);
            },
          },
        });

        redisClient.on("connect", () => {
          logger.log("Connection to Redis established");
        });

        redisClient.on("error", (err) => {
          retryCount++;
          if (retryCount <= maxRetries) {
            logger.error(
              ` + "`" + `Redis connection error (attempt ${retryCount}/${maxRetries}):` + "`" + `,
              err,
            );
          }
        });

        await redisClient.connect();
        return redisClient;
      },
    },
    RedisService,
  ],
  exports: ["REDIS_CLIENT", RedisService],
})
export class RedisModule {}
`

// ─── providers/redis/redis.service.ts ───
var TmplRedisService = `import { Injectable, OnApplicationShutdown, Inject } from "@nestjs/common";
import { RedisClientType } from "redis";

@Injectable()
export class RedisService implements OnApplicationShutdown {
  constructor(
    @Inject("REDIS_CLIENT") private readonly redisClient: RedisClientType,
  ) {}

  async get(key: string): Promise<string | null> {
    return this.redisClient.get(key);
  }

  async set(key: string, value: string, ttlSeconds?: number): Promise<void> {
    if (ttlSeconds) {
      await this.redisClient.setEx(key, ttlSeconds, value);
    } else {
      await this.redisClient.set(key, value);
    }
  }

  async del(key: string): Promise<void> {
    await this.redisClient.del(key);
  }

  async onApplicationShutdown(): Promise<void> {
    if (this.redisClient.isOpen) {
      await this.redisClient.quit();
    }
  }
}
`

// ─── providers/s3/s3.module.ts ───
var TmplS3Module = `import { Module, Global } from "@nestjs/common";
import { S3Client } from "@aws-sdk/client-s3";
import { env } from "../../config";
import { S3Service } from "./s3.service";

@Global()
@Module({
  providers: [
    {
      provide: "S3_CLIENT",
      useFactory: (): S3Client => {
{{if .IsGarage}}        return new S3Client({
          region: env.s3Region,
          endpoint: env.s3Endpoint,
          forcePathStyle: true,
          credentials: {
            accessKeyId: env.s3AccessKeyId,
            secretAccessKey: env.s3SecretAccessKey,
          },
        });
{{else}}        return new S3Client({
          region: env.s3Region,
          credentials: {
            accessKeyId: env.s3AccessKeyId,
            secretAccessKey: env.s3SecretAccessKey,
          },
        });
{{end}}      },
    },
    S3Service,
  ],
  exports: ["S3_CLIENT", S3Service],
})
export class S3Module {}
`

// ─── providers/s3/s3.service.ts ───
var TmplS3Service = `import { Injectable, Inject } from "@nestjs/common";
import {
  S3Client,
  PutObjectCommand,
  GetObjectCommand,
  DeleteObjectCommand,
} from "@aws-sdk/client-s3";
import { Readable } from "stream";

@Injectable()
export class S3Service {
  constructor(@Inject("S3_CLIENT") private readonly s3Client: S3Client) {}

  async uploadFile(
    bucketName: string,
    key: string,
    body: Buffer,
    contentType?: string,
  ): Promise<void> {
    const command = new PutObjectCommand({
      Bucket: bucketName,
      Key: key,
      Body: body,
      ContentType: contentType,
    });
    await this.s3Client.send(command);
  }

  async getFile(bucketName: string, key: string): Promise<Readable> {
    const command = new GetObjectCommand({
      Bucket: bucketName,
      Key: key,
    });
    const response = await this.s3Client.send(command);
    if (response.Body instanceof Readable) {
      return response.Body;
    }
    throw new Error("Failed to retrieve file from S3");
  }

  async deleteFile(bucketName: string, key: string): Promise<void> {
    const command = new DeleteObjectCommand({
      Bucket: bucketName,
      Key: key,
    });
    await this.s3Client.send(command);
  }
}
`
