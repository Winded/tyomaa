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

    authTokenGet(): Promise<TokenGetResponse> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}/auth/token`, {
                headers: this.headers(),
                data: {},
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    authTokenPost(body: TokenPostRequest): Promise<TokenPostResponse> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}/auth/token`, {
                headers: this.headers(),
                data: body,
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }

    clockGet(): Promise<ClockGetResponse> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}/clock`, {
                headers: this.headers(),
                data: {},
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    clockStartPost(body: ClockStartPostRequest): Promise<ClockStartPostResponse> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}/clock/start`, {
                headers: this.headers(),
                data: body,
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    clockStopPost(): Promise<ClockStopPostResponse> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}/clock/stop`, {
                headers: this.headers(),
                data: {},
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }

    entriesGet(query: EntriesGetRequest): Promise<EntriesGetResponse> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}/entries`, {
                headers: this.headers(),
                data: query,
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    entriesPost(body: EntriesPostRequest): Promise<EntriesPostResponse> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}/entries`, {
                headers: this.headers(),
                data: body,
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    entriesSingleGet(entryId: number): Promise<EntriesSingleGetResponse> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}/entries/${entryId}`, {
                headers: this.headers(),
                data: {},
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }
    entriesSinglePost(entryId: number, body: EntriesSinglePostRequest): Promise<EntriesSinglePostResponse> {
        return new Promise((resolve, reject) => {
            this.client.post(`${this.settings.host}/entries/${entryId}`, {
                headers: this.headers(),
                data: body,
            }, (data, _response) => {
                resolve(data);
            }).on('error', (err) => {
                reject(err);
            });
        });
    }

    projectsGet(): Promise<ProjectsGetResponse> {
        return new Promise((resolve, reject) => {
            this.client.get(`${this.settings.host}/projects`, {
                headers: this.headers(),
                data: {},
            }, (data, _response) => {
                resolve(data);
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