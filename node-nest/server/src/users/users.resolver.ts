import { Resolver, Query, Mutation } from '@nestjs/graphql';
import { Context, Roles } from '@3wks/gae-node-nestjs';
import { UserRepository, User } from './users.repository';
import { UsersService } from './users.service';

@Resolver('User')
export class UsersResolver {
  constructor(
    private readonly userRepository: UserRepository,
    private readonly userService: UsersService,
  ) {}

  @Query('users')
  async getUsers(
    _obj: {},
    _args: {},
    context: Context,
  ): Promise<ReadonlyArray<User>> {
    const [users] = await this.userRepository.query(context);

    return users;
  }

  @Query('userById')
  async getUserById(_obj: void, { id }: { id: string }, context: Context) {
    return this.userRepository.get(context, id);
  }

  @Roles('admin')
  @Mutation()
  async updateUser(
    _req: void,
    { id, name, roles }: { id: string; name: string; roles: string[] },
    context: Context,
  ) {
    return await this.userService.update(context, id, { name, roles });
  }

  avatar({ avatar }: User) {
    return {
      url: avatar,
    };
  }

  roles({ roles = [] }: User) {
    return roles;
  }
}