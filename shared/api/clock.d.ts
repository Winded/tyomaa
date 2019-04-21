import { Entry } from './entry';

export interface ClockGetResponse {
    entry: Entry;
}

export interface ClockStartPostRequest {
    project: string;
}
export interface ClockStartPostResponse extends ClockGetResponse {}

export interface ClockStopPostResponse extends ClockGetResponse {}