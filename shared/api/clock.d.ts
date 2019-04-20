import { Entry } from './entry';

export interface ClockGetResponse {
    entry: Entry;
}

export interface ClockStartPostRequest {
    project: string;
}