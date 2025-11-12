import {
  Controller,
  Post,
  Get,
  Body,
  Param,
  Query,
  ParseIntPipe,
  HttpCode,
  HttpStatus,
} from "@nestjs/common";
import { NotificationsService } from "./notifications.service";
import { SendNotificationDto } from "./dto/send-notification.dto";

@Controller("notifications")
export class NotificationsController {
  constructor(private readonly notificationsService: NotificationsService) {}

  @Post()
  @HttpCode(HttpStatus.ACCEPTED)
  async sendNotification(@Body() dto: SendNotificationDto) {
    return await this.notificationsService.sendNotification(dto);
  }

  @Get()
  getAllNotifications(
    @Query("page", new ParseIntPipe({ optional: true })) page: number = 1,
    @Query("limit", new ParseIntPipe({ optional: true })) limit: number = 10,
  ) {
    return this.notificationsService.getAllNotifications(page, limit);
  }

  @Get(":id")
  getNotificationStatus(@Param("id") notification_id: string) {
    return this.notificationsService.getNotificationStatus(notification_id);
  }
}
