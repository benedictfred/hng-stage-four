import { Module } from "@nestjs/common";
import { HttpModule } from "@nestjs/axios";
import { NotificationsController } from "./notifications.controller";
import { NotificationsService } from "./notifications.service";
import { RabbitMQModule } from "../rabbitmq/rabbitmq.module";

@Module({
  imports: [
    HttpModule.register({
      timeout: 5000,
      maxRedirects: 5,
    }),
    RabbitMQModule,
  ],
  controllers: [NotificationsController],
  providers: [NotificationsService],
})
export class NotificationsModule {}
