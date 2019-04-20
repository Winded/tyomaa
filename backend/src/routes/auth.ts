import { Router } from 'express';
import { body } from 'express-validator/check';
import { authenticateUser, createToken } from '../services/auth';
import { TokenGetResponse, TokenPostRequest, TokenPostResponse } from '../../../shared/api/auth';
import { Error as ApiError } from '../../../shared/api/error';
import { validateNameIdentifier } from '../util/identity';
import { validateInput } from '../util/validation';

const router = Router();

router.get('/token', (req, res, _next) => {
    res.send(<TokenGetResponse>{
        token: req.session.token ? req.session.token : null,
        user: req.session.user ? req.session.user.toApiFormat() : null,
    });
});

router.post('/token', [
    body('username').exists().custom(validateNameIdentifier),
    body('password').exists(),
], async (req, res, _next) => {
    if(!validateInput(req, res)) {
        return;
    }
    
    let body = req.body as TokenPostRequest;

    let user = await authenticateUser(body.username, body.password);
    if(!user) {
        res.status(400).send(<ApiError>{
            message: 'Invalid username or password',
        });
        return;
    }

    let t = createToken(user);

    res.send(<TokenPostResponse>{
        token: t,
        message: 'Authentication successful',
    });
});

export default router;