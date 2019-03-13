import { Module, OnModuleInit } from '@nestjs/common';
import { GCloudModule } from '@3wks/gae-node-nestjs';
import { ConfigurationModule } from './config/config.module';
import { UserModule } from './users/users.module';
import { DummiesModule } from './dummies/dummies.module';

@Module({
  imports: [
    GCloudModule.forConfiguration({
      configurationModule: ConfigurationModule,
      userModule: UserModule,
    }),
    ConfigurationModule,
    UserModule,
    DummiesModule
  ],
})
export class AppModule {}
