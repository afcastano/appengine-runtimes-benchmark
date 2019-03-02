import { Module } from '@nestjs/common';
import { ConfigurationModule } from '../config/config.module';
import { DummiesRepository } from './dummies.repository';
import { DummiesResolver } from './dummies.resolver';


@Module({
    imports: [ConfigurationModule],
    providers: [DummiesResolver, DummiesRepository]

})
export class DummiesModule {}