import { Router, Request, Response } from 'express';
import { TimeEntry } from '../db/time-entry';
import { EntriesGetResponse, EntriesPostRequest, EntriesPostResponse, EntriesSingleResponse } from '../../../shared/api/entries';
import { Error as ApiError } from '../../../shared/api/error';
import { body, param } from 'express-validator/check';
import { validateNameIdentifier } from '../util/identity';
import { validateInput } from '../util/validation';

const router = Router();

router.get('/', async (req, res) => {
    let entries = await TimeEntry.findAll({
        where: {
            userId: req.session.user.id,
        },
        order: [['start', 'DESC']],
    });

    res.send(<EntriesGetResponse>{
        entries: entries.map(e => e.toApiFormat()),
    });
});

router.post('/', [
    body('project').exists().custom(validateNameIdentifier),
    body('start').exists().isISO8601(),
    body('end').exists().isISO8601(),
], async (req: Request, res: Response) => {
    if(!validateInput(req, res)) {
        return;
    }

    let body: EntriesPostRequest = {
        project: req.body.project,
        start: new Date(req.body.start),
        end: new Date(req.body.end),
    };

    let entry = new TimeEntry();
    entry.userId = req.session.user.id;
    entry.project = body.project;
    entry.start =body.start;
    entry.end =body.end;
    await entry.save();

    res.send(<EntriesPostResponse>{
        entry: entry.toApiFormat(),
    });
});

router.get('/:entryId', [
    param('entryId').exists().isNumeric(),
], async (req: Request, res: Response) => {
    if(!validateInput(req, res)) {
        return;
    }

    let entry = await TimeEntry.findOne({
        where: {
            id: req.params.entryId,
            userId: req.session.user.id,
        },
    });
    if(!entry) {
        res.status(404).send(<ApiError>{
            message: 'Entry not found',
        });
        return;
    }

    res.send(<EntriesSingleResponse>{
        entry: entry.toApiFormat(),
    });
});

router.put('/:entryId', [
    body('project').exists().custom(validateNameIdentifier),
    body('start').exists().isISO8601(),
    body('end').exists().isISO8601(),
], async (req: Request, res: Response, _next) => {
    if(!validateInput(req, res)) {
        return;
    }

    let entry = await TimeEntry.findOne({
        where: {
            id: req.params.entryId,
            userId: req.session.user.id,
        },
    });
    if(!entry) {
        res.status(404).send(<ApiError>{
            message: 'Entry not found',
        });
        return;
    }

    let body: EntriesPostRequest = {
        project: req.body.project,
        start: new Date(req.body.start),
        end: new Date(req.body.end),
    };

    entry.project = body.project;
    entry.start = body.start;
    entry.end = body.end;
    await entry.save();

    res.send(<EntriesSingleResponse>{
        entry: entry.toApiFormat(),
    });
});

export default router;