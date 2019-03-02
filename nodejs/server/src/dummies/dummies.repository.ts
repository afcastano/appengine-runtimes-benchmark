import { Repository, DatastoreProvider } from '@3wks/gae-node-nestjs';
import * as t from '@3wks/gae-node-nestjs/dist/validator';
import { Injectable } from '@nestjs/common';

const dummySchema = t.interface({
  id: t.string,
  random1: t.string,
  random2: t.number
});

export type DummyEntity = t.TypeOf<typeof dummySchema>;

@Injectable()
export class DummiesRepository extends Repository<DummyEntity> {
  constructor(datastore: DatastoreProvider) {
    super(datastore.datastore, 'DummyEntity', dummySchema, {
      index: {
        random2: true,
      },
      defaultValues: {
        random2: 0,
      },
    });
  }
}
