import { Request, Response } from 'express';
import { validationResult, Result } from 'express-validator/check';

export function validateInput(request: Request, response: Response, onFail?: (errors: Result<{}>) => void): boolean {
    let errors = validationResult(request);
    if(!errors.isEmpty()) {
        if(onFail) {
            onFail(errors);
        }
        response.status(422).send({ errors: errors.array() });
        return false;
    }

    return true;
}