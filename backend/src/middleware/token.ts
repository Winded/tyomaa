import { Request, Response, NextFunction } from 'express';
import { verify } from '../util/token';

export async function tokenSession(request: Request, _response: Response, next: NextFunction): Promise<void> {
    let token: string = request.body.token || request.query.token || request.headers['x-access-token'];
    if(!token || token == 'null') {
        next();
        return;
    }

    let data = await verify(token);
    request.session.token = token;
    request.session.tokenData = data;
    next();
}