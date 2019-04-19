import * as jwt from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';

interface TokenData {
    userId: number;
}

function instanceOfTokenData(object: any): object is TokenData {
    return 'userId' in object;
}

declare global {
    namespace Express {
        interface Request {
            session?: TokenData;
        }
    }
}

const secret = process.env.JWT_SECRET || 'default';
if(process.env.NODE_ENV == 'production' && secret == 'default') {
    throw new Error('Default JWT secret detected in production mode!');
}

export function sign(payload: TokenData): string {
    return jwt.sign(payload, secret, {
        expiresIn: '30d',
    });
}

export function verify(token: string): TokenData {
    let data = jwt.verify(token, secret);
    if(!data || !instanceOfTokenData(data)) {
        return null;
    }

    return data as TokenData;
}

export function tokenMiddleware(request: Request, response: Response, next: NextFunction): void {
    let token: string = request.body.token || request.query.token || request.headers['x-access-token'];
    if(!token || token == 'null') {
        request.session = null;
        next();
        return;
    }

    let data = verify(token);
    request.session = data;
    next();
}