import { Resolver, Query, Mutation } from '@nestjs/graphql';
import { Context, AllowAnonymous } from '@3wks/gae-node-nestjs';
import { DummiesRepository, DummyEntity } from './dummies.repository';

@Resolver('DummyEntity')
export class DummiesResolver {
  constructor(
    private readonly repository: DummiesRepository
  ) {}

  @Query('dummies')
  @AllowAnonymous()
  async getDummies(
    _obj: {},
    _args: {},
    context: Context,
  ): Promise<ReadonlyArray<DummyEntity>> {
    const [users] = await this.repository.query(context);

    return users;
  }

  @Query('getDummyById')
  @AllowAnonymous()
  async getById(_obj: void, { id }: { id: string }, context: Context): Promise<DummyEntity | undefined> {
    return this.repository.get(context, id);
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
