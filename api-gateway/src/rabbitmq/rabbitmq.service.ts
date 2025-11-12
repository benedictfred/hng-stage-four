import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from "@nestjs/common";
import * as amqp from "amqplib";

@Injectable()
export class RabbitMQService implements OnModuleInit, OnModuleDestroy {
  private connection: amqp.Connection;
  private channel: amqp.Channel;
  private readonly logger = new Logger(RabbitMQService.name);

  async onModuleInit() {
    try {
      // Connect to RabbitMQ
      this.connection = await amqp.connect(
        process.env.RABBITMQ_URL || "amqp://localhost:5672",
      );
      this.channel = await this.connection.createChannel();

      // Declare exchange
      await this.channel.assertExchange("notifications.direct", "direct", {
        durable: true,
      });

      // Declare queues
      await this.channel.assertQueue("email.queue", { durable: true });
      await this.channel.assertQueue("push.queue", { durable: true });
      await this.channel.assertQueue("failed.queue", { durable: true });

      // Bind queues to exchange
      await this.channel.bindQueue(
        "email.queue",
        "notifications.direct",
        "email",
      );
      await this.channel.bindQueue(
        "push.queue",
        "notifications.direct",
        "push",
      );
      await this.channel.bindQueue(
        "failed.queue",
        "notifications.direct",
        "failed",
      );

      this.logger.log("✅ RabbitMQ connected and queues configured");
    } catch (error) {
      this.logger.error("❌ Failed to connect to RabbitMQ:", error);
      throw error;
    }
  }

  publishToQueue(
    routingKey: "email" | "push" | "failed",
    message: any,
  ): boolean {
    try {
      const messageBuffer = Buffer.from(JSON.stringify(message));

      this.channel.publish("notifications.direct", routingKey, messageBuffer, {
        persistent: true,
        contentType: "application/json",
        timestamp: Date.now(),
      });

      this.logger.log(`📤 Message published to ${routingKey} queue`);
      return true;
    } catch (error) {
      this.logger.error(`❌ Failed to publish to ${routingKey}:`, error);
      return false;
    }
  }

  async onModuleDestroy() {
    await this.channel?.close();
    await this.connection?.close();
    this.logger.log("RabbitMQ connection closed");
  }
}
