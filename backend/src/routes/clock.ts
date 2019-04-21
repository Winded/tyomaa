import { Request, Response, Router } from "express";
import { validateInput } from "../util/validation";
import { body } from "express-validator/check";
import { validateNameIdentifier } from "../util/identity";
import { TimeEntry } from "../db/time-entry";
import { Error as ApiError } from '../../../shared/api/error';
import { ClockGetResponse, ClockStartPostRequest, ClockStopPostResponse, ClockStartPostResponse } from '../../../shared/api/clock';
import * as cache from "../util/cache";

const router = Router();

router.get('/', async (req, res) => {
    let activeEntry = await TimeEntry.findOne({
        where: {
            userId: req.session.user.id,
            end: null,
        },
    });

    res.send(<ClockGetResponse>{
        entry: activeEntry ? activeEntry.toApiFormat() : null,
    });
});

router.post('/start', [
    body('project').exists().custom(validateNameIdentifier),
], async (req: Request, res: Response) => {
    if(!validateInput(req, res)) {
        return;
    }

    let activeEntry = await TimeEntry.findOne({
        where: {
            userId: req.session.user.id,
            end: null,
        },
    });
    if(activeEntry) {
        res.status(400).send(<ApiError>{
            message: 'Active entry already exists',
        });
        return;
    }

    let body = req.body as ClockStartPostRequest;

    activeEntry = new TimeEntry();
    activeEntry.userId = req.session.user.id;
    activeEntry.project = body.project;
    activeEntry.start = new Date();
    await activeEntry.save();

    cache.clear(req.session.user.id);

    res.send(<ClockStartPostResponse>{
        entry: activeEntry.toApiFormat(),
    });
});

router.post('/stop', async (req: Request, res: Response) => {
    let activeEntry = await TimeEntry.findOne({
        where: {
            userId: req.session.user.id,
            end: null,
        },
    });
    if(!activeEntry) {
        res.status(400).send(<ApiError>{
            message: 'Active entry does not exist',
        });
        return;
    }
    
    activeEntry.end = new Date();
    await activeEntry.save();

    cache.clear(req.session.user.id);

    res.send(<ClockStopPostResponse>{
        entry: activeEntry.toApiFormat(),
    });
});

export default router;