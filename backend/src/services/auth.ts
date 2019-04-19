import * as bcrypt from 'bcrypt';
import { User } from '../db/user';
import * as token from '../util/token';

export async function authenticateUser(name: string, password: string): Promise<User> {
    let user = await User.findOne({
        where: { name: name }
    });
    if(!user) {
        return null;
    }

    if(!await bcrypt.compare(password, user.password)) {
        return null;
    }

    return user;
}

export function createToken(user: User): string {
    return token.sign({
        userId: user.id,
    });
}