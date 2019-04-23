import * as Vorpal from 'vorpal';
import * as path from 'path';
import * as os from 'os';
import * as fs from 'fs';
import * as read from 'read';
import { ApiSettings, ApiClient } from './api-client'
import { IApiClient } from '../../shared/api/client';
import * as moment from 'moment';
import { table } from 'table';

interface CLISettings extends ApiSettings {}

function pad(n: string, width: number, z: string): string {
    z = z || '0';
    n = n + '';
    return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

const inputPrompt = (prompt: string, silent: boolean = false): Promise<string> => new Promise((resolve, reject) => {
    read({ prompt: prompt, silent: silent, replace: '*' }, function (er, value) {
        if (er) {
            reject(er);
        } else {
            resolve(value);
        }
    });
});

function formatDuration(duration: moment.Duration): string {
    return `${pad(duration.hours().toString(), 2, '0')}:${pad(duration.minutes().toString(), 2, '0')}:${pad(duration.seconds().toString(), 2, '0')}`;
}

function createPrompt(settings: CLISettings, apiClient: IApiClient, saveSettingsFunc: () => void): Vorpal {
    const prompt = new Vorpal();

    const authValidation = (args: Vorpal.Args): string | boolean => {
        if(!settings.host || !settings.token) {
            return 'ERROR: No host or auth token found. Please login first.';
        }
        return true;
    }

    prompt.command('auth status', 'Show authentication status')
        .action(async (_args) => {
            if(!settings.host || !settings.token) {
                prompt.log('Not authenticated');
                return;
            }
            
            let result = await apiClient.authTokenGet();
            if(!result.user) {
                prompt.log('Not authenticated');
                return;
            }

            prompt.log(`Authenticated as ${result.user.name}`);
        });
    prompt.command('auth login', 'Authenticate to a server')
        .action(async (_args) => {
            let host = await inputPrompt('Host URL: ');
            let username = await inputPrompt('Username: ');
            let password = await inputPrompt('Password: ', true);

            settings.host = host;
            let result = await apiClient.authTokenPost({
                username: username,
                password: password,
            });

            settings.token = result.token;
            saveSettingsFunc();

            prompt.log('Login successful');
        });

    prompt.command('clock status', 'Show current clock status')
        .validate(authValidation)
        .action(async (args) => {
            try {
                let result = await apiClient.clockGet();
                if(result.entry) {
                    let duration = moment.duration(moment().diff(moment(result.entry.start)));
                    prompt.log(`Clock running for ${result.entry.project} (${formatDuration(duration)})`);
                } else {
                    prompt.log('Clock not running');
                }
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });
    prompt.command('clock start <project>', 'Start clock on specified project')
        .validate(authValidation)
        .action(async (args) => {
            try {
                let result = await apiClient.clockStartPost({
                    project: args['project'],
                });
                prompt.log(`Clock started for project ${result.entry.project}`);
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });
    prompt.command('clock stop', 'Stop active clock')
        .validate(authValidation)
        .action(async (_args) => {
            try {
                let result = await apiClient.clockStopPost();
                prompt.log(`Clock stopped for project ${result.entry.project}`);
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });

    prompt.command('entries list', 'Show entries')
        .validate(authValidation)
        .action(async (_args) => {
            try {
                let result = await apiClient.entriesGet({});
                if(result.entries.length > 0) {
                    let uiData = [
                        ['Entry ID', 'Project name', 'Start time', 'End time'],
                    ];
                    result.entries.forEach(entry => {
                        let start = moment(entry.start).format('DD.MM.YYYY HH:mm:SS');
                        let end = entry.end ? moment(entry.end).format('DD.MM.YYYY HH:mm:SS') : 'None (ongoing)';
                        uiData.push([entry.id.toString(), entry.project, start, end]);
                    });
                    prompt.log(table(uiData));
                } else {
                    prompt.log('No entries found');
                }
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });
    prompt.command('entries create <project> <start> <end>', 'Create a new entry')
        .validate(authValidation)
        .validate((args) => {
            let d1 = new Date(args['start']);
            let d2 = new Date(args['end']);
            if(!isNaN(d1.getTime())) {
                return 'ERROR: Invalid start date';
            } else if(!isNaN(d2.getTime())) {
                return 'ERROR: Invalid end date';
            } else if(d2.getTime() > d1.getTime()) {
                return 'ERROR: End date is before start date';
            }

            return true;
        })
        .action(async (args) => {
            let project = args['project'];
            let start = new Date(args['start']).toISOString();
            let end = new Date(args['end']).toISOString();

            try {
                await apiClient.entriesPost({
                    project: project,
                    start: start,
                    end: end,
                });
                prompt.log('Entry created successfully');
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });
    prompt.command('entries view <entryid>', 'Show a specific entry')
        .validate(authValidation)
        .validate((args) => {
            let entryId = parseInt(args['entryid']);
            if(isNaN(entryId) || entryId < 1) {
                return 'Entry ID must be a number higher or equal to 1';
            }

            return true;
        })
        .action(async (args) => {
            let entryId = parseInt(args['entryid']);

            try {
                let result = await apiClient.entriesSingleGet(entryId);
                if(!result.entry) {
                    prompt.log('Entry not found');
                    return;
                }

                let start = moment(result.entry.start).format('DD.MM.YYYY HH:mm:SS');
                let end = result.entry.end ? moment(result.entry.end).format('DD.MM.YYYY HH:mm:SS') : 'None (ongoing)';
                prompt.log(table([
                    ['Entry ID', 'Project name', 'Start time', 'End time'],
                    [result.entry.id, result.entry.project, start, end],
                ]));
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });
    prompt.command('entries edit <entryid>', 'Modify an existing entry')
        .option('-p,--project <project>', 'Change the project of the entry')
        .option('-s,--start <start>', 'Change the start time of the entry')
        .option('-e,--end <end>', 'Change the end time of the entry')
        .validate(authValidation)
        .action(async (args) => {
            let entryId = parseInt(args['entryid']);

            try {
                let result = await apiClient.entriesSingleGet(entryId);
                if(!result.entry) {
                    prompt.log('ERROR: Entry does not exist');
                    return;
                }

                let project = args.options['project'] || result.entry.project;
                let start = args.options['start'] ? new Date(args['start']).toISOString() : result.entry.start;
                let end = args.options['end'] ? new Date(args['end']).toISOString() : result.entry.end;
                await apiClient.entriesSinglePost(entryId, {
                    project: project,
                    start: start,
                    end: end,
                });
                prompt.log('Entry modified successfully');
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });

    prompt.command('projects list', 'Show projects')
        .validate(authValidation)
        .action(async (_args) => {
            try {
                let result = await apiClient.projectsGet();
                if(result.projects.length > 0) {
                    let uiData = [
                        ['Project name', 'Total time'],
                    ];
                    result.projects.forEach(project => {
                        let totalTime = formatDuration(moment.duration(project.totalTime));
                        uiData.push([project.name, totalTime]);
                    });
                    prompt.log(table(uiData));
                } else {
                    prompt.log('No projects found');
                }
            } catch(err) {
                prompt.log(`ERROR: ${err}`);
            }
        });

    return prompt;
}

export function run(): void {
    let settingsFile = process.env.SETTINGS_FILE ? path.resolve(process.env.SETTINGS_FILE) : path.resolve(os.homedir(), '.tyomaa-cli');
    
    let settings: CLISettings = {
        host: null,
        token: null,
    };

    if(fs.existsSync(settingsFile)) {
        let sdata = fs.readFileSync(settingsFile);
        settings = JSON.parse(sdata.toString()) as CLISettings;
    }

    let client = new ApiClient(settings);
    let prompt = createPrompt(settings, client, () => {
        fs.writeFileSync(settingsFile, JSON.stringify(settings));
    });

    if(process.argv.length > 2) {
        prompt.parse(process.argv);
    }
    else {
        prompt.execSync('help');
    }
}