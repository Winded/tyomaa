import * as express from 'express';
import { sequelize } from './db/index';

let app = express();

app.get('/', (req, res, next) => {
    res.send({ message: 'asddddd!' });
});

export async function runApp(): Promise<void> {
    let port = parseInt(process.env.PORT) || 80;

    await sequelize.sync();

    app.listen(port, () => {
        console.log(`Listening on port ${port}`);
    });
}