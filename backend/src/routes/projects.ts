import { Router } from 'express';
import * as cache from '../util/cache';
import { ProjectsGetResponse, Project } from '../../../shared/api/projects';
import { sequelize } from '../db/index';
import { QueryTypes } from 'sequelize';

const router = Router();

router.get('/', async (req, res) => {
    let projects = cache.get<ProjectsGetResponse>(req.session.user.id, 'projects');
    console.log(projects);
    if(projects) {
        res.send(projects);
        return;
    }
    
    let dbProjects = await sequelize.query<Project>('SELECT "project", SUM("end" - "start") as "totalTime" FROM "TimeEntries" WHERE "userId"=:userId GROUP BY "project"', {
        type: QueryTypes.SELECT,
        replacements: { userId: req.session.user.id },
    });
    // TODO convert postgresql interval

    projects = {
        projects: dbProjects,
    };

    cache.set<ProjectsGetResponse>(req.session.user.id, 'projects', projects);
    res.send(projects);
});

export default router;