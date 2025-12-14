import {
  IsString,
  IsNotEmpty,
  IsOptional,
  IsEnum,
  IsObject,
  ValidateNested,
  MaxLength,
} from 'class-validator';
import { Type } from 'class-transformer';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { CommandType } from '@phd/shared';

class ScriptDataDto {
  @ApiProperty({
    description: 'Command type',
    enum: CommandType,
    example: CommandType.SCRIPT,
  })
  @IsEnum(CommandType)
  commandType: CommandType;

  @ApiProperty({
    description: 'Script content or URL',
    example: 'console.log("Hello World")',
  })
  @IsString()
  @IsNotEmpty()
  data: string;

  @ApiPropertyOptional({
    description: 'Additional metadata',
    example: { author: 'admin', version: '1.0' },
  })
  @IsObject()
  @IsOptional()
  metadata?: Record<string, any>;
}

export class CreateScriptDto {
  @ApiProperty({
    description: 'Script name',
    example: 'Hello World Script',
  })
  @IsString()
  @IsNotEmpty()
  @MaxLength(255)
  name: string;

  @ApiPropertyOptional({
    description: 'Script description',
    example: 'A simple hello world script',
  })
  @IsString()
  @IsOptional()
  @MaxLength(1000)
  description?: string;

  @ApiProperty({
    description: 'Script JSON data',
    type: ScriptDataDto,
  })
  @IsObject()
  @ValidateNested()
  @Type(() => ScriptDataDto)
  jsonData: ScriptDataDto;
}
