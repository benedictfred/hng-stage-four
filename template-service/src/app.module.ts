import { Module } from '@nestjs/common';
import { TemplatesModule } from './templates/templates.module';
import { HealthModule } from './health/health.module';

@Module({
  imports: [TemplatesModule, HealthModule],
})
export class AppModule {}
