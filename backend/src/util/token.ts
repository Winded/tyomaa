import * as jwt from 'jsonwebtoken';

const secret = process.env.JWT_SECRET || 'default';
if(process.env.NODE_ENV == 'production' && secret == 'default') {
    throw new Error('Default JWT secret detected in production mode!');
}

export function sign(payload: string | object): string {
    return jwt.sign(payload, secret, {
        expiresIn: '30d',
    });
}

export function verify(token: string): Promise<string | object> {
    return new Promise((resolve, reject) => {
        jwt.verify(token, secret, (err, data) => {
            if(err) {
                resolve(null);
                return;
            }

            resolve(data);
        });
    });
}