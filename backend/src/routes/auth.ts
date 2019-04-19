import { Router } from 'express';
import * as token from '../util/token';
import { authenticateUser, createToken } from '../services/auth';

const router = Router();

router.get('/token', async (req, res, next) => {
    res.send({ userId: req.session ? req.session.userId : null });
});

router.post('/token', async (req, res, next) => {
    let username = req.body.username;
    let password = req.body.password;

    let user = await authenticateUser(username, password);
    if(!user) {
        res.send({
            success: false,
            token: null,
            message: 'Invalid username or password',
        });
        return;
    }

    let t = createToken(user);

    res.send({
        success: true,
        token: t,
        message: 'Authentication successful',
    });
});

export default router;