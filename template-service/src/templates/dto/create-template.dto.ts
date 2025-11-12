import { IsString, IsArray, IsOptional } from 'class-validator';

export class CreateTemplateDto {
  @IsString()
  name: string;

  @IsOptional()
  @IsString()
  language?: string;

  @IsString()
  subject: string;

  @IsString()
  html_body: string;

  @IsOptional()
  @IsString()
  text_body?: string;

  @IsArray()
  variables: string[];
}
