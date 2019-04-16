process.env.DB_IN_MEMORY = '1';

import 'mocha';
import { expect } from 'chai';
import { sequelize, User } from '../src/db';

describe('Database tests', () => {
    before(async () => {
        await sequelize.sync();
    });

    context('Users', () => {
        it('should create user', async () => {
            let user = await User.create({name: 'test', password: 'test'}) as User;
            expect(user.name).to.equal('test');
        });
    });
});