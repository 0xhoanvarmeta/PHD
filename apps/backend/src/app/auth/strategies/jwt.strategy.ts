import { Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { AdminService } from '../../admin/admin.service';

export interface JwtPayload {
  sub: string; // admin id
  username: string;
  email: string;
}

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(
    private configService: ConfigService,
    private adminService: AdminService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: configService.get<string>(
        'JWT_SECRET',
        'your-secret-key-change-in-production'
      ),
    });
  }

  async validate(payload: JwtPayload) {
    const admin = await this.adminService.findById(payload.sub);

    if (!admin || !admin.isActive) {
      throw new UnauthorizedException('Unauthorized access');
    }

    return admin;
  }
}
