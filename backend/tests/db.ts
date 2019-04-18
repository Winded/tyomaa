process.env.DB_IN_MEMORY = '1';

import 'mocha';
import { expect, use as chaiUse } from 'chai';
import * as chaiAP from 'chai-as-promised';
import { sequelize } from '../src/db/index';
import { User } from '../src/db/user';
import { TimeEntry } from '../src/db/time-entry';

const testUser = {
    name: 'test',
    password: '$2b$08$cb4gSSsDyc1OjEors4n3dOVaClEnEgm6K0qQLWJK4YyM9X4q7PzKG',
};

describe('Database tests', () => {
    before(async () => {
        await sequelize.sync();
        chaiUse(chaiAP);
    });

    afterEach(async () => {
        await User.destroy({ where: {}, force: true });
        await TimeEntry.destroy({ where: {}, force: true });
    });

    context('Users', () => {
        it('should create user', async () => {
            let user = await User.create(testUser);
            let count = await User.count();
            expect(user.name).to.equal('test');
            expect(count).to.equal(1);
        });
        it('should update user name and password', async () => {
            let user = await User.create(testUser);
            user.name = 'test2';
            user.password = '$ff$ff$abcdefghijklmopqrstuwxyz0123456789';
            await user.save();
            user = await User.findOne({ where: { name: 'test2' } });
            let count = await User.count();
            expect(count).to.equal(1);
            expect(user).to.not.equal(null);
            expect(user.name).to.equal('test2');
            expect(user.password).to.equal('$ff$ff$abcdefghijklmopqrstuwxyz0123456789');
        });
        it('should fail to create user with invalid name', async () => {
            expect((async () => {
                return await User.create({...testUser, name: 'test invalid name'});
            })()).to.be.rejectedWith('Name must only contain alphabetic characters, numbers and dashes');
        });
        it('should fail to update user with invalid name', async () => {
            let user = await User.create(testUser);
            expect((async () => {
                user.name = 'test invalid name';
                await user.save();
            })()).to.be.rejectedWith('Name must only contain alphabetic characters, numbers and dashes');
            user = await User.findOne({ where: { name: 'test' } });
            let count = await User.count();
            expect(count).to.equal(1);
            expect(user).to.not.equal(null);
            expect(user.name).to.equal('test');
        });
        it('should delete created user', async () => {
            await User.create(testUser);
            let deleted = await User.destroy({ where: { name: 'test' } });
            let count = await User.count();
            expect(count).to.equal(0);
            expect(deleted).to.equal(1);
        });
    });
});