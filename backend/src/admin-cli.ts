import * as Vorpal from 'vorpal';
import { hashPassword } from './services/auth';

import './db/index';
import { User } from './db/user';

const prompt = new Vorpal();

prompt.command('users list', 'List all users').alias('u l')
    .action(async (args) => {
        let users = await User.findAll();

        prompt.log(`Listing ${users.length} users`);
        prompt.log('-----------------------------');

        users.forEach((user) => {
            prompt.log(user.name);
        });
    });

prompt.command('users create <name> <password>', 'Create new user').alias('u c')
    .action(async (args) => {
        let user = new User();
        user.name = args.name;
        user.password = await hashPassword(args.password);
        await user.save();

        prompt.log(`Created user with ID: ${user.id}`);
    });

export function run(): void {
    prompt.parse(process.argv);
}