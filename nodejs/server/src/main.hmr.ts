import { bootstrap } from './bootstrap';
import { rootLogger } from '@3wks/gae-node-nestjs';

declare const module: any;

bootstrap().then(({ stop }) => {
  rootLogger.info('Setting up hot reloading');

  if (module.hot) {
    module.hot.accept((ex: Error) => {
      rootLogger.error(
        'Unable to hot-reload due to error - killing process!',
        ex,
      );
      process.exit(1);
    });
    module.hot.dispose(() => stop());
  }
});
