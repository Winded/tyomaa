process.env.DB_IN_MEMORY = '1';

import 'mocha';
import { expect } from 'chai';
import { sequelize } from '../src/db/index';
import { User } from '../src/db/user';
import { TimeEntry } from '../src/db/time-entry';

describe('Database tests', () => {
    before(async () => {
        await sequelize.sync();
    });

    afterEach(async () => {
        await User.destroy({ where: {}, force: true });
        await TimeEntry.destroy({ where: {}, force: true });
    });

    context('Users', () => {
        it('should create user', async () => {
            let user = await User.create({name: 'test', password: 'test'});
            let count = await User.count();
            expect(user.name).to.equal('test');
            expect(count).to.equal(1);
        });
        it('should delete created user', async () => {
            await User.create({name: 'test', password: 'test'});
            let deleted = await User.destroy({ where: { name: 'test' } });
            let count = await User.count();
            expect(count).to.equal(0);
            expect(deleted).to.equal(1);
        });
    });
});