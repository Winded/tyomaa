import * as express from 'express';
import * as bodyParser from 'body-parser';
import { sequelize } from './db/index';
import { tokenMiddleware } from './util/token';
import authRouter from './routes/auth';

let app = express();

app.use(bodyParser.json());
app.use(tokenMiddleware);

app.use('/auth', authRouter);

export async function runApp(): Promise<void> {
    let port = parseInt(process.env.PORT) || 80;

    await sequelize.sync();

    app.listen(port, () => {
        console.log(`Listening on port ${port}`);
    });
}