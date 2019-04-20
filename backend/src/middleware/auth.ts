import { Request, Response, NextFunction } from 'express';
import { User } from '../db/user';

export async function authentication(request: Request, _response: Response, next: NextFunction) {
    if(!request.session.tokenData || !(request.session.tokenData instanceof Object)) {
        next();
        return;
    }
    let tokenData = request.session.tokenData as { userId?: number };
    if(tokenData.userId === undefined) {
        next();
        return;
    }

    request.session.user = await User.findOne({ where: { id: tokenData.userId } });
    next();
}

export async function authorization(request: Request, response: Response, next: NextFunction) {
    if(!request.session.user) {
        response.status(403);
        response.send({ message: 'Not logged in' });
        return;
    }

    next();
}