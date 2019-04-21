import * as Vorpal from 'vorpal';
import * as path from 'path';
import * as os from 'os';
import * as fs from 'fs';
import * as read from 'read';
import { ApiSettings, ApiClient } from './api-client'
import { IApiClient } from '../../shared/api/client';

interface CLISettings extends ApiSettings {}

const inputPrompt = (prompt: string, silent: boolean = false): Promise<string> => new Promise((resolve, reject) => {
    read({ prompt: prompt, silent: silent, replace: '*' }, function (er, value) {
        if (er) {
            reject(er);
        } else {
            resolve(value);
        }
    });
});

function createPrompt(settings: CLISettings, apiClient: IApiClient, saveSettingsFunc: () => void): Vorpal {
    const prompt = new Vorpal();

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