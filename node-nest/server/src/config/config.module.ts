import { Module } from '@nestjs/common';
import { ConfigurationProvider } from './config.provider';

export const configurationProvider = new ConfigurationProvider();

@Module({
  providers: [
    { provide: ConfigurationProvider, useValue: configurationProvider },
    { provide: 'Configuration', useValue: configurationProvider },
  ],
  exports: [ConfigurationProvider, 'Configuration'],
})
export class ConfigurationModule {}
