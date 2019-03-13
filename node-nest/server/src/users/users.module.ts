import { Module } from '@nestjs/common';
import { UsersResolver } from './users.resolver';
import { UserRepository } from './users.repository';
import { UsersService } from './users.service';
import { ConfigurationModule } from '../config/config.module';
import { USER_SERVICE } from '@3wks/gae-node-nestjs';

@Module({
  imports: [ConfigurationModule],
  providers: [
    UsersResolver,
    UserRepository,
    UsersService,
    { provide: USER_SERVICE, useClass: UsersService },
  ],
  exports: [UserRepository, USER_SERVICE, UsersService],
})
export class UserModule {}
