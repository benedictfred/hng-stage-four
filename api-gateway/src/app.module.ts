import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { APP_INTERCEPTOR } from "@nestjs/core";
import { NotificationsModule } from "./notifications/notifications.module";
import { RabbitMQModule } from "./rabbitmq/rabbitmq.module";
import { HealthModule } from "./health/health.module";
import { ResponseInterceptor } from "./common/interceptors/response.interceptor";

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    NotificationsModule,
    RabbitMQModule,
    HealthModule,
  ],
  providers: [
    {
      provide: APP_INTERCEPTOR,
      useClass: ResponseInterceptor,
    },
  ],
})
export class AppModule {}
