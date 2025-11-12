import { Controller, Get } from "@nestjs/common";

@Controller("health")
export class HealthController {
  @Get()
  check() {
    return {
      success: true,
      message: "API Gateway is healthy",
      data: {
        status: "healthy",
        service: "api-gateway",
        timestamp: new Date().toISOString(),
      },
    };
  }
}
