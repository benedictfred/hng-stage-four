import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { CreateTemplateDto } from './dto/create-template.dto';
import { Template } from 'generated/prisma';

export interface PaginationMeta {
  total: number;
  limit: number;
  page: number;
  total_pages: number;
  has_next: boolean;
  has_previous: boolean;
}

@Injectable()
export class TemplatesService {
  constructor(private prisma: PrismaService) {}

  async findOne(id: string): Promise<Template> {
    const template = await this.prisma.template.findUnique({
      where: { id },
    });

    if (!template) {
      throw new NotFoundException(`Template with ID ${id} not found`);
    }
    return template;
  }

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

  async findAll(
    page = 1,
    limit = 10,
  ): Promise<{ data: Template[]; meta: PaginationMeta }> {
    const skip = (page - 1) * limit;
    const [data, total] = await Promise.all([
      this.prisma.template.findMany({
        skip,
        take: limit,
        orderBy: { created_at: 'desc' },
      }),
      this.prisma.template.count(),
    ]);

    const total_pages = Math.ceil(total / limit);

    return {
      data,
      meta: {
        total,
        limit,
        page,
        total_pages,
        has_next: page < total_pages,
        has_previous: page > 1,
      },
    };
  }
}
