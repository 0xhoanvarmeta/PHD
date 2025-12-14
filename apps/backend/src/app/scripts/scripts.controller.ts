import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  Query,
  ParseUUIDPipe,
  HttpCode,
  HttpStatus,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
} from '@nestjs/swagger';
import { ScriptsService } from './scripts.service';
import { CreateScriptDto } from './dto/create-script.dto';
import { UpdateScriptDto } from './dto/update-script.dto';
import { QueryScriptsDto } from './dto/query-scripts.dto';
import { CurrentAdmin } from '../auth/decorators/current-admin.decorator';
import { Admin } from '../admin/entities/admin.entity';

@ApiTags('Scripts')
@ApiBearerAuth('JWT-auth')
@Controller('scripts')
export class ScriptsController {
  constructor(private readonly scriptsService: ScriptsService) {}

  @Post()
  @ApiOperation({ summary: 'Create a new script' })
  @ApiResponse({ status: 201, description: 'Script created successfully' })
  @ApiResponse({ status: 400, description: 'Invalid input' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  create(
    @Body() createScriptDto: CreateScriptDto,
    @CurrentAdmin() admin: Admin
  ) {
    return this.scriptsService.create(createScriptDto, admin);
  }

  @Get()
  @ApiOperation({ summary: 'Get all scripts with pagination' })
  @ApiResponse({ status: 200, description: 'Scripts retrieved successfully' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  findAll(@Query() query: QueryScriptsDto) {
    return this.scriptsService.findAll(query);
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get a script by ID' })
  @ApiResponse({ status: 200, description: 'Script found' })
  @ApiResponse({ status: 404, description: 'Script not found' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  findOne(@Param('id', ParseUUIDPipe) id: string) {
    return this.scriptsService.findOne(id);
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Update a script' })
  @ApiResponse({ status: 200, description: 'Script updated successfully' })
  @ApiResponse({ status: 403, description: 'Forbidden - not script owner' })
  @ApiResponse({ status: 404, description: 'Script not found' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  update(
    @Param('id', ParseUUIDPipe) id: string,
    @Body() updateScriptDto: UpdateScriptDto,
    @CurrentAdmin() admin: Admin
  ) {
    return this.scriptsService.update(id, updateScriptDto, admin);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiOperation({ summary: 'Delete a script' })
  @ApiResponse({ status: 204, description: 'Script deleted successfully' })
  @ApiResponse({ status: 403, description: 'Forbidden - not script owner' })
  @ApiResponse({ status: 404, description: 'Script not found' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  async remove(
    @Param('id', ParseUUIDPipe) id: string,
    @CurrentAdmin() admin: Admin
  ) {
    await this.scriptsService.remove(id, admin);
  }

  @Post(':id/trigger')
  @ApiOperation({ summary: 'Trigger a script on blockchain' })
  @ApiResponse({
    status: 200,
    description: 'Script triggered successfully',
    schema: {
      example: {
        success: true,
        transactionHash: '0x123...',
        script: {
          id: '123e4567-e89b-12d3-a456-426614174000',
          name: 'Hello World Script',
          commandType: 0,
        },
      },
    },
  })
  @ApiResponse({ status: 404, description: 'Script not found' })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  @ApiResponse({ status: 500, description: 'Blockchain error' })
  trigger(
    @Param('id', ParseUUIDPipe) id: string,
    @CurrentAdmin() admin: Admin
  ) {
    return this.scriptsService.trigger(id, admin);
  }
}
