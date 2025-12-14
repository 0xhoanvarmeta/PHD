import {
  Controller,
  Get,
  Post,
  Body,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse as SwaggerResponse,
  ApiBearerAuth,
} from '@nestjs/swagger';
import { CommandsService } from './commands.service';
import { CommandPayload, ApiResponse } from '@phd/shared';
import { Public } from '../auth/decorators/public.decorator';

@ApiTags('Commands (Legacy)')
@Controller('commands')
export class CommandsController {
  constructor(private readonly commandsService: CommandsService) {}

  @Post('trigger')
  @ApiBearerAuth('JWT-auth')
  @ApiOperation({
    summary: 'Trigger a command on blockchain',
    description: 'Protected endpoint - requires admin JWT token',
  })
  @SwaggerResponse({ status: 200, description: 'Command triggered successfully' })
  @SwaggerResponse({ status: 401, description: 'Unauthorized' })
  @SwaggerResponse({ status: 500, description: 'Blockchain error' })
  async triggerCommand(@Body() payload: CommandPayload) {
    try {
      return await this.commandsService.triggerCommand(payload);
    } catch (error) {
      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR
      );
    }
  }

  @Public()
  @Get('current')
  @ApiOperation({ summary: 'Get current command from blockchain' })
  @SwaggerResponse({ status: 200, description: 'Current command retrieved' })
  async getCurrentCommand() {
    try {
      const command = await this.commandsService.getCurrentCommand();
      return {
        success: true,
        data: command,
      } as ApiResponse;
    } catch (error) {
      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR
      );
    }
  }

  @Public()
  @Get('latest-id')
  @ApiOperation({ summary: 'Get latest command ID' })
  @SwaggerResponse({ status: 200, description: 'Latest command ID retrieved' })
  async getLatestCommandId() {
    try {
      const commandId = await this.commandsService.getLatestCommandId();
      return {
        success: true,
        data: { commandId },
      } as ApiResponse;
    } catch (error) {
      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR
      );
    }
  }

  @Public()
  @Get('history')
  @ApiOperation({ summary: 'Get command history (last 100 events)' })
  @SwaggerResponse({ status: 200, description: 'Command history retrieved' })
  getCommandHistory() {
    const history = this.commandsService.getCommandHistory();
    return {
      success: true,
      data: history,
    } as ApiResponse;
  }

  @Public()
  @Get('info')
  @ApiOperation({ summary: 'Get blockchain provider info' })
  @SwaggerResponse({ status: 200, description: 'Blockchain info retrieved' })
  async getInfo() {
    try {
      const info = await this.commandsService.getBlockchainInfo();
      return {
        success: true,
        data: info,
      } as ApiResponse;
    } catch (error) {
      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR
      );
    }
  }
}
