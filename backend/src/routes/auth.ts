import { Router } from 'express';
import { check, validationResult } from 'express-validator/check';
import { authenticateUser, createToken } from '../services/auth';
import { TokenGetResponse, TokenPostRequest, TokenPostResponse } from '../../../shared/api/auth';
import { Error as ApiError } from '../../../shared/api/error';

const router = Router();

router.get('/token', async (req, res, _next) => {
    res.send(<TokenGetResponse>{
        token: req.session.token ? req.session.token : null,
        user: req.session.user ? req.session.user.toApiFormat() : null,
    });
});

router.post('/token', [
    check('username').not().isEmpty(),
    check('password').not().isEmpty(),
], async (req, res, _next) => {
    let errors = validationResult(req);
    if(!errors.isEmpty()) {
        res.status(422).send({ errors: errors.array() });
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