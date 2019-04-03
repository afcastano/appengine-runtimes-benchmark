import { Resolver, Query, Mutation } from '@nestjs/graphql';
import { Context, AllowAnonymous } from '@3wks/gae-node-nestjs';
import { DummiesRepository, DummyEntity } from './dummies.repository';
import { ConfigurationProvider, Ring } from '../config/config.provider';

@Resolver('DummyEntity')
export class DummiesResolver {
  constructor(
    private readonly repository: DummiesRepository,
    private readonly configuration: ConfigurationProvider
  ) {}

  @Query('dummies')
  @AllowAnonymous()
  async getDummies(
    _obj: {},
    { index }: { index: number},
    context: Context,
  ): Promise<ReadonlyArray<DummyEntity>> {
    const [entities] = await this.repository.query(context, {filters: {random2: [{op: ">=", value: index}, {op: "<", value: index + 10000}]}, limit: 10});
    return entities;
  }

  @Query('getDummyById')
  @AllowAnonymous()
  async getById(_obj: void, { id }: { id: string }, context: Context): Promise<DummyEntity | undefined> {
    return this.repository.get(context, id);
  }

  @Query('getSecret')
  @AllowAnonymous()
  async getSecret( _obj: {}, _args: { bla: string }, _context: Context): Promise<Ring> {
    const name = await this.configuration.apiKey
    return { name }
  }

  @Mutation('createDummy')
  @AllowAnonymous()
  async createDummy(
    _req: void,
    { id }: { id: string },
    context: Context,
  ) : Promise<DummyEntity> {
    const dummy: DummyEntity = {
      id,
      random1: 'fafsadf',
      random2: 123455
    };
    return await this.repository.save(context, dummy)
  }
}
