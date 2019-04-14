import * as express from 'express';

let app = express();

app.get('/', (req, res, next) => {
    res.send({ message: 'asddddd!' });
});

export function runApp(port: number): void {
    app.listen(port, () => {
        console.log(`Listening on port ${port}`);
    });
}