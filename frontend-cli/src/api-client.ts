import { IApiClient } from '../../shared/api/client';
import { TokenGetResponse, TokenPostRequest, TokenPostResponse } from '../../shared/api/auth';
import { ClockGetResponse, ClockStartPostRequest, ClockStartPostResponse } from '../../shared/api/clock';
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
        throw new Error("Method not implemented.");
    }
    clockStartPost(body: ClockStartPostRequest): Promise<ClockStartPostResponse> {
        throw new Error("Method not implemented.");
    }
    clockStopPost(): Promise<void> {
        throw new Error("Method not implemented.");
    }

    entriesGet(query: EntriesGetRequest): Promise<EntriesGetResponse> {
        throw new Error("Method not implemented.");
    }
    entriesPost(body: EntriesPostRequest): Promise<EntriesPostResponse> {
        throw new Error("Method not implemented.");
    }
    entriesSingleGet(entryId: number): Promise<EntriesSingleGetResponse> {
        throw new Error("Method not implemented.");
    }
    entriesSinglePost(entryId: number, body: EntriesSinglePostRequest): Promise<EntriesSinglePostResponse> {
        throw new Error("Method not implemented.");
    }

    projectsGet(): Promise<ProjectsGetResponse> {
        throw new Error("Method not implemented.");
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