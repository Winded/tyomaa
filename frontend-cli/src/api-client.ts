import { IApiClient } from '../../shared/api/client';
import { TokenGetResponse, TokenPostRequest, TokenPostResponse } from '../../shared/api/auth';
import { ClockGetResponse, ClockStartPostRequest, ClockStartPostResponse, ClockStopPostResponse } from '../../shared/api/clock';
import { EntriesGetResponse, EntriesGetRequest, EntriesPostRequest, EntriesPostResponse, EntriesSingleGetResponse, EntriesSinglePostRequest, EntriesSinglePostResponse } from '../../shared/api/entries';
import { ProjectsGetResponse } from '../../shared/api/projects';
import * as rest from 'node-rest-client';

export interface ApiSettings {
    host: string,
    token: string;
}

const defaultHeaders = {
    'Content-Type': 'application/json',
};

function resolveResponse(data, response) {
    
}

export class ApiClient implements IApiClient {
    private settings: ApiSettings;
    private client: rest.Client;

    constructor(settings: ApiSettings) {
        this.settings = settings;
        this.client = new rest.Client();
    }

    get token(): string {
        return this.settings.token;
    }
    set token(value: string) {
        this.settings.token = value;
    }

    async authTokenGet(): Promise<TokenGetResponse> {
        return await this.getRequest('/auth/token');
    }
    async authTokenPost(body: TokenPostRequest): Promise<TokenPostResponse> {
        return await this.postRequest('/auth/token', body);
    }

    async clockGet(): Promise<ClockGetResponse> {
        return await this.getRequest('/clock');
    }
    async clockStartPost(body: ClockStartPostRequest): Promise<ClockStartPostResponse> {
        return await this.postRequest('/clock/start', body);
    }
    async clockStopPost(): Promise<ClockStopPostResponse> {
        return await this.postRequest('/clock/stop');
    }

    async entriesGet(query: EntriesGetRequest): Promise<EntriesGetResponse> {
        return await this.getRequest('/entries', query);
    }
    async entriesPost(body: EntriesPostRequest): Promise<EntriesPostResponse> {
        return await this.postRequest('/entries', body);
    }
    async entriesSingleGet(entryId: number): Promise<EntriesSingleGetResponse> {
        return await this.getRequest(`/entries/${entryId}`);
    }
    async entriesSinglePost(entryId: number, body: EntriesSinglePostRequest): Promise<EntriesSinglePostResponse> {
        return await this.postRequest(`/entries/${entryId}`, body);
    }

    async projectsGet(): Promise<ProjectsGetResponse> {
        return await this.getRequest(`/projects`);
    }

    private getRequest(url: string, data: any = {}): Promise<any> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}${url}`, {
                headers: this.headers(),
                data: data,
            }, (data, response) => {
                if(Math.floor(response.statusCode / 100) == 2) {
                    resolve(data);
                } else {
                    reject(data);
                }
            }).on('error', (err) => {
                reject(err);
            });
        });
    }

    private postRequest(url: string, data: any = {}): Promise<any> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}${url}`, {
                headers: this.headers(),
                data: data,
            }, (data, response) => {
                if(Math.floor(response.statusCode / 100) == 2) {
                    resolve(data);
                } else {
                    reject(data);
                }
            }).on('error', (err) => {
                reject(err);
            });
        });
    }

    private headers(): object {
        let h = {
            ...defaultHeaders,
        };
        if(this.settings.token) {
            h['x-access-token'] = this.settings.token;
        }
        return h;
    }
}