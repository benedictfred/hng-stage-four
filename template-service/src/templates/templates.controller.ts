import { Body, Post, Controller } from '@nestjs/common';
import { TemplatesService } from './templates.service';
import { CreateTemplateDto } from './dto/create-template.dto';

@Controller('templates')
export class TemplatesController {
  constructor(private readonly service: TemplatesService) {}

  @Post()
  async create(@Body() dto: CreateTemplateDto) {
    const data = await this.service.create(dto);
    return {
      success: true,
      data,
      message: 'Template created successfully',
    };
  }
}
