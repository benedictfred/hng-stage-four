import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { CreateTemplateDto } from './dto/create-template.dto';
import { Template } from 'generated/prisma';

@Injectable()
export class TemplatesService {
  constructor(private prisma: PrismaService) {}

  async create(dto: CreateTemplateDto): Promise<Template> {
    const latest: Template | null = await this.prisma.template.findFirst({
      where: { name: dto.name, language: dto.language || 'en' },
      orderBy: { version: 'desc' },
    });

    const version: number = (latest?.version ?? 0) + 1;

    return this.prisma.template.create({
      data: {
        ...dto,
        language: dto.language || 'en',
        version,
        is_active: true,
      },
    });
  }
}
