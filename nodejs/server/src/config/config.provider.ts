import * as t from 'io-ts';
import { reporter } from 'io-ts-reporters';
import { ThrowReporter } from 'io-ts/lib/ThrowReporter';
import * as _ from 'lodash';
import * as Logger from 'bunyan';
import { Configuration, createLogger } from '@3wks/gae-node-nestjs';

const auth = t.partial({
  local: t.interface({
    enabled: t.boolean,
  }),
  google: t.interface({
    enabled: t.boolean,
    clientId: t.string,
    secret: t.string,
    signUpEnabled: t.boolean,
    signUpDomains: t.array(t.string),
    signUpRoles: t.array(t.string),
  }),
  saml: t.interface({
    enabled: t.boolean,
    cert: t.string,
    identityProviderUrl: t.string,
  }),
});

const devHooks = t.partial({
  disableLocalMailLogger: t.boolean,
  divertEmailTo: t.array(t.string),
  emailSubjectPrefix: t.string,
});

// tslint:disable-next-line:variable-name
const Config = t.intersection([
  t.interface({
    projectId: t.string,
    host: t.string,
    bucket: t.string,
    location: t.string,
    gmailUser: t.string,
    systemSecret: t.string,
    cookieSecret: t.string,
    auth,
  }),
  t.partial({
    APP_ENGINE_ENVIRONMENT: t.string,
    GOOGLE_CLOUD_PROJECT: t.string,
    devHooks,
    bootstrap: t.boolean,
    apiEndpoint: t.string,
    twilioNumber: t.string,
    twilioAccountSID: t.string,
    twilioAuthToken: t.string,
  }),
]);

interface SessionConfiguration {
  apiEndpoint?: string;
  projectId?: string;
  secret: string;
}

export class ConfigurationProvider implements Configuration {
  configuration: t.TypeOf<typeof Config>;
  logger: Logger;

  constructor() {
    this.logger = createLogger('configuration-provider');

    if (process.env.GOOGLE_CLOUD_PROJECT) {
      const projectId = process.env.GOOGLE_CLOUD_PROJECT;
      process.env.NODE_CONFIG_ENV = _.last(projectId.split('-'));
    } else if (!process.env.NODE_CONFIG_ENV) {
      process.env.NODE_CONFIG_ENV = 'development';
    }

    const nodeConfig = require('config');
    const mergedConfig: object = {};
    const configSources: any = nodeConfig.util.getConfigSources();

    configSources.forEach((config: any) => {
      this.logger.info(`Loading config from ${config.name}`);
      nodeConfig.util.extendDeep(mergedConfig, config.parsed);
    });

    const withEnvironment = nodeConfig.util.extendDeep(
      mergedConfig,
      process.env,
    );

    const decodedConfig = Config.decode(withEnvironment);

    if (decodedConfig.isLeft()) {
      this.logger.error(reporter(decodedConfig));
      throw ThrowReporter.report(decodedConfig);
    }

    this.configuration = decodedConfig.value;
  }

  get projectId(): string {
    return this.configuration.projectId;
  }

  get environment(): 'development' | 'appengine' {
    if (this.configuration.APP_ENGINE_ENVIRONMENT) {
      return 'appengine';
    }

    return 'development';
  }

  isDevelopment(): boolean {
    return this.environment === 'development';
  }

  get host(): string {
    return this.configuration.host;
  }

  get bucket(): string {
    return this.configuration.bucket;
  }

  get bootstrap(): boolean {
    return this.configuration.bootstrap === true;
  }

  get location(): string {
    return this.configuration.location;
  }

  get apiEndpoint(): string | undefined {
    return this.configuration.apiEndpoint;
  }

  get gmailUser(): string {
    return this.configuration.gmailUser;
  }

  get systemSecret(): Buffer {
    return Buffer.from(this.configuration.systemSecret, 'base64');
  }

  get auth() {
    return this.configuration.auth;
  }

  get devHooks() {
    return this.configuration.devHooks;
  }

  get session(): SessionConfiguration {
    return {
      secret: this.configuration.cookieSecret,
      apiEndpoint: this.apiEndpoint,
      projectId: this.projectId,
    };
  }

}
