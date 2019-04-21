import * as express from 'express';
import * as bodyParser from 'body-parser';
import { sequelize } from './db/index';
import { tokenSession } from './middleware/token';
import { authentication, authorization } from './middleware/auth';
import authRouter from './routes/auth';
import entriesRouter from './routes/entries';
import clockRouter from './routes/clock';
import projectsRouter from './routes/projects';
import { User } from './db/user';

interface SessionData {
    token?: string,
    tokenData?: string | object;
    user?: User;
}

declare global {
    namespace Express {
        interface Request {
            session?: SessionData;
        }
    }
}

const app = express();

app.use((req, _res, next) => {
    req.session = {};
    next();
});
app.use(bodyParser.urlencoded());
app.use(bodyParser.json());
app.use(tokenSession);

const baseRouter = express.Router();
baseRouter.use(authentication);
baseRouter.use('/auth', authRouter);
baseRouter.use(authorization);
baseRouter.use('/entries', entriesRouter);
baseRouter.use('/clock', clockRouter);
baseRouter.use('/projects', projectsRouter);

app.use(process.env.ROOT_URL || '', baseRouter);

export async function runServer(): Promise<void> {
    let port = parseInt(process.env.PORT) || 80;

    await sequelize.sync();

    app.listen(port, () => {
        console.log(`Listening on port ${port}`);
    });
}