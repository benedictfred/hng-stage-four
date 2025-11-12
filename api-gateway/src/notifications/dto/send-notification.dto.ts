import {
  IsString,
  IsNotEmpty,
  IsEnum,
  IsObject,
  IsOptional,
} from "class-validator";

export enum NotificationType {
  EMAIL = "email",
  PUSH = "push",
  BOTH = "both",
}

export class SendNotificationDto {
  @IsString()
  @IsNotEmpty()
  user_id: string;

  @IsEnum(NotificationType)
  @IsNotEmpty()
  type: NotificationType;

  @IsString()
  @IsNotEmpty()
  template_id: string;

  @IsObject()
  @IsOptional()
  data?: Record<string, any>;

  @IsString()
  @IsOptional()
  priority?: "high" | "normal" | "low";
}

export class NotificationResponse {
  notification_id: string;
  user_id: string;
  type: NotificationType;
  status: "queued" | "processing" | "sent" | "failed";
  created_at: Date;
}
