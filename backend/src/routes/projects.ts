import { Router } from 'express';
import * as cache from '../util/cache';
import { ProjectsGetResponse, Project } from '../../../shared/api/projects';
import { sequelize } from '../db/index';
import { QueryTypes } from 'sequelize';
import { IPostgresInterval } from 'postgres-interval';

const router = Router();

router.get('/', async (req, res) => {
    let projects = cache.get<ProjectsGetResponse>(req.session.user.id, 'projects');
    console.log(projects);
    if(projects) {
        res.send(projects);
        return;
    }
    
    let dbProjects = await sequelize.query('SELECT "project", SUM("end" - "start") as "totalTime" FROM "TimeEntries" WHERE "userId"=:userId GROUP BY "project"', {
        type: QueryTypes.SELECT,
        replacements: { userId: req.session.user.id },
    });

    projects = {
        projects: dbProjects.map((p: {project: string, totalTime: IPostgresInterval}) => (<Project>{
            name: p.project,
            totalTime: p.totalTime.toISO(),
        })),
    };

    cache.set<ProjectsGetResponse>(req.session.user.id, 'projects', projects);
    res.send(projects);
});

export default router;